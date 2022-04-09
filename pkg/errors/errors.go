package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"runtime"
)

func callers() []uintptr {
	var pcs [32]uintptr
	l := runtime.Callers(3, pcs[:])
	return pcs[:l]
}

// Error with caller stack information
type Error interface {
	error
	t()
}

type item struct {
	msg   string
	stack []uintptr
}

func (i *item) Error() string {
	return i.msg
}

func (i *item) t() {}

func (i *item) Format(s fmt.State, verb rune) {
	_, _ = io.WriteString(s, i.msg)
	_, _ = io.WriteString(s, "\n")

	for _, pc := range i.stack {
		_, err := fmt.Fprintf(s, "%+v\n", errors.Frame(pc))
		if err != nil {
			panic(errors.New("error stack print error"))
		}
	}
}

// New create a new error
func New(msg string) Error {
	return &item{msg: msg, stack: callers()}
}

// Errorf create a new error with format
func Errorf(format string, args ...interface{}) Error {
	return &item{msg: fmt.Sprintf(format, args), stack: callers()}
}

// Wrap can add some extra message into err
func Wrap(err error, msg string) Error {
	if err == nil {
		return nil
	}
	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// Wrapf can add some extra message into err with format
func Wrapf(err error, format string, args ...interface{}) Error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, args...)

	e, ok := err.(*item)
	if !ok {
		return &item{msg: fmt.Sprintf("%s; %s", msg, err.Error()), stack: callers()}
	}

	e.msg = fmt.Sprintf("%s; %s", msg, e.msg)
	return e
}

// WithStack add caller stack information
func WithStack(err error) Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*item); ok {
		return e
	}

	return &item{msg: err.Error(), stack: callers()}
}
