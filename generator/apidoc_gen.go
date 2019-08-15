package generator

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"regexp"
	"text/template"
)

// APIDoc entity
type APIDoc struct {
	IsSuccessArray bool
	Method         string
	Path           string
	Title          string
	APIName        string
	Group          string
	Version        string
	Params         map[string]Param
	Success        map[string]Param
}

// Param stands for param
type Param struct {
	Name        string
	Type        string
	Description string
	Child       []Param
}

const apiDocTpl = `
/**
 * 
 * @api {{"{"}}{{.Method}}{{"}"}} {{.Path}} {{.Title}}
{{if .APIName}} * @apiName {{.APIName}}{{end}}
{{if .Group}} * @apiGroup {{.Group}}{{end}}
{{if .Version}} * @apiVersion {{.Version}}{{end}}
 * 
 {{range $k, $v := .Params}}
 	{{if eq $v.Type "Object"}}
		{{range $ck, $cv := $v.Child}}
 * @apiParam { {{$cv.Type}} } {{$k}}.{{$cv.Name}} {{$cv.Description}}{{"\n"}} 
		{{end}}
	{{else}}
 * @apiParam { {{$v.Type}} } {{$k}} {{$v.Description}}{{"\n"}} 
	{{end}}
 {{end}} 

 {{range $k, $v := .Success}}
 	{{if eq $v.Type "Object"}}
		{{range $ck, $cv := $v.Child}}
 * @apiSuccess { {{$cv.Type}} } {{$k}}.{{$cv.Name}} {{$cv.Description}}{{"\n"}} 
		{{end}}
	{{else}}
 * @apiSuccess { {{$v.Type}} } {{$k}} {{$v.Description}}{{"\n"}} 
	{{end}}
 {{end}}
 * 
 * @apiParamExample  { Object } Request-Example:
{
{{range $k, $v := .Params}}
{{ $length := len $v.Child }} {{ if eq $length 0 }}
	"{{$k}}" : "{{$v.Type}}",{{"\n\t\t\t"}}
	{{else}}
	"{{$k}}" : 
		{{ if eq $v.Type "Object[]"}}
		[{
				{{range $ck, $cv := $v.Child}} 
			"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
				{{end}}
		},]
		{{else}}
		{
			{{range $ck, $cv := $v.Child}} 
			"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
			{{end}}
		}
		{{end}}
	{{end}} 
{{end}}
}
 * 
 * @apiSuccessExample { Object } Success-Response:
{
	"code": 1,
	"msg": "",
	"data": {{ if .IsSuccessArray }} [
		{
		{{range $k, $v := .Success}} 
			{{ $length := len $v.Child }} {{ if eq $length 0 }}
			"{{$k}}" : "{{$v.Type}}",{{"\n\t\t\t"}}
			{{else}}
			"{{$k}}" : 
				{{ if eq $v.Type "Object[]"}}
				[{
						{{range $ck, $cv := $v.Child}} 
					"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
						{{end}}
				},]
				{{else}}
				{
					{{range $ck, $cv := $v.Child}} 
					"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
					{{end}}
				}
				{{end}}
			{{end}}
		{{end}} 
		},
	]
	{{else}}{
	{{range $k, $v := .Success}} 
		{{ $length := len $v.Child }} {{ if eq $length 0 }}
			"{{$k}}" : "{{$v.Type}}",{{"\n\t\t\t"}}
		{{else}}
			"{{$k}}" : 
				{{ if eq $v.Type "Object[]"}}
					[{
						{{range $ck, $cv := $v.Child}} 
							"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
						{{end}}
					},]
				{{else}}
				{
					{{range $ck, $cv := $v.Child}} 
						"{{$cv.Name}}" : "{{$cv.Type}}",{{"\n\t\t\t\t\t"}}
					{{end}}
				}
			{{end}}
		{{end}}
	{{end}} 
	}
	{{ end }}
}
 * 
 * 
 */`

// Gen genera apidoc like doc
func Gen(docs []APIDoc, filename string) {
	emptyLineRegex := regexp.MustCompile(`(?m)^\s*$[\r\n]*|[\r\n]+\s+\z`)
	var b bytes.Buffer
	temp := bufio.NewWriter(&b)
	for _, doc := range docs {
		if doc.APIName == "" {
			continue
		}
		t := template.New("apidoc")
		t, _ = t.Parse(apiDocTpl)
		err := t.Execute(temp, doc)
		if err != nil {
			log.Fatalln(err)
		}
		temp.Flush()
	}
	s := emptyLineRegex.ReplaceAllString(b.String(), "")
	err := ioutil.WriteFile(filename, []byte(s), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
