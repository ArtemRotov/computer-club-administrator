package model

import "time"

type ClubConfiguration struct {
	TableCount  int
	Price       int
	OpeningTime time.Time
	ClosingTime time.Time
}
