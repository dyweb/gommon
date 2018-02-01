package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestECharts_Render(t *testing.T) {
	t.Skip("deprecated, benchgraph will be merged to https://github.com/benchhub/benchboard")

	c1 := EChartOption{
		Name:   "totalRequestChart",
		Title:  "Total Success Request in 5 seconds",
		Legend: []string{"Xephon-K(Mem)", "Xephon-K(Cassandra)", "KairosDB", "InfluxDB"},
		Axis:   []string{"10", "100", "1000", "5000"},
	}
	s1 := Series{
		Name: "Xephon-K(Mem)",
		Type: "bar",
		Data: []float64{12327, 21099, 31791, 12279},
	}
	s2 := Series{
		Name: "Xephon-K(Cassandra)",
		Type: "bar",
		Data: []float64{7931, 11336, 14590, 8703},
	}
	s3 := Series{
		Name: "KairosDB",
		Type: "bar",
		Data: []float64{15561, 26154, 26939, 16506},
	}
	s4 := Series{
		Name: "InfluxDB",
		Type: "bar",
		Data: []float64{118, 139, 131, 130},
	}
	c1.Series = append(c1.Series, s1, s2, s3, s4)

	c2 := c1
	c2.Name = "totalRequestChart2"
	c2.Title = "Total Success Request in 5 seconds Dup"

	charts := ECharts{Title: "TSDB Bench"}
	charts.Charts = append(charts.Charts, c1, c2)
	b, err := charts.Render()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(b))
	ioutil.WriteFile("tsdb-bench.html", b, os.ModePerm)
}
