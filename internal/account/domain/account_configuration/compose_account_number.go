package account_configuration

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func (ac AccountConfiguration) ComposeAccountNumber(codeLengths []int) (string, error) {
	if len(ac.numberHierarchy) > len(codeLengths) {
		return "", errors.Errorf("account number hierarchy %d exceeds max depth %d", len(ac.numberHierarchy), len(codeLengths))
	}

	for i := 0; i < len(ac.numberHierarchy); i++ {
		if ac.numberHierarchy[i] < 1 {
			return "", errors.Errorf("account number %d at level %d cannot be smaller than 1", ac.numberHierarchy[i], i)
		}
		if len(strconv.Itoa(ac.numberHierarchy[i])) > codeLengths[i] {
			return "", errors.Errorf("account number %d at level %d exceeds max length (%d)", ac.numberHierarchy[i], i, codeLengths[i])
		}
	}

	var builder strings.Builder
	for i, number := range ac.numberHierarchy {
		builder.WriteString(fmt.Sprintf("%0*d", codeLengths[i], number))
	}

	return builder.String(), nil
}
