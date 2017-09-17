package main

import (
	"encoding/json"
	"text/template"
	"bytes"
)

type Series struct {
	Name string    `json:"name"`
	Type string    `json:"type"`
	Data []float64 `json:"data"` // TODO: float or int ?
}

type ECharts struct {
	Title string
	Charts []EChartOption
}

type EChartOption struct {
	Name   string
	Title  string
	Legend []string
	YAxis  []string
	Series []Series
}

func jsonPipe(d interface{}) (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (charts *ECharts) Render() ([]byte, error) {
	funcMap := template.FuncMap{
		"json": jsonPipe,
	}
	tmpl, err := template.New("chart").Funcs(funcMap).Parse(chartsTemplate)
	if err != nil {
		return nil, err
	}
	var b bytes.Buffer
	err = tmpl.Execute(&b, charts)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//func (opt *EChartOption) Render() ([]byte, error) {
//	funcMap := template.FuncMap{
//		"json": jsonPipe,
//	}
//	tmpl, err := template.New("chart").Funcs(funcMap).Parse(chartTemplate)
//	if err != nil {
//		return nil, err
//	}
//	err = tmpl.Execute(os.Stdout, opt)
//	if err != nil {
//		return nil, err
//	}
//	return nil, nil
//}
