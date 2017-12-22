package stdlib

import (
	"testing"
	"fmt"
	"os"
)

func printArgs(args ...interface{}) {
	fmt.Println(args) // WRONG, it will have [arg1 arg2 arg3]
	fmt.Println(args...)
}

func TestFmt_Sprintf(t *testing.T) {
	fmt.Print("This a debug message", 1, "\n")
	fmt.Println("This a debug message", 1) // NOTE: Println add space between words, and \n at last
	fmt.Fprintln(os.Stdout, fmt.Sprint("This a debug message", 1))
	printArgs("Tell me", 123)
}
