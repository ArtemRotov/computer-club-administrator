package model

import "time"

type Desk struct {
	Id           int
	IsBusy       bool
	RentDuration []time.Duration
	currTime     time.Time
}

func (d *Desk) Take(t time.Time) {
	if !d.IsBusy {
		d.IsBusy = true
		d.currTime = t
	}
}

func (d *Desk) Free(t time.Time) {
	if d.IsBusy {
		d.IsBusy = false
		d.RentDuration = append(d.RentDuration, t.Sub(d.currTime))
		d.currTime = time.Time{}
	}
}
