package account

import (
	"fmt"
	"strconv"
	"strings"
)

func ComposeAccountNumber(numberHierarchy, codeLengths []int) (string, error) {
	if len(numberHierarchy) > len(codeLengths) {
		return "", fmt.Errorf("account number hierarchy %d exceeds max depth %d", len(numberHierarchy), len(codeLengths))
	}

	for i := 0; i < len(numberHierarchy); i++ {
		if numberHierarchy[i] < 1 {
			return "", fmt.Errorf("account number %d at level %d cannot be smaller than 1", numberHierarchy[i], i)
		}
		if len(strconv.Itoa(numberHierarchy[i])) > codeLengths[i] {
			return "", fmt.Errorf("account number %d at level %d exceeds max length (%d)", numberHierarchy[i], i, codeLengths[i])
		}
	}

	var builder strings.Builder
	for i, number := range numberHierarchy {
		builder.WriteString(fmt.Sprintf("%0*d", codeLengths[i], number))
	}

	return builder.String(), nil
}
