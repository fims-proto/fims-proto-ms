package domain

import "strconv"

type Formatter struct {
	prefix string
	sufix  string
}

func NewFormatter(prefix string, sufix string) Formatter {
	return Formatter{
		prefix: prefix,
		sufix:  sufix,
	}
}

func (f Formatter) format(count uint) string {
	return f.prefix + strconv.Itoa(int(count)) + f.sufix
}
