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
}

func (e *Event) String() string {
	if e.DeskId > 0 {
		return fmt.Sprintf("%s %d %s %d", e.Time.Format(time.TimeOnly)[:5], e.Id, e.ClientName, e.DeskId)
	}

	return fmt.Sprintf("%s %d %s", e.Time.Format(time.TimeOnly)[:5], e.Id, e.ClientName)
}
