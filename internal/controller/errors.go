package controller

import (
	"errors"
	"fmt"
)

var (
	ErrCannotParseTimeValue = errors.New("cannot parse time value")
	ErrNotMatchEventPattern = errors.New("string does not match event patterns")
	ErrCannotParseEventId   = errors.New("cannot parse event id")
	ErrCannotParseDeskId    = errors.New("cannot parse desk id")
)

func newLineError(line int, str string) error {
	return errors.New(fmt.Sprintf("error at line %d: '%s'", line, str))
}
