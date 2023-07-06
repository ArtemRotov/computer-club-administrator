package model

import "time"

type ClubConfiguration struct {
	DeskCount   int
	Price       int
	OpeningTime time.Time
	ClosingTime time.Time
}
