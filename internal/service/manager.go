package service

import (
	"fmt"
	"github.com/ArtemRotov/computer-club-administrator/internal/model"
	"sort"
	"time"
)

type ManagerService struct {
	desks   map[int]*model.Desk
	clients map[string]*model.Client
	queue   []*model.Client
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		desks:   nil,
		clients: nil,
		queue:   nil,
	}
}

func (s *ManagerService) Run(config *model.ClubConfiguration, events []*model.Event) error {
	s.initialize(config)
	s.startMessage(config)

events:
	for _, e := range s.incomingEvents {
		fmt.Println(e)
		switch e.Id {
		case model.ClientHasCome:
			_, ok := s.clients[e.ClientName]
			if ok {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "YouShallNotPass",
				}
				fmt.Println(event)
				continue
			}
			if e.Time.Compare(s.configuration.OpeningTime) < 0 {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "NotOpenYet",
				}
				fmt.Println(event)
				continue
			}
			s.clients[e.ClientName] = &model.Client{
				Name:          e.ClientName,
				DeskId:        0,
				IsAlreadyHere: true,
			}
		case model.ClientTookTheTable:
			_, ok := s.clients[e.ClientName]
			if !ok {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "ClientUnknown",
				}
				fmt.Println(event)
				continue
			}
			if s.desk[e.DeskId].IsBusy {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "PlaceIsBusy",
				}
				fmt.Println(event)
				continue
			}
			cl := s.clients[e.ClientName]
			if cl.DeskId > 0 {
				s.desk[cl.DeskId].Free(e.Time)
			}
			s.desk[e.DeskId].Take(e.Time)
			s.clients[e.ClientName].DeskId = e.DeskId
		case model.ClientIsWaiting:
			for _, d := range s.desk {
				if !d.IsBusy {
					event := &model.Event{
						Time:   e.Time,
						Id:     model.EventError,
						ErrMsg: "ICanWaitNoLonger!",
					}
					fmt.Println(event)
					continue events
				}
			}
			if len(s.queue) > s.configuration.DeskCount {
				event := &model.Event{
					Time:       e.Time,
					Id:         model.ClientLeftAtClosing,
					ClientName: e.ClientName,
				}
				delete(s.clients, e.ClientName)
				fmt.Println(event)
				continue
			}
			cl := s.clients[e.ClientName]
			s.queue = append(s.queue, cl)
		case model.ClientLeft:
			_, ok := s.clients[e.ClientName]
			if !ok {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "ClientUnknown",
				}
				fmt.Println(event)
				continue
			}
			id := s.clients[e.ClientName].DeskId
			s.desk[id].Free(e.Time)
			delete(s.clients, e.ClientName)

			if len(s.queue) > 0 {
				cl := s.queue[0]
				s.desk[id].Take(e.Time)
				cl.DeskId = id
				event := &model.Event{
					Time:       e.Time,
					ClientName: cl.Name,
					Id:         model.ClientTookTheTableAfterWaiting,
					DeskId:     id,
				}
				s.queue = s.queue[1:]
				fmt.Println(event)
				continue
			}
		}
	}

	names := make([]string, 0)
	for _, cl := range s.clients {
		s.desk[cl.DeskId].Free(s.configuration.ClosingTime)
		names = append(names, cl.Name)
	}
	sort.Strings(names)
	for _, n := range names {
		event := &model.Event{
			Time:       s.configuration.ClosingTime,
			ClientName: n,
			Id:         model.ClientLeftAtClosing,
		}
		fmt.Println(event)
	}
	fmt.Println(s.configuration.ClosingTime.Format(time.TimeOnly)[:5])

	for i := 0; i < s.configuration.DeskCount; i++ {
		fmt.Println(fmt.Sprintf("%d %d %s", i+1, s.revenue(i+1), s.rentDuration(i+1)))
	}
	return nil
}

func (s *ManagerService) initialize(config *model.ClubConfiguration) {
	s.desks = make(map[int]*model.Desk, 0)
	for i := 0; i < config.DeskCount; i++ {
		s.desks[i+1] = &model.Desk{
			Id:           i + 1,
			IsBusy:       false,
			RentDuration: make([]time.Duration, 0),
		}
	}
	s.clients = make(map[string]*model.Client, 0)
	s.queue = make([]*model.Client, 0)
}

func (s *ManagerService) startMessage(config *model.ClubConfiguration) {
	fmt.Println(config.OpeningTime.Format(time.TimeOnly)[:5])
}

func (s *ManagerService) revenue(deskId int) int {
	total := 0
	for _, dur := range s.desk[deskId].RentDuration {
		h := int(dur.Hours())
		m := int(dur.Minutes()) % 60
		if m > 0 {
			total += (int(h) + 1) * s.configuration.Price
		} else {
			total += int(h) * s.configuration.Price
		}
	}
	return total
}

func (s *ManagerService) rentDuration(deskId int) string {
	var t time.Duration
	for _, dur := range s.desk[deskId].RentDuration {
		t += dur
	}
	return fmt.Sprintf("%02d:%02d", int(t.Hours()), int(t.Minutes())%60)
}
