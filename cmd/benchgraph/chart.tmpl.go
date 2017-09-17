package main

var chartTemplate = `
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
//            orient: 'vertical',
		orient: 'horizontal',
//            left: 'right',
		top: 'bottom'
	},
	xAxis: {
		name: 'number of concurrent clients',
		nameLocation: 'middle',
		nameGap: 20,
		data: {{ json .XAxis }}
	},
	yAxis: {},
	series: {{ json .Series }}
};

{{.Name}}.setOption({{.Name}}Option);
`