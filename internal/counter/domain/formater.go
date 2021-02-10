package counter

import (
	"github.com/pkg/errors"
	"fmt"
)

type Formater struct {
	length uint // length of zero-padded number
	prefix string 
	sufix string
}

func NewFormater(len uint, prefix string, sufix string) Formater{
	return Formater{
		length: len,
		prefix: prefix,
		sufix: sufix,
	}
}

func (f Formater) format(count uint) (string, error){
	count_str := fmt.Sprintf("%0*d",f.length,count)
	if len(count_str) > int(f.length){
		return "", errors.Errorf("Counter Number exceeds the predefined length %d",f.length)
	}
	return f.prefix+count_str+f.sufix, nil
}

func (f Formater) SetLen(len uint) {
	f.length = len
}