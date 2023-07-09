package service

import (
	"github.com/ArtemRotov/computer-club-manager/internal/model"
	"testing"
	"time"
)

func TestManagerService_rentDuration(t *testing.T) {
	s := &ManagerService{
		desks: make(map[int]*model.Desk, 0),
	}

	tc := []struct {
		name     string
		desk     *model.Desk
		expected string
	}{
		{
			name:     "rentDuration #1 (zero duration)",
			desk:     &model.Desk{RentDuration: []time.Duration{}},
			expected: "00:00",
		},
		{
			name: "rentDuration #2",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*2 + time.Minute*38),
			}},
			expected: "02:38",
		},
		{
			name: "rentDuration #3",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*13 + time.Minute*0),
			}},
			expected: "13:00",
		},
		{
			name: "rentDuration #4",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*2 + time.Minute*38),
				time.Duration(time.Hour*13 + time.Minute*0),
			}},
			expected: "15:38",
		},

		{
			name: "rentDuration #5",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*2 + time.Minute*38),
				time.Duration(time.Hour*13 + time.Minute*0),
				time.Duration(time.Hour*24 + time.Minute*0),
			}},
			expected: "39:38",
		},
	}
	for i, tt := range tc {
		s.desks[i] = tt.desk
		t.Run(tt.name, func(t *testing.T) {
			if got := s.rentDuration(i); got != tt.expected {
				t.Errorf("rentDuration() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestManagerService_revenue(t *testing.T) {
	s := &ManagerService{
		desks:  make(map[int]*model.Desk, 0),
		config: &model.ClubConfiguration{Price: 10},
	}

	tc := []struct {
		name     string
		desk     *model.Desk
		expected int
	}{
		{
			name:     "revenue #1",
			desk:     &model.Desk{RentDuration: []time.Duration{}},
			expected: 0,
		},
		{
			name: "revenue #2",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*1 + time.Minute*0),
			}},
			expected: 10,
		},
		{
			name: "revenue #3",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*1 + time.Minute*0),
				time.Duration(time.Hour*0 + time.Minute*1),
			}},
			expected: 20,
		},
		{
			name: "revenue #4",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*1 + time.Minute*0),
				time.Duration(time.Hour*0 + time.Minute*1),
				time.Duration(time.Hour*0 + time.Minute*59),
			}},
			expected: 30,
		},
		{
			name: "revenue #5",
			desk: &model.Desk{RentDuration: []time.Duration{
				time.Duration(time.Hour*1 + time.Minute*23),
				time.Duration(time.Hour*4 + time.Minute*17),
				time.Duration(time.Hour*2 + time.Minute*22),
				time.Duration(time.Hour*4 + time.Minute*20),
				time.Duration(time.Hour*1 + time.Minute*0),
			}},
			expected: 160,
		},
	}
	for i, tt := range tc {
		s.desks[i] = tt.desk
		t.Run(tt.name, func(t *testing.T) {
			if got := s.revenue(i); got != tt.expected {
				t.Errorf("revenue() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
