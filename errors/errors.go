package errors

type TracedError interface {
	ErrorStack() *Stack
}

var _ error = (*FreshError)(nil)
var _ TracedError = (*FreshError)(nil)

type FreshError struct {
	msg   string
	stack *Stack
}

func (fresh *FreshError) Error() string {
	return fresh.msg
}

func (fresh *FreshError) ErrorStack() *Stack {
	return fresh.stack
}

var _ error = (*WrappedError)(nil)
var _ TracedError = (*WrappedError)(nil)

type WrappedError struct {
	msg   string
	cause error
	stack *Stack
}

func Wrap(err error, msg string) error {
	// NOTE: sometimes we call wrap without check if the error is nil, it is cleaner if it is the last statement in func
	//
	// i.e. return errors.Wrap(f.Close(), "failed to close file")
	//
	// 		if err := f.Close(); err != nil {
	//			return errors.Wrap(err, "failed to close file")
	//      }
	//      return nil
	if err == nil {
		return nil
	}
	var stack *Stack
	if t, ok := err.(TracedError); ok {
		stack = t.ErrorStack()
	} else {
		stack = callers()
	}
	return &WrappedError{
		msg:   msg,
		cause: err,
		stack: stack,
	}
}

func (wrapped *WrappedError) Error() string {
	return wrapped.msg + ": " + wrapped.cause.Error()
}

func (wrapped *WrappedError) ErrorStack() *Stack {
	return wrapped.stack
}
