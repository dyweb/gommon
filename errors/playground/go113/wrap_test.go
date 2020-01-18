package go113

// tests go1.13 error wrapping

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestWrap(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e1")
	e3 := fmt.Errorf("e2 wrap e1: %w", e1)
	t.Log(errors.Is(e2, e1)) // false
	t.Log(errors.Is(e3, e1)) // true
}

func TestAs(t *testing.T) {
	var perr *os.PathError
	e1 := fmt.Errorf("path error %w", &os.PathError{Path: "foo"})
	t.Log(errors.As(e1, &perr))
	t.Log(perr.Path)
}
