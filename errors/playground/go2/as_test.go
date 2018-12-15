package go2

import (
	"encoding/json"
	"log"
	"reflect"
	"testing"
)

// test go2 As implementation in go1 (without polymorphism)
//
// From https://go.googlesource.com/proposal/+/master/design/go2draft-error-inspection.md
//
// If Go 2 does not choose to adopt polymorphism or if we need a function to use in the interim, we could write a temporary helper:
// // instead of pe, ok := err.(*os.PathError)
// var pe *os.PathError
// if errors.AsValue(&pe, err) { ... pe.Path ... }

func AsValue(v interface{}, err error) {
	// TODO: I think I need to use reflect

	// TODO: it does not compare package import path? ...
	// type of v **go2.PathError err *go2.PathError
	log.Printf("type of v %s err %s", reflect.TypeOf(v), reflect.TypeOf(err))

	// NOTE: val must be a pointer so we can unmarshal to it
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		// do nothing, panic so the test is easier
		panic("v must be a pointer")
		return
	}
	iv := reflect.Indirect(rv)
	log.Printf("iv is nil? %t", iv.IsNil())
	log.Printf("type of iv %s", reflect.TypeOf(iv))
	ivi := iv.Interface()
	log.Printf("type of ivi %s", reflect.TypeOf(ivi))
	// TODO: https://github.com/hashicorp/errwrap/blob/master/errwrap.go#L108-L130
	//rt := reflect.TypeOf(v)
	//rv.Set(reflect.ValueOf(err))
	iv.Set(reflect.ValueOf(err))
}

// copied from os.PathError
// PathError records an error and the operation and file path that caused it.
type PathError struct {
	Op   string
	Path string
	Err  error
}

func (e *PathError) Error() string { return e.Op + " " + e.Path + ": " + e.Err.Error() }

// Timeout reports whether this error represents a timeout.
//func (e *PathError) Timeout() bool {
//	t, ok := e.Err.(timeout)
//	return ok && t.Timeout()
//}

func TestAsValue(t *testing.T) {
	e1 := &PathError{Op: "open", Path: "/dev/null"}
	e2 := &PathError{}
	t.Logf("e2 Path %s", e2.Path)
	AsValue(&e2, e1)
	t.Logf("e2 Path %s", e2.Path)
}

// TODO: can I use reflect type as first parameter to just do the match

func TestJson(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		b := []byte(`{"a":1}`)
		type s struct {
			A int
		}
		var s1 s
		err := json.Unmarshal(b, &s1)
		t.Logf("s.A %d err %v", s1.A, err)
	})
	// you can decode single number as json
	t.Run("single number", func(t *testing.T) {
		var a int
		err := json.Unmarshal([]byte("1"), &a)
		t.Logf("a %d err %v", a, err)
	})
}
