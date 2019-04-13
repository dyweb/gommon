package errors

import (
	"fmt"
	"os"
)

// FIXME: it seems example code need to be in a different package?
// https://github.com/dyweb/gommon/issues/107
// https://github.com/alecthomas/gometalinter/issues/218
// but I never saw it in previous commit

func ExampleWrap() {
	err := Wrap(os.ErrNotExist, "oops")
	fmt.Println(err)
	// Output:
	// oops: file does not exist
}

func ExampleMultiErr() {
	// TODO: demo the return value of append
	err := NewMultiErr()
	err.Append(os.ErrPermission)
	err.Append(os.ErrNotExist)
	fmt.Println(err.Error())
	fmt.Println(err.Errors())
	// Output:
	// 2 errors; permission denied; file does not exist
	// [permission denied file does not exist]
}
