package main

import "testing"

func TestEchartOption_Render(t *testing.T) {
	o := EchartOption{
		Name:   "totalRequestChart",
		Title:  "Total Success Request in 5 seconds",
		Legend: []string{"Xephon-K(Mem)", "Xephon-K(Cassandra)", "KairosDB", "InfluxDB"},
		XAxis:  []string{"10", "100", "1000", "5000"},
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

	o.Series = append(o.Series, s1, s2, s3, s4)
	o.Render()
}
