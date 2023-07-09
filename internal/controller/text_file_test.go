package controller

import (
	"github.com/ArtemRotov/computer-club-manager/internal/model"
	"reflect"
	"testing"
	"time"
)

func TestTextFileHandler_parseDeskCount(t *testing.T) {
	h := &TextFileHandler{
		currentRow: 0,
	}
	tc := []struct {
		name      string
		input     []string
		expected  int
		expectErr bool
	}{
		{
			name:      "parseDeskCount #1",
			input:     []string{""},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "parseDeskCount #2",
			input:     []string{"d"},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "parseDeskCount #3",
			input:     []string{"1"},
			expected:  1,
			expectErr: false,
		},
	}
	for _, tt := range tc {
		h.currentRow = 0
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.parseDeskCount(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseDeskCount() error = %v, expect error %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parseDeskCount() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestTextFileHandler_parseTimes(t *testing.T) {
	h := &TextFileHandler{
		currentRow: 1,
	}

	t1, _ := time.Parse(time.TimeOnly, "09:00:00")
	t2, _ := time.Parse(time.TimeOnly, "19:00:00")

	tc := []struct {
		name      string
		input     []string
		expected1 time.Time
		expected2 time.Time
		expectErr bool
	}{
		{
			name:      "parseTimes #1",
			input:     []string{""},
			expected1: time.Time{},
			expected2: time.Time{},
			expectErr: true,
		},
		{
			name:      "parseTimes #2",
			input:     []string{"sf212"},
			expected1: time.Time{},
			expected2: time.Time{},
			expectErr: true,
		},
		{
			name:      "parseTimes #3 (no zero at begin)",
			input:     []string{"9:00 19:00"},
			expected1: time.Time{},
			expected2: time.Time{},
			expectErr: true,
		},
		{
			name:      "parseTimes #4 (opening time after closing)",
			input:     []string{"20:00 19:00"},
			expected1: time.Time{},
			expected2: time.Time{},
			expectErr: true,
		},
		{
			name:      "parseTimes #5",
			input:     []string{"09:00 19:00"},
			expected1: t1,
			expected2: t2,
			expectErr: false,
		},
	}
	for _, tt := range tc {
		h.currentRow = 0
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := h.parseTimes(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("parseTimes() error = %v, expect error %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected1) {
				t.Errorf("parseTimes() got = %v, expected %v", got, tt.expected1)
			}
			if !reflect.DeepEqual(got1, tt.expected2) {
				t.Errorf("parseTimes() got1 = %v, expected %v", got1, tt.expected2)
			}
		})
	}
}

func TestTextFileHandler_parsePrice(t *testing.T) {
	h := &TextFileHandler{
		currentRow: 0,
	}
	tc := []struct {
		name      string
		input     []string
		expected  int
		expectErr bool
	}{
		{
			name:      "parsePrice #1",
			input:     []string{""},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "parsePrice #2",
			input:     []string{"d"},
			expected:  0,
			expectErr: true,
		},
		{
			name:      "parsePrice #3",
			input:     []string{"1"},
			expected:  1,
			expectErr: false,
		},
	}
	for _, tt := range tc {
		h.currentRow = 0
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.parsePrice(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("parsePrice() error = %v, expect error %v", err, tt.expectErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parsePrice() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestTextFileHandler_newEvent(t *testing.T) {
	h := &TextFileHandler{
		currentRow: 0,
	}
	t1, _ := time.Parse(time.TimeOnly, "09:00:00")

	tc := []struct {
		name      string
		input     string
		expected  *model.Event
		expectErr bool
	}{
		{
			name:      "newEvent #1",
			input:     "",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "newEvent #2",
			input:     "1e1wrw2ed3sf12312",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "newEvent #3",
			input:     "08:48 client1",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "newEvent #4",
			input:     "08:48 client1 1",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "newEvent #5",
			input:     "8:48 1 client1 1",
			expected:  nil,
			expectErr: true,
		},
		{
			name:  "newEvent #6",
			input: "09:00 1 client1",
			expected: &model.Event{
				Time:       t1,
				Id:         1,
				ClientName: "client1",
			},
			expectErr: false,
		},
		{
			name:  "newEvent #7",
			input: "09:00 1 client1 1",
			expected: &model.Event{
				Time:       t1,
				Id:         1,
				ClientName: "client1",
				DeskId:     1,
			},
			expectErr: false,
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			got, err := h.newEvent(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("newEvent() error = %v, expect error %v", err, tt.expectErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("newEvent() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}
