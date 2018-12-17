# 2018-12-14 Reference

This doc go over the references mentioned in [error-categorization](2018-12-10-error-categorization.md)

## Go 2

- notes: https://github.com/at15/papers-i-read/commit/a44b13fd8e1bcc51135ebe128f391ebc674904c4
- source: https://go.googlesource.com/proposal/+/master/design/go2draft-error-values-overview.md
 
TODO

- [ ] it mentioned https://github.com/spacemonkeygo/errors which has class etc.
  - it allows attach key value pairs to error
  - defines common set of error https://github.com/spacemonkeygo/errors/blob/master/errors.go#L563-L600
- [x] https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html by rob pike
  - defines error kinds https://godoc.org/upspin.io/errors#Kind
  - https://godoc.org/upspin.io/errors#MarshalError to transfer error across the wire

https://go.googlesource.com/proposal/+/master/design/go2draft-error-printing.md

- error printing is only for read by human
- it print trace of error (w/ or w/o stack?)
  - [ ] now in gommon/errors we only do wrapping in first error, however the call stack of init error should be different from the wrapping stack
- mentioned list of error (which is multi error)

https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md

- `Is` is just same as upspin, and I think it only works for sentinel errors like `io.EOF`

````go
func Is(err, target error) bool {
	for {
		if err == target {
			return true
		}
		wrapper, ok := err.(Wrapper)
		if !ok {
			return false
		}
		err = wrapper.Unwrap()
		if err == nil {
			return false
		}
	}
}
````

- `As` requires contracts ... (pass type as a parameter)
  - without polymorphism
  - [ ] I didn't quite get this part ..
 
````go
func As(type E)(err error) (e E, ok bool) {
	for {
		if e, ok := err.(E); ok {
			return e, true
		}
		wrapper, ok := err.(Wrapper)
		if !ok {
			return e, false
		}
		err = wrapper.Unwrap()
		if err == nil {
			return e, false
		}
	}
}
````

````go
// instead of pe, ok := err.(*os.PathError)
var pe *os.PathError
if errors.AsValue(&pe, err) { ... pe.Path ... }
````

- can use reflect to implement this, to compare type, just use `reflect.TypeOf(err).String`

## Upspin

https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html

- didn't use stack trace for error, show operational trace

> There is a tension between making errors helpful and concise for the end user versus making them expansive and analytic for the implementer. 
Too often the implementer wins and the errors are overly verbose, to the point of including stack traces or other overwhelming detail

> Upspin's errors are an attempt to serve both the users and the implementers. 
The reported errors are reasonably concise, concentrating on information the user should find helpful. 
But they also contain internal details such as method names an implementer might find diagnostic but not in a way that overwhelms the user. 
In practice we find that the tradeoff has worked well

> In contrast, a stack trace-like error is worse in both respects. 
The user does not have the context to understand the stack trace, 
and an implementer shown a stack trace is denied the information that could be presented 
if the server-side error was passed to the client. This is why Upspin error nesting behaves as an operational trace, 
showing the path through the elements of the system, rather than as an execution trace, showing the path through the code. 
The distinction is vital

- `func Is(kind Kind, err error) bool` I think is almost same as the go 2 proposal `Is`, compare through the cause chain
- `func Match(template, err error) bool` is almost same as go 2 proposal `func As(type E)(err error) (e E, ok bool) {}`

````go
// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}

// Match only operates on upspin's Error type, so it does not need the polymorphism like go 2
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	// un wrap and compare etc.
}
````

> Errors are for users, not just for programmers

## errorwrap

- https://github.com/hashicorp/errwrap/blob/master/errwrap.go
- do unwrapping 
- allow match using string, exact match on error string, not strings.Contains
- allow get a slice of errors from error chain after matching `func GetAll(err error, msg string) []error `
- match type use reflect `func GetAllType(err error, v interface{}) []error`

````go
// GetAllType gets all the errors that are the same type as v.
//
// The order of the return value is the same as described in GetAll.
func GetAllType(err error, v interface{}) []error {
	var result []error

	var search string
	if v != nil {
		search = reflect.TypeOf(v).String()
	}
	Walk(err, func(err error) {
		var needle string
		if err != nil {
			needle = reflect.TypeOf(err).String()
		}

		if needle == search {
			result = append(result, err)
		}
	})

	return result
}

// Walk walks all the wrapped errors in err and calls the callback. If
// err isn't a wrapped error, this will be called once for err. If err
// is a wrapped error, the callback will be called for both the wrapper
// that implements error as well as the wrapped error itself.
func Walk(err error, cb WalkFunc) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *wrappedError:
		cb(e.Outer)
		Walk(e.Inner, cb)
	case Wrapper:
		cb(err)

		for _, err := range e.WrappedErrors() {
			Walk(err, cb)
		}
	default:
		cb(err)
	}
}

