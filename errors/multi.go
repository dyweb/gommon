package errors

import (
	"strconv"
	"sync"
)

type MultiErr interface {
	error
	// Append returns true if the appended error is not nil, inspired by https://github.com/uber-go/multierr/issues/21
	Append(error) bool
	Errors() []error
	// Error returns itself or nil if there are no errors, inspired by https://github.com/hashicorp/go-multierror
	ErrorOrNil() error
	// HasError is ErrorOrNil != nil
	HasError() bool
}

func NewMultiErr() MultiErr {
	return &multiErr{}
}

func NewMultiErrSafe() MultiErr {
	return &multiErrSafe{}
}

var _ MultiErr = (*multiErr)(nil)

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

var _ MultiErr = (*multiErrSafe)(nil)

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

func formatErrors(errs []error) string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	buf := make([]byte, 0, len(errs)*10)
	buf = strconv.AppendInt(buf, int64(len(errs)), 10)
	buf = append(buf, " errors; "...)
	for i := range errs {
		buf = append(buf, errs[i].Error()...)
		buf = append(buf, MultiErrSep...)
	}
	return string(buf[:len(buf)-2])
}
