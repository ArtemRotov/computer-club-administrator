package model

import "time"

type Desk struct {
	Id           int
	IsBusy       bool
	RentDuration []time.Duration
}
