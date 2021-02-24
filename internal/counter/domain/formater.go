package counter

import "strconv"

type Formater struct {
	prefix string
	sufix  string
}

func NewFormater(prefix string, sufix string) Formater {
	return Formater{
		prefix: prefix,
		sufix:  sufix,
	}
}

func (f Formater) format(count uint) string {
	return f.prefix + strconv.Itoa(int(count)) + f.sufix
}
