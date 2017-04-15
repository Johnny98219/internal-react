package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	templateHeader = "Template generated by " + reactGenCmd
)

func doinit(wd string, tmpl string) {
	if tmpl != "minimal" {
		panic(fmt.Errorf("unknown template %q", tmpl))
	}

	for fn, c := range minimal {
		b := bytes.NewBuffer(nil)

		tmpl := struct {
			DirName string
		}{
			DirName: filepath.Base(wd),
		}

		t := template.New(fn)
		_, err := t.Parse(c)
		if err != nil {
			panic(fmt.Errorf("failed to parse template %v: %v", fn, err))
		}

		err = t.Execute(b, tmpl)
		if err != nil {
			panic(fmt.Errorf("failed to execute template %v: %v", fn, err))
		}

		toWrite := b.Bytes()

		fp := filepath.Join(wd, fn)

		if strings.HasSuffix(fn, ".go") {

			out, err := fmtBuf(b)
			if err == nil {
				toWrite = out.Bytes()
			}
		}

		err = ioutil.WriteFile(fp, toWrite, 0644)
		if err != nil {
			panic(fmt.Errorf("failed to write file %v: %v", fp, err))
		}
	}

	gg := exec.Command("go", "generate")
	gg.Dir = wd
	gg.Stderr = os.Stderr
	gg.Stdout = os.Stdout

	err := gg.Run()
	if err != nil {
		fatalf("failed to run go generate: %v", err)
	}
}

var minimal = map[string]string{
	// app.go
	"app.go": `// ` + templateHeader + `

package main

import (
	r "myitcv.io/react"
)

type AppDef struct {
	r.ComponentDef
}

func App() *AppDef {
	res := new(AppDef)
	r.BlessElement(res, nil)
	return res
}

func (a *AppDef) Render() r.Element {
	return r.Div(nil,
		r.H1(nil,
			r.S("Hello World"),
		),
		r.P(nil,
			r.S("This is my first GopherJS React App."),
		),
	)
}
`,

	// main.go
	"main.go": `// ` + templateHeader + `

package main

import (
	r "myitcv.io/react"

	"honnef.co/go/js/dom"
)

  //go:generate reactGen

var document = dom.GetWindow().Document()

func main() {
	domTarget := document.GetElementByID("app")

	r.Render(App(), domTarget)
}
`,

	// index.html
	"index.html": `<!--` + templateHeader + `-->

<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Hello World</title>
  </head>
  <body>
    <div id="app"></div>
    <script src="{{.DirName}}.js"></script>
  </body>
</html>
`,
}
