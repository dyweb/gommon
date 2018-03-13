package testutil

import (
	"testing"
	"time"
)

type Foo struct {
	Bar string    `json:"bar"`
	T   time.Time `json:"time"`
}

func TestPrintTidyJson(t *testing.T) {
	a := Foo{Bar: "bar2", T: time.Now()}
	PrintTidyJson(t, a)
}
