package log

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	sd = "[ debug ]: "
	si = "[ info ] : "
	se = "[ error ]: "
)

func TestLog(t *testing.T) {
	t.Parallel()

	logAssert := assert.New(t)

	mle := &mockLogEnabler{}
	mo := &mockOutput{}

	InitLoggers(mle.getEnabler(), mo)

	mle.setLvl(0)
	Debugln("test")
	logAssert.True(strings.HasPrefix(mo.Line, sd))
	logAssert.True(strings.HasSuffix(mo.Line, "test\n"))
	Infoln("test")
	logAssert.False(strings.HasPrefix(mo.Line, si))
	Errorln("test")
	logAssert.False(strings.HasPrefix(mo.Line, se))

	mle.setLvl(1)
	Infoln("test")
	logAssert.True(strings.HasPrefix(mo.Line, si))
	logAssert.True(strings.HasSuffix(mo.Line, "test\n"))
	Debugln("test")
	logAssert.False(strings.HasPrefix(mo.Line, sd))
	Errorln("test")
	logAssert.False(strings.HasPrefix(mo.Line, se))

	mle.setLvl(2)
	Errorln("test")
	logAssert.True(strings.HasPrefix(mo.Line, se))
	logAssert.True(strings.HasSuffix(mo.Line, "test\n"))
	Infoln("test")
	logAssert.False(strings.HasPrefix(mo.Line, si))
	Debugln("test")
	logAssert.False(strings.HasPrefix(mo.Line, sd))
}

type mockLogEnabler struct {
	lvl int
}

func (le *mockLogEnabler) getEnabler() func(lvl int, v ...interface{}) bool {
	return func(lvl int, v ...interface{}) bool {
		return le.lvl == lvl
	}
}

func (le *mockLogEnabler) setLvl(lvl int) {
	le.lvl = lvl
}

type mockOutput struct {
	Line string
}

func (o *mockOutput) Write(p []byte) (n int, err error) {
	o.Line = string(p)
	return len(p), nil
}
