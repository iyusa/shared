package tool

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// FixedString for fixed string handing
type FixedString struct {
	dict             map[string]string
	offset           int
	source           string
	fixedSource      string
	errorDescription string
}

// NewFixedString create new FixedString struct
func NewFixedString(source string) *FixedString {
	return &FixedString{
		source: source,
		offset: 0,
		dict:   make(map[string]string),
	}
}

// Add get string within width and insert into dict with name = key
func (f *FixedString) Add(key string, width int) (result string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("index out of bounds: %v", r)
		}
	}()

	result = f.source[f.offset : f.offset+width]
	f.dict[key] = result
	f.offset = f.offset + width
	f.fixedSource = f.fixedSource + result

	return
}

// Put directly to dictionary
func (f *FixedString) Put(key, value string) {
	f.dict[key] = value
}

// AddExclude like Add but without added to fixedsource
func (f *FixedString) AddExclude(key string, width int) (result string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("index out of bounds: %v", r)
		}
	}()

	result = f.source[f.offset : f.offset+width]
	f.dict[key] = result
	f.offset = f.offset + width

	return
}

// AddRemaining to fixed source
func (f *FixedString) AddRemaining() {
	if f.offset < len(f.source) {
		f.fixedSource = f.fixedSource + f.source[f.offset:]
	}
}

// Get get value form dictionary
func (f *FixedString) Get(key string) string {
	val, ok := f.dict[key]
	if ok {
		f.errorDescription = ""
		return strings.TrimSpace(val)
	}
	f.errorDescription = fmt.Sprintf("Key %s not found", key)
	return ""
}

// GetNumber with optional last two is decimal
func (f *FixedString) GetNumber(key string, lastTwoDecimal bool) string {
	s := f.Get(key)
	n := len(s)
	if lastTwoDecimal && n > 2 {
		rem := n - 2
		x := s[0:rem]
		return StringToNumberString(x)
	}
	return StringToNumberString(s)
}

// GetInt convert value of key to integer, 0 if fail
func (f *FixedString) GetInt(key string) int {
	v, err := strconv.Atoi(f.Get(key))
	if err != nil {
		f.errorDescription = err.Error()
		return 0
	}
	f.errorDescription = ""
	return v
}

// Error return error if errorDescrition is not empty
func (f *FixedString) Error() error {
	if f.errorDescription == "" {
		return nil
	}
	return errors.New(f.errorDescription)
}

// Map get internal map dict
func (f *FixedString) Map() map[string]string {
	return f.dict
}