// wrappedError is an implementation of error that has both the
// outer and inner errors.
type wrappedError struct {
	Outer error
	Inner error
}

func (w *wrappedError) Error() string {
	return w.Outer.Error()
}

func (w *wrappedError) WrappedErrors() []error {
	return []error{w.Outer, w.Inner}
}
````

## TiDB

- https://github.com/pingcap/tidb/blob/master/terror/terror.go
- https://github.com/pingcap/errors fork of pkg/errors with juju adapter
- `ErroClass` is pretty application specific (common source of database errors)
- can convert to MySQL error code

````go
// Error implements error interface and adds integer Class and Code, so
// errors with different message can be compared.
type Error struct {
	class   ErrClass
	code    ErrCode
	message string
	args    []interface{}
	file    string
	line    int
}

// Class returns ErrClass
func (e *Error) Class() ErrClass {
	return e.class
}

// Code returns ErrCode
func (e *Error) Code() ErrCode {
	return e.code
}

// ErrCode represents a specific error type in a error class.
// Same error code can be used in different error classes.
type ErrCode int

// ErrClass represents a class of errors.
type ErrClass int

// Error classes.
const (
	ClassAutoid ErrClass = iota + 1
	ClassDDL
	ClassDomain
	ClassEvaluator
	ClassExecutor
	ClassExpression
	ClassAdmin
	ClassKV
	ClassMeta
	ClassOptimizer
	ClassParser
	ClassPerfSchema
	ClassPrivilege
	ClassSchema
	ClassServer
	ClassStructure
	ClassVariable
	ClassXEval
	ClassTable
	ClassTypes
	ClassGlobal
	ClassMockTikv
	ClassJSON
	ClassTiKV
	ClassSession
	// Add more as needed.
)
````

## gRPC

- https://github.com/grpc/grpc-go/blob/master/codes/codes.go
- https://github.com/grpc/grpc-go/blob/master/status/status.go
  - `Code` return code defined in codes
  - status can contains detail, which are message encoded into string
- https://github.com/grpc-ecosystem/grpc-opentracing/blob/master/go/otgrpc/errors.go

````go
// A Class is a set of types of outcomes (including errors) that will often
// be handled in the same way.
type Class string

const (
	Unknown Class = "0xx"
	// Success represents outcomes that achieved the desired results.
	Success Class = "2xx"
	// ClientError represents errors that were the client's fault.
	ClientError Class = "4xx"
	// ServerError represents errors that were the server's fault.
	ServerError Class = "5xx"
)

// A Code is an unsigned 32-bit error code as defined in the gRPC spec.
type Code uint32

const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled (typically by the caller).
	Canceled Code = 1

	// Unknown error. An example of where this error may be returned is
	// if a Status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown Code = 2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument Code = 3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded Code = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound Code = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied Code = 7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Code = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	//
	// A litmus test that may help a service implementor in deciding
	// between FailedPrecondition, Aborted, and Unavailable:
	//  (a) Use Unavailable if the client can retry just the failing call.
	//  (b) Use Aborted if the client should retry at a higher-level
	//      (e.g., restarting a read-modify-write sequence).
	//  (c) Use FailedPrecondition if the client should not retry until
	//      the system state has been explicitly fixed. E.g., if an "rmdir"
	//      fails because the directory is non-empty, FailedPrecondition
	//      should be returned since the client should not retry unless
	//      they have first fixed up the directory by deleting files from it.
	//  (d) Use FailedPrecondition if the client performs conditional
	//      REST Get/Update/Delete on a resource and the resource on the
	//      server does not match the condition. E.g., conflicting
	//      read-modify-write on the same resource.
	FailedPrecondition Code = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted Code = 10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes. For example, a 32-bit file
	// system will generate InvalidArgument if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate
	// OutOfRange if asked to read from an offset past the current
	// file size.
	//
	// There is a fair bit of overlap between FailedPrecondition and
	// OutOfRange. We recommend using OutOfRange (the more specific
	// error) when it applies so that callers who are iterating through
	// a space can easily look for an OutOfRange error to detect when
	// they are done.
	OutOfRange Code = 11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Code = 12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal Code = 13

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable Code = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated Code = 16

	_maxCode = 17
)
````

## docker

https://github.com/moby/moby/blob/master/errdefs/doc.go 

`Errors that cross the package boundary should implement one (and only one) of these interfaces.`

- `defs.go` defines interfaces like `ErrNotFound`
- `helpers.go` have util method to create a wrapper, with wrap the error without any extra message, just to indicate type
  - this is actually very useful, sometimes I just want to say what type of error it is without adding extra message

