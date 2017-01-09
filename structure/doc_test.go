package structure

import "fmt"

func ExampleNewSet() {
	s := NewSet("a", "b", "c")
	fmt.Println(s.Size())
	// Output: 3
}
