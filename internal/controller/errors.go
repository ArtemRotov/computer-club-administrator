package controller

import (
	"errors"
	"fmt"
)

var (
	ErrCannotParseTimeValue  = errors.New("cannot parse time value")
	ErrCannotBadTimeInterval = errors.New("bad time interval")
	ErrNotMatchEventPattern  = errors.New("string does not match event patterns")
	ErrCannotParseEventId    = errors.New("cannot parse event id")
	ErrCannotParseDeskId     = errors.New("cannot parse desk id")
)

func newLineError(line int, err error) error {
	if err != nil {
		return fmt.Errorf("error at line %d - %w", line, err)
	}
	return errors.New(fmt.Sprintf("error at line %d", line))
}
