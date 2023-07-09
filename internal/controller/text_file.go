package controller

import (
	"bufio"
	"github.com/ArtemRotov/computer-club-manager/internal/model"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	inputDeskCount    = iota // row with the number of desks
	inputTimeInterval        // row with the time interval
	inputPrice               // row with rent price
	inputEvents              // row with first incoming event
)

type ManagerService interface {
	Run(config *model.ClubConfiguration, events []*model.Event) error
}

type TextFileHandler struct {
	managerService ManagerService
	currentRow     int
}

func NewTextFileHandler(m ManagerService) *TextFileHandler {
	return &TextFileHandler{
		managerService: m,
		currentRow:     inputDeskCount,
	}
}

// Handle - input file handler function.
// Returns an error if the input type is violated.
func (h *TextFileHandler) Handle(file *os.File) error {
	scanner := bufio.NewScanner(file)
	input := make([]string, 0)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	h.currentRow = inputDeskCount

	deskCount, err := h.parseDeskCount(input)
	if err != nil {
		return err
	}

	openTime, closeTime, err := h.parseTimes(input)
	if err != nil {
		return err
	}

	price, err := h.parsePrice(input)
	if err != nil {
		return err
	}

	events, err := h.parseEvents(input)
	if err != nil {
		return err
	}

	return h.managerService.Run(
		&model.ClubConfiguration{
			DeskCount:   deskCount,
			Price:       price,
			OpeningTime: openTime,
			ClosingTime: closeTime,
		},
		events)
}

// parseDeskCount - parses string to create desk count.
// Returns an error if the input type is violated.
func (h *TextFileHandler) parseDeskCount(input []string) (int, error) {
	if len(input) <= h.currentRow {
		return 0, newLineError(h.currentRow+1, nil)
	}

	deskCount, err := strconv.Atoi(input[h.currentRow])
	if err != nil {
		return 0, newLineError(h.currentRow+1, err)
	}

	h.currentRow = inputTimeInterval
	return deskCount, nil
}

// parseTimes - parses strings to create time interval.
// Returns an error if the input type is violated.
func (h *TextFileHandler) parseTimes(input []string) (time.Time, time.Time, error) {
	if len(input) <= h.currentRow {
		return time.Time{}, time.Time{}, newLineError(h.currentRow+1, nil)
	}

	const Pattern = "^[0-9]{2}:[0-9]{2} [0-9]{2}:[0-9]{2}$"
	matched, err := regexp.MatchString(Pattern, input[h.currentRow])
	if err != nil || !matched {
		return time.Time{}, time.Time{}, newLineError(h.currentRow+1, ErrCannotParseTimeValue)
	}

	times := strings.Split(input[h.currentRow], " ")
	// opening time
	openingTime, err := time.Parse(time.TimeOnly, times[0]+":00")
	if err != nil {
		return time.Time{}, time.Time{}, newLineError(h.currentRow+1, ErrCannotParseTimeValue)
	}
	// closing time
	closingTime, err := time.Parse(time.TimeOnly, times[1]+":00")
	if err != nil {
		return time.Time{}, time.Time{}, newLineError(h.currentRow+1, ErrCannotParseTimeValue)
	}
	if openingTime.Compare(closingTime) != -1 {
		return time.Time{}, time.Time{}, newLineError(h.currentRow+1, ErrCannotBadTimeInterval)
	}

	h.currentRow = inputPrice
	return openingTime, closingTime, nil
}

// parsePrice - parses string to create price.
// Returns an error if the input type is violated.
func (h *TextFileHandler) parsePrice(input []string) (int, error) {
	if len(input) <= h.currentRow {
		return 0, newLineError(h.currentRow+1, nil)
	}

	price, err := strconv.Atoi(input[h.currentRow])
	if err != nil {
		return 0, newLineError(h.currentRow+1, err)
	}

	h.currentRow = inputEvents
	return price, nil
}

// parseEvents - parses strings to create events.
// Returns an error if the input type is violated.
func (h *TextFileHandler) parseEvents(input []string) ([]*model.Event, error) {
	events := make([]*model.Event, 0)

	for ; h.currentRow < len(input); h.currentRow++ {
		ev, err := h.newEvent(input[h.currentRow])
		if err != nil {
			return nil, newLineError(h.currentRow+1, err)
		}
		events = append(events, ev)
	}

	return events, nil
}

// newEvent - creates an event model by processing the input string.
// Returns an error if the input type is violated.
func (h *TextFileHandler) newEvent(str string) (*model.Event, error) {
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
