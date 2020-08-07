package errors

import (
	"fmt"
	"strconv"
	"sync"
)

// ErrorList is a list of errors that do not fit into the Wrapper model
// because they are at same level and don't have direct causal relationships.
// For example, a user request that lacks both username and password can have two parallel errors.
//
// The interface is mainly used for error unwrapping, for create error list, use MultiErr
// TODO: a better name than ErrorList
type ErrorList interface {
	Errors() []error
}

// MultiErr is a slice of Error. It has two implementation,
// NewMultiErr return a non thread safe version,
// NewMultiErrSafe return a thread safe version using mutex.
type MultiErr interface {
	error
	fmt.Formatter
	// Append adds error to current selection, it will flatten the error being added if it is also MultiErr
	// It returns true if the appended error is not nil, inspired by https://github.com/uber-go/multierr/issues/21
	Append(error) bool
	// Errors returns errors stored, if no error
	// TODO: multiErr returns internal []error while multiErrSafe returns a copy
	// but even for thread safe error, Errors() normally get called in one go routine.
	Errors() []error
	// ErrorOrNil returns itself or nil if there are no errors, inspired by https://github.com/hashicorp/go-multierror
	ErrorOrNil() error
	// HasError is ErrorOrNil != nil
	HasError() bool
	// Len returns length of underlying errors, it's a shortcut for `len(merr.Errors)`
	Len() int
}

// NewMultiErr returns a non thread safe implementation
// Deprecated: use NewMulti instead, this is an error package and *Err is redundant.
func NewMultiErr() MultiErr {
	return &multiErr{}
}

// NewMulti returns a non thread safe MultiErr implementation.
// Use NewMultiSafe if you need one protected with sync.Mutex.
// Or you can use your own locking for access from multiple goroutines.
func NewMulti() MultiErr {
	return &multiErr{}
}

// NewMultiErrSafe returns a thread safe implementation which protects the underlying slice using mutex.
// It returns a copy of slice when Errors is called
func NewMultiErrSafe() MultiErr {
	return &multiErrSafe{}
}

// NewMultiSafe returns a MultiErr protected with sync.Mutex.
// Use NewMulti if you use multi error in one go routine or have your own locking.
// It returns a copy of slice when Errors is called. TODO(at15): consider change this behavior or update interface doc.
func NewMultiSafe() MultiErr {
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
	if len(m.errs) == 0 {
		return nil
	}
	return m
}

func (m *multiErr) HasError() bool {
	if len(m.errs) == 0 {
		return false
	}
	return true
}

func (m *multiErr) Format(s fmt.State, verb rune) {
	// TODO: support different verb
	switch verb {
	case 'v', 's':
		s.Write([]byte(formatErrors(m.errs)))
	case 'q':
		fmt.Fprintf(s, "%q", formatErrors(m.errs))
	}
}

func (m *multiErr) Len() int {
	return len(m.errs)
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
	// nothing
	if m.errs == nil {
		m.mu.Unlock()
		return nil
	}

	// make a copy if there are content
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
	m.mu.Lock()
	if len(m.errs) == 0 {
		m.mu.Unlock()
		return false
	} else {
		m.mu.Unlock()
		return true
	}
}

func (m *multiErrSafe) Len() int {
	m.mu.Lock()
	l := len(m.errs)
	m.mu.Unlock()
	return l
}

func (m *multiErrSafe) Format(s fmt.State, verb rune) {
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
	// TODO(at15): might use strings.Join implementation, two loops, first calculate the length of the strings to
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
