package main

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
<div id="{{.Name}}" style="width: 600px;height:400px;"></div>
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
</script>
{{ end }}
</body>
</html>
`