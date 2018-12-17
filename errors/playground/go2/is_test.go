package go2

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

// test go2 Is without un wrapping logic
//
// From https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md

// check if same error
func Is(err, target error) bool {
	if err == target {
		return true
	}
	return false
}

// check if same error type, all string based error are same type (created from errors.New)
func IsType(err, target error) bool {
	// found from hashicorp/errwrap
	if reflect.TypeOf(err).String() == reflect.TypeOf(target).String() {
		return true
	}
	return false
}

// return error if it is of same type, this is quite useless because we don't have unwrap logic
func GetType(err, target error) (error, bool) {
	if reflect.TypeOf(err).String() == reflect.TypeOf(target).String() {
		return err, true
	}
	return nil, false
}

func giveMeEOF() error {
	return io.EOF
}

func TestIs(t *testing.T) {
	t.Logf("%t", Is(fmt.Errorf("EOF"), io.EOF)) // false
	t.Logf("%t", Is(errors.New("EOF"), io.EOF)) // false
	t.Logf("%t", Is(giveMeEOF(), io.EOF))       // true
}

var ErrTypeString = errors.New("I am just a string")
var ErrTypePathError = &os.PathError{}

var ErrTypeFoo = errFoo{msg: "I am just a template"}

type errFoo struct {
	msg string
}

func (e errFoo) Error() string {
	return e.msg
}

func TestIsType(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		t.Logf("%t", IsType(errors.New("miao"), errors.New("wang"))) // true
		t.Logf("%t", IsType(fmt.Errorf("miao"), errors.New("wang"))) // true
		t.Logf("%t", IsType(fmt.Errorf("miao"), ErrTypeString))      // true
		t.Logf("%t", IsType(os.ErrNotExist, ErrTypeString))          // true
	})

	t.Run("struct pointer", func(t *testing.T) {
		_, err := os.Open("404")                    //os.PathError{}
		t.Logf("%t", IsType(err, ErrTypeString))    // false
		t.Logf("%t", IsType(err, &os.PathError{}))  // true
		t.Logf("%t", IsType(err, ErrTypePathError)) // true
	})

	t.Run("struct", func(t *testing.T) {
		t.Logf("%t", IsType(errFoo{"miao"}, ErrTypeFoo)) // true
	})

	t.Run("example", func(t *testing.T) {
		var err error
		err = errFoo{msg: "a"}
		// NOTE: it can be converted to a switch case as well
		if IsType(err, ErrTypeFoo) {
			fooErr := err.(errFoo)
			t.Logf("%s", fooErr.msg)
		}
	})
}

func TestGetType(t *testing.T) {
	e, ok := GetType(errors.New("miao"), ErrTypeString)
	t.Logf("%v %t", e, ok) // true
	e, ok = GetType(errors.New("miao"), ErrTypePathError)
	t.Logf("%v %t", e, ok)   // false
	_, err := os.Open("404") //os.PathError{}
	e, ok = GetType(err, ErrTypeString)
	t.Logf("%v %t", e, ok) // false
	e, ok = GetType(err, ErrTypePathError)
	t.Logf("%v %t", e, ok) // true
}
