package main

import (
	"encoding/json"
	"text/template"
	"os"
)

type Series struct {
	Name string    `json:"name"`
	Type string    `json:"type"`
	Data []float64 `json:"data"` // TODO: float or int ?
}

type EchartOption struct {
	Name   string
	Title  string
	Legend []string
	XAxis  []string
	Series []Series
	//Legend struct{
	//	Data []string
	//	//Orient string
	//	//Top string
	//}
}

func jsonPipe(d interface{}) (string, error) {
	b, err := json.Marshal(d)
	return string(b), err
}

func (opt *EchartOption) Render() ([]byte, error) {
	funcMap := template.FuncMap{
		"json": jsonPipe,
	}
	tmpl, err := template.New("chart").Funcs(funcMap).Parse(chartTemplate)
	if err != nil {
		return nil, err
	}
	err = tmpl.Execute(os.Stdout, opt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
