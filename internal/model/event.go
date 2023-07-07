package model

import (
	"fmt"
	"time"
)

type EventId int

const (
	ClientHasCome EventId = iota + 1
	ClientTookTheTable
	ClientIsWaiting
	ClientLeft

	ClientLeftAtClosing = iota + 7
	ClientTookTheTableAfterWaiting
	EventError
)

type Event struct {
	Time       time.Time
	Id         EventId
	ClientName string
	DeskId     int
	ErrMsg     string
}

func (e *Event) String() string {
	var s string

	s = fmt.Sprintf("%s %d", e.Time.Format(time.TimeOnly)[:5], e.Id)

	if len(e.ClientName) > 0 {
		s += fmt.Sprintf(" %s", e.ClientName)
		if e.DeskId > 0 {
			s += fmt.Sprintf(" %d", e.DeskId)
		}
	} else if len(e.ErrMsg) > 0 {
		s += fmt.Sprintf(" %s", e.ErrMsg)
	}

	return s
}
