package tool

import (
	"bytes"
	"strings"
)

// ReportString create text receipt with maximum length
type ReportString struct {
	buffer         bytes.Buffer
	max            int
	rightAlignment bool
}

// NewReportString create new reportstring with initialized max
func NewReportString(rightAlignment bool) *ReportString {
	return &ReportString{max: 32, rightAlignment: rightAlignment}
}

// Add string with max
func (r *ReportString) Add(value string) {
	length := len(value)

	if length == 0 {
		chars := strings.Repeat(" ", r.max)
		r.buffer.WriteString(chars)
	} else if length == 1 {
		chars := strings.Repeat(value, r.max)
		r.buffer.WriteString(chars)
	} else if length == r.max {
		r.buffer.WriteString(value)
	} else if length > r.max {
		chars := value[0:r.max]
		r.buffer.WriteString(chars)
	} else {
		chars := PadRight(value, " ", r.max)
		r.buffer.WriteString(chars)
	}
}

// AddKV add key value pair
func (r *ReportString) AddKV(key string, value string) {
	n := r.max - len(key)
	chars := PadLeft(value, " ", n)
	if r.rightAlignment {
		chars = PadRight(value, " ", n)
	}
	r.Add(key + chars)
	// r.Add(key + PadLeft(value, " ", n))
}

// AddCenter ed string
func (r *ReportString) AddCenter(value string) {
	length := len(value)

	if length >= r.max {
		r.buffer.WriteString(value[0:r.max])
	} else {
		remaining := r.max - length
		start := remaining / 2
		r.Add(strings.Repeat(" ", start) + value)
	}
}

// String return string representation
func (r *ReportString) String() string {
	return r.buffer.String()
}

/*
// PadRight return empty string on the right
func PadRight(value string, maxLength int) string {
	if len(value) >= maxLength {
		return value[0:maxLength]
	}
	remaining := maxLength - len(value)
	return value + strings.Repeat(" ", remaining)
}

// PadLeft return empty string on the left
func PadLeft(value string, maxLength int) string {
	if len(value) >= maxLength {
		return value[0:maxLength]
	}
	remaining := maxLength - len(value)
	return strings.Repeat(" ", remaining) + value
}
*/
