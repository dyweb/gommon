package errors

import (
	"fmt"
	"strconv"
	"sync"
)

// TODO: add interface for error list so we can do flatten when unwrap, Errors.Errors is not good, or maybe just an not exported interface

// MultiErr is a slice of Error. It has two implementation, NewMultiErr return a non thread safe version,
// NewMultiErrSafe return a thread safe version using mutex
type MultiErr interface {
	error
	fmt.Formatter
	// Append adds error to current selection, it will flatten the error being added if it is also MultiErr
	// It returns true if the appended error is not nil, inspired by https://github.com/uber-go/multierr/issues/21
	Append(error) bool
	// Errors returns errors stored, if no error
	Errors() []error
	// ErrorOrNil returns itself or nil if there are no errors, inspired by https://github.com/hashicorp/go-multierror
	ErrorOrNil() error
	// HasError is ErrorOrNil != nil
	HasError() bool
}

// NewMultiErr returns a non thread safe implementation
func NewMultiErr() MultiErr {
	return &multiErr{}
}

// NewMultiErrSafe returns a thread safe implementation which protects the underlying slice using mutex.
// It returns a copy of slice when Errors is called
func NewMultiErrSafe() MultiErr {
	return &multiErrSafe{}
}

var (
	_ MultiErr = (*multiErr)(nil)
	_ MultiErr = (*multiErrSafe)(nil)
)

// multiErr is NOT thread (goroutine) safe
type multiErr struct {
	errs []error
}

func (m *multiErr) Append(err error) bool {
	if err == nil {
		return false
	}
	if mErr, ok := err.(MultiErr); ok {
		m.errs = append(m.errs, mErr.Errors()...)
	} else {
		m.errs = append(m.errs, err)
	}
	return true
}

func (m *multiErr) Errors() []error {
	return m.errs
}

func (m *multiErr) Error() string {
	return formatErrors(m.errs)
}

func (m *multiErr) ErrorOrNil() error {
	if m == nil || len(m.errs) == 0 {
		return nil
	}
	return m
}

func (m *multiErr) HasError() bool {
	if m == nil || len(m.errs) == 0 {
		return false
	}
	return true
}

func (m *multiErr) Format(s fmt.State, verb rune) {
	// TODO:
	switch verb {
	case 'v', 's':
		s.Write([]byte(formatErrors(m.errs)))
	case 'q':
		fmt.Fprintf(s, "%q", formatErrors(m.errs))
	}
}

// multiErrSafe is thread safe
type multiErrSafe struct {
	mu   sync.Mutex
	errs []error
}

func (m *multiErrSafe) Append(err error) bool {
	if err == nil {
		return false
	}
	m.mu.Lock()
	if mErr, ok := err.(MultiErr); ok {
		m.errs = append(m.errs, mErr.Errors()...)
	} else {
		m.errs = append(m.errs, err)
	}
	m.mu.Unlock()
	return true
}

func (m *multiErrSafe) Errors() []error {
	m.mu.Lock()
	t := make([]error, len(m.errs))
	copy(t, m.errs)
	m.mu.Unlock()
	return t
}

func (m *multiErrSafe) Error() string {
	m.mu.Lock()
	s := formatErrors(m.errs)
	m.mu.Unlock()
	return s
}

func (m *multiErrSafe) ErrorOrNil() error {
	if m == nil {
		return nil
	}
	m.mu.Lock()
	if len(m.errs) == 0 {
		m.mu.Unlock()
		return nil
	} else {
		m.mu.Unlock()
		return m
	}
}

func (m *multiErrSafe) HasError() bool {
	if m == nil {
		return false
	}
	m.mu.Lock()
	if len(m.errs) == 0 {
		m.mu.Unlock()
		return false
	} else {
		m.mu.Unlock()
		return true
	}
}

func (m *multiErrSafe) Format(s fmt.State, verb rune) {
	if m == nil {
		return
	}
	m.mu.Lock()
	switch verb {
	case 'v', 's':
		s.Write([]byte(formatErrors(m.errs)))
	case 'q':
		fmt.Fprintf(s, "%q", formatErrors(m.errs))
	}
	m.mu.Unlock()
}

// TODO: it should support format flag so we can pass it down to sub errors
func formatErrors(errs []error) string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	// TODO: (at15) might use strings.Join implementation, two loops, first calculate the length of the strings to
	// get total buf size, then copy the first element, the second loop copy separator and the n-1 elements, and there is
	// no need to trim the last separator
	buf := make([]byte, 0, len(errs)*10)
	buf = strconv.AppendInt(buf, int64(len(errs)), 10)
	buf = append(buf, " errors; "...)
	for i := range errs {
		buf = append(buf, errs[i].Error()...)
		buf = append(buf, MultiErrSep...)
	}
	return string(buf[:len(buf)-2]) // -2 trim len(MultiErrSep)
}
