package service

import "errors"

var (
	ErrNotMatchEventPattern = errors.New("string does not match event patterns")
	ErrCannotParseTimeValue = errors.New("cannot parse time value")
	ErrCannotParseEventId   = errors.New("cannot parse event id")
	ErrCannotParseDeskId    = errors.New("cannot parse desk id")
)