````go
// helpers.go
type errNotFound struct{ error }

func (errNotFound) NotFound() {}

func (e errNotFound) Cause() error {
	return e.error
}

// NotFound is a helper to create an error of the class with the same name from any error type
func NotFound(err error) error {
	if err == nil {
		return nil
	}
	return errNotFound{err}
}

// is.go
type causer interface {
	Cause() error
}

func getImplementer(err error) error {
	switch e := err.(type) {
	case
		ErrNotFound,
		ErrInvalidParameter,
		ErrConflict,
		ErrUnauthorized,
		ErrUnavailable,
		ErrForbidden,
		ErrSystem,
		ErrNotModified,
		ErrAlreadyExists,
		ErrNotImplemented,
		ErrCancelled,
		ErrDeadline,
		ErrDataLoss,
		ErrUnknown:
		return err
	case causer:
		return getImplementer(e.Cause())
	default:
		return err
	}
}

// IsNotFound returns if the passed in error is an ErrNotFound
func IsNotFound(err error) bool {
	_, ok := getImplementer(err).(ErrNotFound)
	return ok
}

// defs.go

// ErrNotFound signals that the requested object doesn't exist
type ErrNotFound interface {
	NotFound()
}

// ErrInvalidParameter signals that the user input is invalid
type ErrInvalidParameter interface {
	InvalidParameter()
}

// ErrConflict signals that some internal state conflicts with the requested action and can't be performed.
// A change in state should be able to clear this error.
type ErrConflict interface {
	Conflict()
}

// ErrUnauthorized is used to signify that the user is not authorized to perform a specific action
type ErrUnauthorized interface {
	Unauthorized()
}

// ErrUnavailable signals that the requested action/subsystem is not available.
type ErrUnavailable interface {
	Unavailable()
}

// ErrForbidden signals that the requested action cannot be performed under any circumstances.
// When a ErrForbidden is returned, the caller should never retry the action.
type ErrForbidden interface {
	Forbidden()
}

// ErrSystem signals that some internal error occurred.
// An example of this would be a failed mount request.
type ErrSystem interface {
	System()
}

// ErrNotModified signals that an action can't be performed because it's already in the desired state
type ErrNotModified interface {
	NotModified()
}

// ErrAlreadyExists is a special case of ErrConflict which signals that the desired object already exists
type ErrAlreadyExists interface {
	AlreadyExists()
}

// ErrNotImplemented signals that the requested action/feature is not implemented on the system as configured.
type ErrNotImplemented interface {
	NotImplemented()
}

// ErrUnknown signals that the kind of error that occurred is not known.
type ErrUnknown interface {
	Unknown()
}

// ErrCancelled signals that the action was cancelled.
type ErrCancelled interface {
	Cancelled()
}

// ErrDeadline signals that the deadline was reached before the action completed.
type ErrDeadline interface {
	DeadlineExceeded()
}

// ErrDataLoss indicates that data was lost or there is data corruption.
type ErrDataLoss interface {
	DataLoss()
}
````

## Teleport

https://github.com/gravitational/trace/blob/master/errors.go

