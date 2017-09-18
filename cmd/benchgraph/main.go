package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"golang.org/x/tools/benchmark/parse"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	benchmarkOutput := os.Args[1]
	bs := parseFile(benchmarkOutput)
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
	fmt.Printf("output has %d groups\n", len(groups))
	charts := ECharts{Title: benchmarkOutput}
	sortedGroups := make([]string, 0, len(groups))
	for groupName := range groups {
		sortedGroups = append(sortedGroups, groupName)
	}
	sort.Strings(sortedGroups)
	// one graph for each group
	for _, groupName := range sortedGroups {
		group := groups[groupName]
		// nanoseconds per iteration
		nsChart := EChartOption{
			// TODO: groupName might not be a valid javascript variable name, sanitize?
			Name:   groupName + "Ns",
			Title:  groupName + " NsPerOp",
			Legend: []string{"nanosecond per iteration"},
		}
		bChart := EChartOption{
			Name:   groupName + "Bytes",
			Title:  groupName + " BytesPerOp",
			Legend: []string{"bytes allocated per iteration"},
		}
		nsSeries := Series{
			Name: "nanosecond per iteration",
			Type: "bar",
		}
		bSeries := Series{
			Name: "bytes allocated per iteration",
			Type: "bar",
		}
		keys := make([]string, 0, len(group))
		for subName := range group {
			keys = append(keys, subName)
		}
		sort.Strings(keys)
		for _, subName := range keys {
			b := group[subName]
			nsChart.YAxis = append(nsChart.YAxis, subName)
			nsSeries.Data = append(nsSeries.Data, float64(b.NsPerOp))
			bChart.YAxis = append(bChart.YAxis, subName)
			bSeries.Data = append(bSeries.Data, float64(b.AllocedBytesPerOp))
		}
		nsChart.Series = append(nsChart.Series, nsSeries)
		bChart.Series = append(bChart.Series, bSeries)
		charts.Charts = append(charts.Charts, nsChart, bChart)
	}
	b, err := charts.Render()
	if err != nil {
		fatal(err)
	}
	htmlLocation := benchmarkOutput + ".html"
	if len(os.Args) > 2 {
		htmlLocation = os.Args[2]
	}
	if err := ioutil.WriteFile(htmlLocation, b, os.ModePerm); err != nil {
		fatal(err)
	}
	fmt.Printf("charts %s generated from %s\n", htmlLocation, benchmarkOutput)
}

func usage() {
	fmt.Printf("usage: %s bench.txt bech.html\n\n", os.Args[0])
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
	// TODO: I think benchcmp might allow concat output of different runs, so there is more than one result for one benchmark, it has select best ....
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
