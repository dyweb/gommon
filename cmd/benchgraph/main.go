package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	bs := parseFile(os.Args[1])
	var groups map[string]map[string]*parse.Benchmark
	groups = make(map[string]map[string]*parse.Benchmark)
	for name, bb := range bs {
		parent, sub := extractParent(name)
		// group by parent
		//fmt.Print(parent, " ", sub, len(bb), " ")
		if len(bb) != 1 {
			fmt.Printf("expect only one benchmark, but got %d", len(bb))
		}
		b := bb[0]
		//fmt.Print(b.N, b.NsPerOp, b.AllocedBytesPerOp, b.AllocsPerOp, b.MBPerS, b.Measured, b.Ord, "\n")
		group, ok := groups[parent]
		if !ok {
			group = make(map[string]*parse.Benchmark)
		}
		group[sub] = b
		groups[parent] = group
	}
	fmt.Println("groups", len(groups))
	// one graph for each group
}

func usage() {
	fmt.Printf("usage: %s bench.txt\n\n", os.Args[0])
}

func fatal(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseFile(path string) parse.Set {
	f, err := os.Open(path)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	b, err := parse.ParseSet(f)
	if err != nil {
		fatal(err)
	}
	return b
}

func extractParent(bechmarkName string) (string, string) {
	slash := strings.Index(bechmarkName, "/")
	if slash == -1 {
		return bechmarkName, ""
	}
	return bechmarkName[:slash], bechmarkName[slash+1:]
}
