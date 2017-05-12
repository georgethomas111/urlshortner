// File stores the contents in the template file so that its part of the package.
package couchdb

import (
	"strings"
	"text/template"
)

var DESIGNTMPL *template.Template = template.Must(template.New("design").Parse(strings.Replace(strings.Replace(`
{
	"_id": "{{.Id}}",
	{{if .RevStatus}}
	"_rev":"{{.Rev}}",
	{{end}}
	"views": {
		{{range .Views}}

		   "{{.Name}}": {
		   {{if .RawJson}}
			{{.RawJson}}
		   {{else}}
			"map": " \
			function({{.VariableName}}) { \
\
				function get_value(doc, keys) { \
					var value = {}; \
					for(key in keys) { \
						value[keys[key]] = doc[keys[key]];\
					} \
					return value; \
					}\
				{{if .CondStatus}}\
					if({{.Condition}}) { \
						emit({{.KeyName}}, get_value({{.VariableName}},  [{{.EmitStr}}]));\
						} \
				{{else}} \
						emit({{.KeyName}}, get_value({{.VariableName}},  [{{.EmitStr}}]));\
				{{end}} \
			}"
		   {{end}}
		   },

		{{end}}
		   "{{.LastView.Name}}": {
	            {{if .LastView.RawJson}}
			{{.LastView.RawJson}}
		    {{else}}
			"map": " \
			function({{.LastView.VariableName}}) { \
\
				function get_value(doc, keys) { \
					var value = {}; \
					for(key in keys) { \
						value[keys[key]] = doc[keys[key]];\
					} \
					return value; \
					}\
				{{if .LastView.CondStatus}}\
					if({{.LastView.Condition}}) { \
						emit({{.LastView.KeyName}}, get_value({{.LastView.VariableName}},  [{{.LastView.EmitStr}}]));\
						} \
				{{else}} \
						emit({{.LastView.KeyName}}, get_value({{.LastView.VariableName}},  [{{.LastView.EmitStr}}]));\
				{{end}} \
			}"
		    {{end}}
		   }
	},
	"language": "javascript"
}
`, "\\\n", "", -1), "\t", "", -1)))
