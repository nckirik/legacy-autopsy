package cdl

import (
	"errors"
	"fmt"
)

type codedError struct {
	code string
	msg  string
}

func (e codedError) Error() string { return e.code + ": " + e.msg }

func cdlerr(code, format string, args ...any) error {
	return codedError{code: code, msg: fmt.Sprintf(format, args...)}
}

func diagCode(err error) string {
	var c codedError
	if errors.As(err, &c) {
		return c.code
	}
	return "CDL_PARSE"
}

func diagMessage(err error) string {
	var c codedError
	if errors.As(err, &c) {
		return c.msg
	}
	return err.Error()
}

// errAborted signals that a nested parse already recorded a diagnostic.
var errAborted = errors.New("cdl: parse aborted")
