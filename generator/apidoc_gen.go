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
}

const apiDocTpl = `
/**
 * 
 * @api {{"{"}}{{.Method}}{{"}"}} {{.Path}} {{.Title}}
{{if .APIName}} * @apiName {{.APIName}}{{end}}
{{if .Group}} * @apiGroup {{.Group}}{{end}}
{{if .Version}} * @apiVersion {{.Version}}{{end}}
 * 
 {{range $k, $v := .Params}}* @apiParam  { {{$v.Type}} } {{$k}} {{$v.Description}}{{"\n"}} {{end}}
 * 
 {{range $k, $v := .Success}}* @apiSuccess { {{$v.Type}} } {{$k}} {{$v.Description}}{{"\n"}} {{end}}
 * 
 * @apiParamExample  { Object } Request-Example:
{
 {{range $k, $v := .Params}}"{{$k}}" : "{{$v.Type}}",{{"\n"}} {{end}}
}
 * 
 * @apiSuccessExample { Object } Success-Response:
{
	"data": {{ if .IsSuccessArray }} [
		{
			{{range $k, $v := .Success}} "{{$k}}" : "{{$v.Type}}",{{"\n\t\t\t"}}{{end}} 
		},
	]
	{{else}}{
		{{range $k, $v := .Success}} "{{$k}}" : "{{$v.Type}}",{{"\n\t\t\t"}}{{end}} 
	}
	{{ end }}
}
 * 
 * 
 */`

// Gen genera apidoc like doc
func Gen(docs []APIDoc, filename string) {
	regex, err := regexp.Compile("\n\n|,\n $")
	if err != nil {
		return
	}
	var b bytes.Buffer
	temp := bufio.NewWriter(&b)
	for _, doc := range docs {
		if doc.APIName == "" {
			continue
		}
		t := template.New("apidoc")
		t, _ = t.Parse(apiDocTpl)
		t.Execute(temp, doc)
		temp.Flush()
	}
	s := regex.ReplaceAllString(b.String(), "")
	err = ioutil.WriteFile(filename, []byte(s), 0644)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
