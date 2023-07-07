package service

import (
	"bufio"
	"fmt"
	"github.com/ArtemRotov/computer-club-administrator/internal/model"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AdministratorService struct {
	configuration  *model.ClubConfiguration
	desk           map[int]model.Desk
	incomingEvents []*model.Event
	clients        map[string]*model.Client
}

// New - parses the file by filling in the configuration and the slice of events.
// Returns an error if the input type is violated.
func New(file *os.File) (*AdministratorService, error) {
	var err error
	scanner := bufio.NewScanner(file)
	input := make([]string, 0)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	currentStr := 0
	srv := &AdministratorService{
		configuration:  &model.ClubConfiguration{},
		desk:           map[int]model.Desk{},
		incomingEvents: []*model.Event{},
		clients:        map[string]*model.Client{},
	}

	// parse desk count
	if len(input) > currentStr {
		srv.configuration.DeskCount, err = strconv.Atoi(input[0])
		if err != nil {
			return nil, newLineError(currentStr+1, err)
		}
		currentStr++
	} else {
		return nil, newLineError(currentStr+1, nil)
	}

	// parse times
	if len(input) > currentStr {
		times := strings.Split(input[1], " ")
		if len(times) != 2 {
			return nil, newLineError(currentStr+1, ErrCannotParseTimeValue)
		}
		// opening time
		srv.configuration.OpeningTime, err = time.Parse(time.TimeOnly, times[0]+":00")
		if err != nil {
			return nil, newLineError(currentStr+1, ErrCannotParseTimeValue)
		}
		// closing
		srv.configuration.ClosingTime, err = time.Parse(time.TimeOnly, times[1]+":00")
		if err != nil {
			return nil, newLineError(currentStr+1, ErrCannotParseTimeValue)
		}
		currentStr++
	} else {
		return nil, newLineError(currentStr+1, nil)
	}

	// parse price
	if len(input) > currentStr {
		srv.configuration.Price, err = strconv.Atoi(input[2])
		if err != nil {
			return nil, newLineError(currentStr+1, err)
		}
		currentStr++
	} else {
		return nil, newLineError(currentStr+1, nil)
	}

	// parse events
	for i := currentStr; i < len(input); i++ {
		ev, err := srv.newEvent(input[i])
		if err != nil {
			return nil, newLineError(currentStr+1, err)
		}
		srv.incomingEvents = append(srv.incomingEvents, ev)
	}

	// create desks
	for i := 0; i < srv.configuration.DeskCount; i++ {
		srv.desk[i+1] = model.Desk{
			Id:     i + 1,
			IsBusy: false,
		}
	}

	return srv, nil
}

// NewEvent - creates an event model by processing the input string.
// Returns an error if the input type is violated.
func (s *AdministratorService) newEvent(str string) (*model.Event, error) {
	const Pattern1 = "^[0-9]{2}:[0-9]{2} [0-9]* [a-z0-9_-]*$"        // like "08:48 1 client1"
	const Pattern2 = "^[0-9]{2}:[0-9]{2} [0-9]* [a-z0-9_-]* [0-9]*$" // like "08:48 1 client1 1"

	matched, err := regexp.MatchString(Pattern1, str)
	if err != nil {
		return nil, err
	}
	if !matched {
		matched, err = regexp.MatchString(Pattern2, str)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, ErrNotMatchEventPattern
		}
	}
	strs := strings.Split(str, " ")
	e := &model.Event{}

	e.Time, err = time.Parse(time.TimeOnly, strs[0]+":00")
	if err != nil {
		return nil, ErrCannotParseTimeValue
	}

	id, err := strconv.Atoi(strs[1])
	if err != nil {
		return nil, ErrCannotParseEventId
	}
	e.Id = model.EventId(id)

	e.ClientName = strs[2]

	if len(strs) == 4 {
		e.DeskId, err = strconv.Atoi(strs[3])
		if err != nil {
			return nil, ErrCannotParseDeskId
		}
	}

	return e, nil
}

func (s *AdministratorService) Run() error {
	fmt.Println(s.configuration.OpeningTime.Format(time.TimeOnly)[:5])

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
			}
			if e.Time.Compare(s.configuration.OpeningTime) < 0 {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "NotOpenYet",
				}
				fmt.Println(event)
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
			}
			if s.desk[e.DeskId].IsBusy {
				event := &model.Event{
					Time:   e.Time,
					Id:     model.EventError,
					ErrMsg: "PlaceIsBusy",
				}
				fmt.Println(event)
			}
		case model.ClientIsWaiting:
		case model.ClientLeft:
		}

	}

	fmt.Println(s.configuration.ClosingTime.Format(time.TimeOnly)[:5])

	for i := 0; i < s.configuration.DeskCount; i++ {
		fmt.Println(fmt.Sprintf("%d %d %s", i+1, s.revenue(i+1), s.rentDuration(i+1)))
	}
	return nil
}

func (s *AdministratorService) revenue(deskId int) int {
	total := 0
	for _, dur := range s.desk[deskId].RentDuration {
		h := dur.Hours()
		m := dur.Minutes()
		if m > 0 {
			total += (int(h) + 1) * s.configuration.Price
		} else {
			total += int(h) * s.configuration.Price
		}
	}
	return total
}

func (s *AdministratorService) rentDuration(deskId int) string {
	var t time.Duration
	for _, dur := range s.desk[deskId].RentDuration {
		t += dur
	}
	return fmt.Sprintf("%02d:%02d", int(t.Hours()), int(t.Minutes()))
}
