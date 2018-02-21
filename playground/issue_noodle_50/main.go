package main

import "fmt"

// https://github.com/dyweb/gommon/issues/50

type I interface {
	GetName() string
}

var _ I = (*Info)(nil)

type Info struct {
	Name string
}

func (i *Info) GetName() string {
	return i.Name
}

type Dir struct {
	Entries []Info
}

func main() {
	dir := Dir{
		Entries:[]Info{
			{Name: "a"},
			{Name: "b"},
		},
	}
	files := make([]I, 0, len(dir.Entries))
	for _, f := range dir.Entries {
		files = append(files, &f)
	}
	for _, f := range files {
		fmt.Println(f.GetName())
	}
}