- interface is defined in function (so user can't use it) ...
- has a convert system error

````go
// ConvertSystemError converts system error to appropriate trace error
// if it is possible, otherwise, returns original error
func ConvertSystemError(err error) error {
	innerError := Unwrap(err)

	if os.IsExist(innerError) {
		return WrapWithMessage(&AlreadyExistsError{Message: innerError.Error()}, innerError.Error())
	}
	if os.IsNotExist(innerError) {
		return WrapWithMessage(&NotFoundError{Message: innerError.Error()}, innerError.Error())
	}
	if os.IsPermission(innerError) {
		return WrapWithMessage(&AccessDeniedError{Message: innerError.Error()}, innerError.Error())
	}
	switch realErr := innerError.(type) {
	case *net.OpError:
		return WrapWithMessage(&ConnectionProblemError{
			Message: realErr.Error(),
			Err:     realErr}, realErr.Error())
	case *os.PathError:
		message := fmt.Sprintf("failed to execute command %v error:  %v", realErr.Path, realErr.Err)
		return WrapWithMessage(&AccessDeniedError{
			Message: message,
		}, message)
	case x509.SystemRootsError, x509.UnknownAuthorityError:
		return wrapWithDepth(&TrustError{Err: innerError}, 2)
	}
	if _, ok := innerError.(net.Error); ok {
		return WrapWithMessage(&ConnectionProblemError{
			Message: innerError.Error(),
			Err:     innerError}, innerError.Error())
	}
	return err
}
````

````go
// NotImplemented returns a new instance of NotImplementedError
func NotImplemented(message string, args ...interface{}) error {
	return WrapWithMessage(&NotImplementedError{
		Message: fmt.Sprintf(message, args...),
	}, message, args...)
}

// NotImplementedError defines an error condition to describe the result
// of a call to an unimplemented API
type NotImplementedError struct {
	Message string `json:"message"`
}

// Error returns log friendly description of an error
func (e *NotImplementedError) Error() string {
	return e.Message
}

// OrigError returns original error
func (e *NotImplementedError) OrigError() error {
	return e
}

// IsNotImplementedError indicates that this error is of NotImplementedError type
func (e *NotImplementedError) IsNotImplementedError() bool {
	return true
}

// IsNotImplemented returns whether this error is of NotImplementedError type
func IsNotImplemented(e error) bool {
	type ni interface {
		IsNotImplementedError() bool
	}
	err, ok := Unwrap(e).(ni)
	return ok && err.IsNotImplementedError()
}
````

## Dgraph

https://github.com/dgraph-io/dgraph/blob/master/x/error.go

- they say there are moving to x.Trace, which is using opencensus

````go
// Check logs fatal if err != nil.
func Check(err error) {
	if err != nil {
		log.Fatalf("%+v", Wrap(err))
	}
}

// Check2 acts as convenience wrapper around Check, using the 2nd argument as error.
func Check2(_ interface{}, err error) {
	Check(err)
}

// Ignore function is used to ignore errors deliberately, while keeping the
// linter happy.
func Ignore(_ error) {
	// Do nothing.
}
````

## Sentry

- https://github.com/getsentry/raven-go
- https://docs.sentry.io/clients/go/
- `raven.CaptureErrorAndWait(err, nil)` block call, send error and then exit
- `raven.CaptureError(err, map[string]string{"browser": "Firefox"})`
- use `map[string]interface{}` as meta data
- serialize to json and send to server

````go
type causer interface {
	Cause() error
}

type errWrappedWithExtra struct {
	err       error
	extraInfo map[string]interface{}
}
````

- https://github.com/getsentry/raven-go/blob/master/stacktrace.go extract package etc.

````go
// https://docs.getsentry.com/hosted/clientdev/interfaces/#failure-interfaces
type Stacktrace struct {
	// Required
	Frames []*StacktraceFrame `json:"frames"`
}

type StacktraceFrame struct {
	// At least one required
	Filename string `json:"filename,omitempty"`
	Function string `json:"function,omitempty"`
	Module   string `json:"module,omitempty"`

	// Optional
	Lineno       int      `json:"lineno,omitempty"`
	Colno        int      `json:"colno,omitempty"`
	AbsolutePath string   `json:"abs_path,omitempty"`
	ContextLine  string   `json:"context_line,omitempty"`
	PreContext   []string `json:"pre_context,omitempty"`
	PostContext  []string `json:"post_context,omitempty"`
	InApp        bool     `json:"in_app"`
}

````

- https://github.com/getsentry/raven-go/blob/master/client.go

````go
// https://docs.getsentry.com/hosted/clientdev/#building-the-json-packet
type Packet struct {
	// Required
	Message string `json:"message"`

	// Required, set automatically by Client.Send/Report via Packet.Init if blank
	EventID   string    `json:"event_id"`
	Project   string    `json:"project"`
	Timestamp Timestamp `json:"timestamp"`
	Level     Severity  `json:"level"`
	Logger    string    `json:"logger"`

	// Optional
	Platform    string            `json:"platform,omitempty"`
	Culprit     string            `json:"culprit,omitempty"`
	ServerName  string            `json:"server_name,omitempty"`
	Release     string            `json:"release,omitempty"`
	Environment string            `json:"environment,omitempty"`
	Tags        Tags              `json:"tags,omitempty"`
	Modules     map[string]string `json:"modules,omitempty"`
	Fingerprint []string          `json:"fingerprint,omitempty"`
	Extra       Extra             `json:"extra,omitempty"`

	Interfaces []Interface `json:"-"`
}

// CaptureErrors formats and delivers an error to the Sentry server.
// Adds a stacktrace to the packet, excluding the call to this method.
func (client *Client) CaptureError(err error, tags map[string]string, interfaces ...Interface) string {
	if client == nil {
		return ""
	}

	if err == nil {
		return ""
	}

	if client.shouldExcludeErr(err.Error()) {
		return ""
	}

	extra := extractExtra(err)
	cause := pkgErrors.Cause(err)

	packet := NewPacketWithExtra(err.Error(), extra, append(append(interfaces, client.context.interfaces()...), NewException(cause, GetOrNewStacktrace(cause, 1, 3, client.includePaths)))...)
	eventID, _ := client.Capture(packet, tags)

	return eventID
}

````