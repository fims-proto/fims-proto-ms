package account

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	SegmentWidth = 6
	MaxLevels    = 10
)

// PadRawAccountNumber converts a hierarchy slice to a max-padded string.
// Each segment is zero-padded to exactly SegmentWidth (6) digits.
// Example: [3104, 2] → "003104000002"
func PadRawAccountNumber(hierarchy []int) (string, error) {
	if len(hierarchy) == 0 {
		return "", fmt.Errorf("account hierarchy must not be empty")
	}
	if len(hierarchy) > MaxLevels {
		return "", fmt.Errorf("account hierarchy exceeds max levels (%d)", MaxLevels)
	}

	var sb strings.Builder
	for _, seg := range hierarchy {
		if seg < 0 || seg >= 1_000_000 {
			return "", fmt.Errorf("segment %d out of range [0, 999999]", seg)
		}
		if _, err := fmt.Fprintf(&sb, "%06d", seg); err != nil {
			return "", fmt.Errorf("failed to format segment %d: %w", seg, err)
		}
	}
	return sb.String(), nil
}

// HierarchyFromRaw extracts the integer segments from a raw account number.
// Each segment is SegmentWidth (6) digits wide.
// Example: "003104000002" → [3104, 2]
func HierarchyFromRaw(raw string) ([]int, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("raw account number must not be empty")
	}
	if len(raw)%SegmentWidth != 0 {
		return nil, fmt.Errorf("invalid raw account number length: %d (must be multiple of %d)", len(raw), SegmentWidth)
	}

	levels := len(raw) / SegmentWidth
	result := make([]int, levels)
	for i := 0; i < levels; i++ {
		seg := raw[i*SegmentWidth : (i+1)*SegmentWidth]
		val, err := strconv.Atoi(seg)
		if err != nil {
			return nil, fmt.Errorf("invalid segment %q: %w", seg, err)
		}
		result[i] = val
	}
	return result, nil
}

// ReadableFromRaw derives the human-readable account number given SoB code lengths.
// Each segment from the raw number is extracted and then zero-padded to the width
// specified in codeLengths.
// Example: raw="003104000002", codeLengths=[4,2] → "310402"
func ReadableFromRaw(raw string, codeLengths []int) (string, error) {
	hierarchy, err := HierarchyFromRaw(raw)
	if err != nil {
		return "", err
	}

	if len(codeLengths) < len(hierarchy) {
		return "", fmt.Errorf("codeLengths has fewer entries (%d) than hierarchy levels (%d)",
			len(codeLengths), len(hierarchy))
	}

	var sb strings.Builder
	for i, val := range hierarchy {
		if _, err = fmt.Fprintf(&sb, "%0*d", codeLengths[i], val); err != nil {
			return "", fmt.Errorf("failed to format hierarchy %d: %w", i, err)
		}
	}
	return sb.String(), nil
}

// RawFromReadable converts a human-readable account number to raw format given codeLengths.
// Accepts any valid prefix: with codeLengths=[4,2,2], accepts "1001", "100101", or "10010101"
func RawFromReadable(readable string, codeLengths []int) (string, error) {
	if len(readable) == 0 {
		return "", fmt.Errorf("readable account number must not be empty")
	}

	var hierarchy []int
	pos := 0

	for i, cl := range codeLengths {
		end := pos + cl
		if end > len(readable) {
			break
		}

		val, err := strconv.Atoi(readable[pos:end])
		if err != nil {
			return "", fmt.Errorf("non-numeric segment at level %d: %w", i, err)
		}
		hierarchy = append(hierarchy, val)
		pos = end
	}

	if pos != len(readable) {
		return "", fmt.Errorf("readable number has %d trailing characters not matching codeLengths %v", len(readable)-pos, codeLengths)
	}

	if len(hierarchy) == 0 {
		return "", fmt.Errorf("readable account number does not match any level in codeLengths %v", codeLengths)
	}

	return PadRawAccountNumber(hierarchy)
}

// LevelFromRaw derives the depth level (1-based) from a raw account number.
func LevelFromRaw(raw string) int {
	return len(raw) / SegmentWidth
}

// AppendRawAccountNumber appends levelNumber as a new level segment to superiorRaw,
// producing the child's raw account number.
// Pass superiorRaw = "" for a root (level-1) account.
func AppendRawAccountNumber(superiorRaw string, levelNumber int) (string, error) {
	if superiorRaw == "" {
		return PadRawAccountNumber([]int{levelNumber})
	}
	hierarchy, err := HierarchyFromRaw(superiorRaw)
	if err != nil {
		return "", fmt.Errorf("invalid superior raw account number: %w", err)
	}
	return PadRawAccountNumber(append(hierarchy, levelNumber))
}
