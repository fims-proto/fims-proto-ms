package domain

import "strconv"

type Formatter struct {
	prefix string
	suffix string
}

func NewFormatter(prefix string, suffix string) Formatter {
	return Formatter{
		prefix: prefix,
		suffix: suffix,
	}
}

func (f Formatter) format(count uint) string {
	return f.prefix + strconv.Itoa(int(count)) + f.suffix
}
