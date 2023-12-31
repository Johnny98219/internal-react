package main

import (
	"go/ast"
	"go/format"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"myitcv.io/gogenerate"
)

type propsGen struct {
	*coreGen

	Recv  string
	Name  string
	TName string
	Doc   string

	Fields []field
}

func (g *gen) genProps(defName string, t typeFile) {
	name := strings.TrimPrefix(defName, propsTypeTmplPrefix)

	r, _ := utf8.DecodeRuneInString(name)

	pg := &propsGen{
		coreGen: newCoreGen(g),
		Name:    name,
		Recv:    string(unicode.ToLower(r)),
	}

	var doc string

	if v := t.ts.Doc; v != nil {
		for _, c := range v.List {
			doc = doc + strings.Replace(c.Text, defName, name, -1)
		}
	}

	pg.Doc = doc

	fe := &fieldExploder{
		first:  true,
		pkgStr: g.pkgImpPath,
		sn:     t.ts.Name.Name,
		imps:   make(map[*ast.ImportSpec]struct{}),
	}

	err := fe.explode(&t)
	if err != nil {
		fatalf("could not explode fields: %v", err)
	}

	sort.Slice(fe.fields, func(i, j int) bool {
		return fe.fields[i].TName < fe.fields[j].TName
	})

	if len(fe.fields) > 2 {
		for i := 1; i < len(fe.fields)-1; i++ {
			if fe.fields[i].IsEvent != fe.fields[i-1].IsEvent {
				fe.fields[i].GapBefore = true
			}
		}
	}

	pg.Fields = fe.fields

	pg.pf("// Code generated by %v. DO NOT EDIT.\n", reactGenCmd)
	pg.pln()
	pg.pf("package %v\n", pg.pkg)

	if !g.isReactCore {
		pg.pf("import \"%v\"\n", reactPkg)
	}

	for is := range fe.imps {
		if is.Name != nil {
			pg.pf("import %v %v", is.Name.Name, is.Path.Value)
		} else {
			pg.pf("import %v", is.Path.Value)
		}
	}

	pg.pln()

	// TODO the ID/Key hack below feels fragile...
	pg.pt(`
	{{ $recv := .Recv }}

	{{.Doc}}
	type {{.Name}} struct {
	{{range $i, $v := .Fields}}
		{{if $v.GapBefore}}
		{{end -}}
		{{$v.Name}} {{$v.Type}}
	{{- end}}
	}

	func ({{$recv}} *{{.Name}}) assign(v *_{{.Name}}) {
		{{- range .Fields}}
			{{ if eq .TName "Ref" }}
			if {{$recv}}.Ref != nil {
				v.o.Set("ref", {{$recv}}.Ref.Ref)
			}
			{{ else if eq .TName "DataSet" }}
			if {{$recv}}.DataSet != nil {
				for dk, dv := range {{$recv}}.DataSet {
					v.o.Set("data-"+dk, dv)
				}
			}
			{{ else if eq .TName "AriaSet" }}
			if {{$recv}}.AriaSet != nil {
				for dk, dv := range {{$recv}}.AriaSet {
					v.o.Set("aria-"+dk, dv)
				}
			}
			{{else}}
			{{ if .Omit }}
				if {{$recv}}.{{.TName}} != "" {
					v.{{.TName}} = {{$recv}}.{{.TName}}
				}
			{{else}}
			{{if .IsEvent}}
				if {{$recv}}.{{.TName}} != nil {
					v.o.Set("{{.FName}}", {{$recv}}.{{.TName}}.{{.TName}})
				}
			{{else if eq .Name "Style"}}
				// TODO: until we have a resolution on
				// https://github.com/gopherjs/gopherjs/issues/236
				v.{{.TName}} = {{$recv}}.{{.TName}}.hack()
			{{else}}
				v.{{.TName}} = {{$recv}}.{{.TName}}
			{{end}}
			{{end}}
			{{end}}
		{{- end}}
	}
	`, pg)

	ofName := gogenerate.NameFile(name, reactGenCmd)
	toWrite := pg.buf.Bytes()

	out, err := format.Source(toWrite)
	if err == nil {
		toWrite = out
	}

	wrote, err := gogenerate.WriteIfDiff(toWrite, ofName)
	if err != nil {
		fatalf("could not write %v: %v", ofName, err)
	}

	if wrote {
		infof("writing %v", ofName)
	} else {
		infof("skipping writing of %v; it's identical", ofName)
	}

}
