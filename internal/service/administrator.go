package service

import (
	"bufio"
	"errors"
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
}

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
	srv := &AdministratorService{}

	// parse desk count
	if len(input) > currentStr {
		srv.configuration.DeskCount, err = strconv.Atoi(input[0])
		if err != nil {
			return nil, fmt.Errorf("error at line %d : cannot parse DeskCount (%w)", currentStr, err)
		}
		currentStr++
	} else {
		return nil, errors.New(fmt.Sprintf("error at line %d", currentStr))
	}

	// parse times
	if len(input) > currentStr {
		times := strings.Split(input[1], " ")
		if len(times) != 2 {
			return nil, fmt.Errorf("error at line %d (%w)", currentStr, ErrCannotParseTimeValue)
		}
		srv.configuration.OpeningTime, err = time.Parse(model.TimeLayout, times[0])
		if err != nil {
			return nil, fmt.Errorf("error at line %d (%w)", currentStr, ErrCannotParseTimeValue)
		}
		srv.configuration.ClosingTime, err = time.Parse(model.TimeLayout, times[1])
		if err != nil {
			return nil, fmt.Errorf("error at line %d (%w)", currentStr, ErrCannotParseTimeValue)
		}
		currentStr++
	} else {
		return nil, errors.New(fmt.Sprintf("error at line %d", currentStr))
	}

	// parse price
	if len(input) > currentStr {
		srv.configuration.Price, err = strconv.Atoi(input[2])
		if err != nil {
			return nil, fmt.Errorf("error at line %d : cannot parse price (%w)", currentStr, err)
		}
		currentStr++
	} else {
		return nil, errors.New(fmt.Sprintf("error at line %d", currentStr))
	}

	return nil, nil
}

func (s *AdministratorService) NewEvent(str string) (*model.Event, error) {
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

	e.Time, err = time.Parse(model.TimeLayout, strs[0])
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
