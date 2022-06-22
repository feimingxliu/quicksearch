package errors

import (
	"runtime/debug"
	"strings"
)

var sb = &strings.Builder{}

type withStack struct {
	err   error
	stack []byte
}

func (e *withStack) Error() string {
	sb.WriteString(e.err.Error())
	sb.WriteString("\n")
	sb.Write(e.stack)
	s := sb.String()
	sb.Reset()
	return s
}

func WithStack(err error) error {
	return &withStack{
		err:   err,
		stack: debug.Stack(),
	}
}
