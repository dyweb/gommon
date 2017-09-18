package main

import (
	"bytes"
	"encoding/json"
	"text/template"
)

const (
	defaultWidth  = 1024
	defaultHeight = 600
)

type ECharts struct {
	Title  string
	Charts []EChartOption
}

type EChartOption struct {
	Name   string
	Title  string
	Width  int // TODO: width and hight is not used in template yet
	Height int
	Legend []string
	YAxis  []string
	Series []Series
}

type Series struct {
	Name string    `json:"name"`
	Type string    `json:"type"`
	Data []float64 `json:"data"`
}

func NewEChartOption() *EChartOption {
	return &EChartOption{
		Width:  defaultWidth,
		Height: defaultHeight,
	}
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

var chartsTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/echarts/3.7.1/echarts.min.js"></script>
</head>
<body>
{{ range .Charts }}
<div id="{{.Name}}" style="width: 1024px;height:600px;"></div>
<br/>
{{ end }}
{{ range .Charts }}
<script type="text/javascript">
var {{.Name}} = echarts.init(document.getElementById('{{.Name}}'));

var {{.Name}}Option = {
	title: {
		text: '{{.Title}}'
	},
	tooltip: {},
	toolbox: {
		feature: {
			dataView: {show: true, readOnly: true},
			magicType: {show: true, type: ['line', 'bar']},
			restore: {show: true},
			saveAsImage: {show: true}
		}
	},
	legend: {
		data: {{ json .Legend }},
		orient: 'horizontal',
		top: 'bottom'
	},
	yAxis: {
		type: 'category',
		data: {{ json .YAxis }},
		axisLabel: {
			interval: 0,
			rotate: 30
		}
	},
	xAxis: {
		type: 'value'
	},
	series: {{ json .Series }}
};

{{.Name}}.setOption({{.Name}}Option);
</script>
{{ end }}
</body>
</html>
`
