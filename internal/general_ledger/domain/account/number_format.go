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
