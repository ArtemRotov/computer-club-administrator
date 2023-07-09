package service

import (
	"fmt"
	"github.com/ArtemRotov/computer-club-manager/internal/model"
	"sort"
	"time"
)

type ManagerService struct {
	config  *model.ClubConfiguration
	desks   map[int]*model.Desk
	clients map[string]*model.Client
	queue   []*model.Client
}

func NewManagerService() *ManagerService {
	return &ManagerService{
		config:  nil,
		desks:   nil,
		clients: nil,
		queue:   nil,
	}
}

// Run - handles all incoming events
func (s *ManagerService) Run(config *model.ClubConfiguration, events []*model.Event) error {
	s.initialize(config)

	s.start()
events:
	for _, e := range events {
		if e.Time.Compare(s.config.ClosingTime) != -1 {
			s.stop()
			return nil
		}
		fmt.Println(e)
		switch e.Id {
		case model.ClientHasCome:
			_, ok := s.clients[e.ClientName]
			if ok {
				fmt.Println(s.newErrorEvent(e.Time, eventErrorYouShallNotPass))
				continue
			}
			if e.Time.Compare(s.config.OpeningTime) < 0 {
				fmt.Println(s.newErrorEvent(e.Time, eventErrorNotOpenYet))
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
				fmt.Println(s.newErrorEvent(e.Time, eventErrorClientUnknown))
				continue
			}
			if s.desks[e.DeskId].IsBusy {
				fmt.Println(s.newErrorEvent(e.Time, eventErrorPlaceIsBusy))
				continue
			}
			cl := s.clients[e.ClientName]
			if cl.DeskId > 0 {
				s.desks[cl.DeskId].Free(e.Time)
			}
			s.desks[e.DeskId].Take(e.Time)
			s.clients[e.ClientName].DeskId = e.DeskId

		case model.ClientIsWaiting:
			for _, d := range s.desks {
				if !d.IsBusy {
					fmt.Println(s.newErrorEvent(e.Time, eventErrorCanWaitNoLonger))
					continue events
				}
			}
			if len(s.queue) > s.config.DeskCount {
				fmt.Println(s.newEvent(e.Time, model.ClientLeftAtClosing, e.ClientName))
				delete(s.clients, e.ClientName)
				continue
			}
			cl := s.clients[e.ClientName]
			s.queue = append(s.queue, cl)

		case model.ClientLeft:
			_, ok := s.clients[e.ClientName]
			if !ok {
				fmt.Println(s.newErrorEvent(e.Time, eventErrorClientUnknown))
				continue
			}
			id := s.clients[e.ClientName].DeskId
			s.desks[id].Free(e.Time)
			delete(s.clients, e.ClientName)

			if len(s.queue) > 0 {
				cl := s.queue[0]
				s.desks[id].Take(e.Time)
				cl.DeskId = id
				s.queue = s.queue[1:]
				fmt.Println(s.newEventWithDesk(e.Time, model.ClientTookTheTableAfterWaiting, id, cl.Name))
			}
		}
	}
	s.stop()

	return nil
}

// initialize - initializes service structures
func (s *ManagerService) initialize(config *model.ClubConfiguration) {
	s.config = config
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

// start - performs actions before starting processing
func (s *ManagerService) start() {
	fmt.Println(s.config.OpeningTime.Format(time.TimeOnly)[:5])
}

// stop - end of processing
func (s *ManagerService) stop() {
	names := make([]string, 0)
	for _, cl := range s.clients {
		s.desks[cl.DeskId].Free(s.config.ClosingTime)
		names = append(names, cl.Name)
	}
	sort.Strings(names)

	for _, n := range names {
		fmt.Println(s.newEvent(s.config.ClosingTime, model.ClientLeftAtClosing, n))
	}
	fmt.Println(s.config.ClosingTime.Format(time.TimeOnly)[:5])

	for i := 0; i < s.config.DeskCount; i++ {
		fmt.Println(fmt.Sprintf("%d %d %s", i+1, s.revenue(i+1), s.rentDuration(i+1)))
	}
}

// newErrorEvent - returns a new outgoing error event
func (s *ManagerService) newErrorEvent(t time.Time, msg string) *model.Event {
	return &model.Event{
		Time:   t,
		Id:     model.EventError,
		ErrMsg: msg,
	}
}

// newEvent - returns a new outgoing event with time, id and client name
func (s *ManagerService) newEvent(t time.Time, id int, name string) *model.Event {
	return &model.Event{
		Time:       t,
		Id:         model.EventId(id),
		ClientName: name,
	}
}

// newEventWithDesk - returns a new outgoing event with time, id , deskId and client name
func (s *ManagerService) newEventWithDesk(t time.Time, id, deskId int, name string) *model.Event {
	return &model.Event{
		Time:       t,
		Id:         model.EventId(id),
		ClientName: name,
		DeskId:     deskId,
	}
}

// revenue - returns the revenue of a particular desk
func (s *ManagerService) revenue(deskId int) int {
	total := 0
	for _, dur := range s.desks[deskId].RentDuration {
		h := int(dur.Hours())
		m := int(dur.Minutes()) % 60
		if m > 0 {
			total += (int(h) + 1) * s.config.Price
		} else {
			total += int(h) * s.config.Price
		}
	}
	return total
}

// rentDuration - returns total rent duration of a particular desk
func (s *ManagerService) rentDuration(deskId int) string {
	var t time.Duration
	for _, dur := range s.desks[deskId].RentDuration {
		t += dur
	}
	return fmt.Sprintf("%02d:%02d", int(t.Hours()), int(t.Minutes())%60)
}
