// cssGen is a temporary code generator for the myitcv.io/react.CSS type
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"myitcv.io/gogenerate"
)

var attrs = map[string]typ{
	"Height":    typ{"height", "string"},
	"MaxHeight": typ{"maxHeight", "string"},
	"MinHeight": typ{"minHeight", "string"},
	"Overflow":  typ{"overflow", "string"},
	"Resize":    typ{"resize", "string"},
	"Width":     typ{"width", "string"},
}

const (
	cssGenCmd = "cssGen"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(cssGenCmd + ": ")

	flag.Parse()

	wd, err := os.Getwd()
	if err != nil {
		fatalf("unable to get working directory: %v", err)
	}

	envFileName, ok := os.LookupEnv(gogenerate.GOFILE)
	if !ok {
		fatalf("env not correct; missing %v", gogenerate.GOFILE)
	}

	fp := filepath.Join(wd, envFileName)

	ofName, ok := gogenerate.NameFileFromFile(fp, cssGenCmd)
	if !ok {
		fatalf("could not generate generated filename from %q (with cmd %q)", envFileName, cssGenCmd)
	}

	buf := bytes.NewBuffer(nil)

	t, err := template.New("t").Parse(tmpl)
	if err != nil {
		fatalf("could not parse template: %v", err)
	}

	err = t.Execute(buf, attrs)
	if err != nil {
		fatalf("could not execute template: %v", err)
	}

	toWrite := buf.Bytes()
	out, err := format.Source(toWrite)
	if err == nil {
		toWrite = out
	}

	_, err = gogenerate.WriteIfDiff(toWrite, ofName)
	if err != nil {
		fatalf("could not write %v: %v", ofName, err)
	}
}

type typ struct {
	Attr string
	Type string
}

var tmpl = `
 // Code generated by cssGen. DO NOT EDIT.

package react

import "github.com/gopherjs/gopherjs/js"

// CSS defines CSS attributes for HTML components
//
type CSS struct {
	o *js.Object

	{{range $k, $v := .}}
	{{$k}} {{$v.Type}}
	{{end}}
}

// TODO: until we have a resolution on
// https://github.com/gopherjs/gopherjs/issues/236 we define hack() below

func (c *CSS) hack() *CSS {
	if c == nil {
		return nil
	}

	o := object.New()

	{{range $k, $v := .}}
	o.Set("{{$v.Attr}}", c.{{$k}})
	{{end}}

	return &CSS{o: o}
}
`

func fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
