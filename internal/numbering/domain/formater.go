package counter

import (
	"github.com/pkg/errors"
	"fmt"
)

type Formater struct {
	length uint
}

func (f Formater) format(count uint) (string, error){
	count_str := fmt.Sprintf("%0*d",f.length,count)
	if len(count_str) > int(f.length){
		return "", errors.New("Voucher Number exceeds the predefined length")
	}
	return count_str, nil
}

func (f Formater) SetLen(len uint) {
	f.length = len
}