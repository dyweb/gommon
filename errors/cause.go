package errors

// cause.go contains interface and unwrap func for the deprecated causer interface we adopted from pkg/errors

// causer returns the underlying error, a error without cause should return itself.
// It is based on the private `causer` interface in pkg/errors, so errors wrapped using pkg/errors can also be handled
// Deprecated: Use Wrapper interface instead,
type causer interface {
	Cause() error
}

// Cause returns root cause of the error (if any), it stops at the last error that does not implement causer interface.
// If you want get direct cause, use DirectCause.
// If error is nil, it will return nil. If error is not wrapped it will return the error itself.
// error wrapped using https://github.com/pkg/errors also satisfies this interface and can be unwrapped as well.
// TODO: might consider rename it to Unwrap since we are deprecating causer interface
func Cause(err error) error {
	if err == nil {
		return nil
	}
	for err != nil {
		switch err.(type) {
		case Wrapper:
			err = err.(Wrapper).Unwrap()
		case causer:
			err = err.(causer).Cause()
		default:
			return err
		}
	}
	return err
}

// DirectCause returns the direct cause of the error (if any).
// It does NOT follow the cause chain all the way down, just the first one (if any),
// If you want to get root cause, use Cause
func DirectCause(err error) error {
	if err == nil {
		return nil
	}
	switch err.(type) {
	case Wrapper:
		return err.(Wrapper).Unwrap()
	case causer:
		return err.(causer).Cause()
	default:
		return err
	}
}
