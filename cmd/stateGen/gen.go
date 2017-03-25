package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/imports"

	"github.com/myitcv/gogenerate"
)

const (
	rootVar      = "root"
	leafTypeName = "Leaf"
	nodePrefix   = "_Node_"
)

type gen struct {
	fset *token.FileSet

	dir string

	buf *bytes.Buffer

	rootType *node

	roots   []*ast.ValueSpec
	nodes   map[string]node
	pkg     *ast.Package
	pkgName string

	imports map[*ast.ImportSpec]bool

	file *ast.File

	// a map from the _Node_XYZ name
	seenNodes    map[string]bool
	nodesToVisit []node

	seenLeaves    map[string]bool
	leavesToVisit []leafField

	stderr io.Writer
	failed bool
}

func dogen(stderr io.Writer, dir, license string) bool {
	fset := token.NewFileSet()

	notGenByUs := func(fi os.FileInfo) bool {
		return !gogenerate.FileGeneratedBy(fi.Name(), stateGenCmd)
	}

	pkgs, err := parser.ParseDir(fset, dir, notGenByUs, 0)
	if err != nil {
		panic(fmt.Errorf("unable to parse directory %v: %v", dir, err))
	}

	failed := false

	for pn, pkg := range pkgs {
		g := &gen{
			fset:    fset,
			dir:     dir,
			pkg:     pkg,
			pkgName: pn,

			buf: bytes.NewBuffer(nil),

			imports: make(map[*ast.ImportSpec]bool),

			nodes:      make(map[string]node),
			seenNodes:  make(map[string]bool),
			seenLeaves: make(map[string]bool),

			stderr: stderr,
		}

		g.parse()
		if !g.ok() {
			failed = true
			continue
		}

		if g.rootType == nil {
			continue
		}

		g.pf("package %v\n", pn)

		g.pf(`
		import "path"
		`)

		for i := range g.imports {
			if i.Name != nil {
				g.pf("import %v %v\n", i.Name.Name, i.Path.Value)
			} else {
				g.pf("import %v\n", i.Path.Value)
			}
		}

		g.gen()

		fn := gogenerate.NameFile(pn, stateGenCmd)
		fp := filepath.Join(dir, fn)

		toWrite := g.buf.Bytes()

		res, err := imports.Process(fn, toWrite, nil)
		if err == nil {
			toWrite = res
		}

		wrote, err := gogenerate.WriteIfDiff(toWrite, fp)
		if err != nil {
			panic(fmt.Errorf("unable to write to %v: %v", fp, err))
		}

		if wrote {
			infof("wrote %v\n", fp)
		} else {
			infof("skipped writing %v; identical\n", fp)
		}
	}

	return !failed
}

type node struct {
	Name     string
	children []field
	leaves   []leafField
}

type field struct {
	Name string
	Type string
}

type leafField struct {
	Name     string
	Type     string
	LeafType string
}

func (g *gen) parse() {
	for _, f := range g.pkg.Files {
		g.file = f

		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok {
				continue
			}

			switch gd.Tok {
			case token.TYPE:
				for _, s := range gd.Specs {
					g.parseNode(s.(*ast.TypeSpec))
				}
			case token.VAR:
				for _, s := range gd.Specs {
					s := s.(*ast.ValueSpec)

					if len(s.Names) != 1 {
						continue
					}

					if s.Names[0].Name != rootVar {
						continue
					}

					_, ok := s.Type.(*ast.Ident)
					if !ok {
						continue
					}

					g.roots = append(g.roots, s)
				}
			}
		}
	}
}

func (g *gen) parseNode(s *ast.TypeSpec) {
	st, ok := s.Type.(*ast.StructType)
	if !ok {
		return
	}

	tn := s.Name.Name

	if !strings.HasPrefix(tn, nodePrefix) {
		return
	}

	n := strings.TrimPrefix(tn, nodePrefix)

	var children []field
	var leaves []leafField

	for _, f := range st.Fields.List {
		var id *ast.Ident

		switch typ := f.Type.(type) {
		case *ast.Ident:
			id = typ
		case *ast.StarExpr:
			if v, ok := typ.X.(*ast.Ident); ok {
				id = v
			}
		}

		if id != nil && strings.HasPrefix(id.Name, nodePrefix) {
			for _, n := range f.Names {
				children = append(children, field{
					Name: n.Name,
					Type: strings.TrimPrefix(id.Name, nodePrefix),
				})
			}
		} else {
			typ, leafTyp := g.addImports(f.Type)
			for _, n := range f.Names {
				leaves = append(leaves, leafField{
					Name:     n.Name,
					Type:     typ,
					LeafType: leafTyp,
				})
			}
		}
	}

	g.nodes[s.Name.Name] = node{
		Name:     n,
		children: children,
		leaves:   leaves,
	}
}

func (g *gen) addImports(exp ast.Expr) (string, string) {
	finder := &importFinder{
		imports: g.file.Imports,
		matches: g.imports,
	}

	ast.Walk(finder, exp)

	es := g.expString(exp)

	s := strings.Replace(es, ".", "", -1)

	if strings.HasPrefix(s, "*") {
		s = strings.TrimPrefix(s, "*")
		s = s + "P"
	}

	r, l := utf8.DecodeRune([]byte(s))

	return es, string(unicode.ToUpper(r)) + string(s[l:]) + leafTypeName
}

func (g *gen) ok() bool {
	if len(g.roots) == 0 {
		return true
	}

	if v := len(g.roots); v > 1 {
		g.errorf("expected 1 root, found %v\n", v)
		for _, v := range g.roots {
			g.errorf("  %v\n", g.fset.Position(v.Pos()))
		}

		return false
	}

	r := g.roots[0]
	rtn := r.Type.(*ast.Ident)

	rt, ok := g.nodes[rtn.Name]
	if !ok {
		g.errorf("need root type to be node type; instead it was %v", rtn.Name)
	}

	g.rootType = &rt

	return !g.failed
}

func (g *gen) gen() {
	var h node
	var l leafField

	g.nodesToVisit = append(g.nodesToVisit, *g.rootType)

	for len(g.nodesToVisit) != 0 {
		h, g.nodesToVisit = g.nodesToVisit[0], g.nodesToVisit[1:]
		g.genNode(h)
	}

	for len(g.leavesToVisit) != 0 {
		l, g.leavesToVisit = g.leavesToVisit[0], g.leavesToVisit[1:]
		g.genLeaf(l)
	}

	g.pt(`
	func NewRoot() *{{.Name}} {
		r := &rootNode{
			store: make(map[string]interface{}),
			cbs:   make(map[string]map[*Sub]struct{}),
			subs:  make(map[*Sub]struct{}),
		}

		return new{{.Name}}(r, "")
	}
	`, g.rootType)

	g.pf(`
	type Node interface {
		Subscribe(cb func()) *Sub
	}

	type Sub struct {
		*rootNode
		prefix string
		cb     func()
	}

	func (s *Sub) Clear() {
		s.rootNode.unsubscribe(s)
	}

	var NoSuchSubErr = errors.New("No such sub")

	type rootNode struct {
		store map[string]interface{}
		cbs   map[string]map[*Sub]struct{}
		subs  map[*Sub]struct{}
	}

	func (r *rootNode) subscribe(prefix string, cb func()) *Sub {

		res := &Sub{
			cb:     cb,
			prefix: prefix,
			rootNode: r,
		}

		l, ok := r.cbs[prefix]
		if !ok {
			l = make(map[*Sub]struct{})
			r.cbs[prefix] = l
		}

		l[res] = struct{}{}
		r.subs[res] = struct{}{}

		return res
	}

	func (r *rootNode) unsubscribe(s *Sub) {
		if _, ok := r.subs[s]; !ok {
			panic(NoSuchSubErr)
		}

		l, ok := r.cbs[s.prefix]
		if !ok {
			panic("Real problems...")
		}

		delete(l, s)
		delete(r.subs, s)
	}

	func (r *rootNode) get(k string) (interface{}, bool) {
		v, ok := r.store[k]
		return v, ok
	}

	func (r rootNode) set(k string, v interface{}) {
		if curr, ok := r.store[k]; ok && v == curr {
			return
		}

		r.store[k] = v

		parts := strings.Split(k, "/")

		var subs []*Sub

		var kk string

		for _, p := range parts {
			kk = path.Join(kk, p)

			if ll, ok := r.cbs[kk]; ok {
				for k := range ll {
					subs = append(subs, k)
				}
			}

		}

		for _, s := range subs {
			s.cb()
		}
	}
	`)
}

func (g *gen) genLeaf(n leafField) {
	g.pt(`
	type {{.LeafType}} struct {
		*rootNode
		prefix string
	}

	var _ Node = new({{.LeafType}})

	func new{{.LeafType}}(r *rootNode, prefix string) *{{.LeafType}} {
		prefix = path.Join(prefix, "{{.LeafType}}")

		return &{{.LeafType}}{
			rootNode:   r,
			prefix: prefix,
		}
	}

	func (m *{{.LeafType}}) Get() {{.Type}} {
		var res {{.Type}}
		if v, ok := m.rootNode.get(m.prefix); ok {
			return v.({{.Type}})
		}
		return res
	}

	func (m *{{.LeafType}}) Set(v {{.Type}}) {
		m.rootNode.set(m.prefix, v)
	}

	func (m *{{.LeafType}}) Subscribe(cb func()) *Sub {
		return m.rootNode.subscribe(m.prefix, cb)
	}
	`, n)

}

func (g *gen) genNode(n node) {
	g.pt(`
	var _ Node = new({{.Name}})

	type {{.Name}} struct {
		*rootNode
		prefix string

	`, n)

	for _, c := range n.children {
		if !g.seenNodes[c.Type] {
			g.nodesToVisit = append(g.nodesToVisit, g.nodes[nodePrefix+c.Type])
			g.seenNodes[c.Type] = true
		}

		g.pt(`
		_{{.Name}} *{{.Type}}
		`, c)
	}

	for _, l := range n.leaves {
		if !g.seenLeaves[l.LeafType] {
			g.leavesToVisit = append(g.leavesToVisit, l)
			g.seenLeaves[l.LeafType] = true
		}

		g.pt(`
		_{{.Name}} *{{.LeafType}}
		`, l)
	}

	g.pt(`
	}

	func new{{.Name}}(r *rootNode, prefix string) *{{.Name}} {
		prefix = path.Join(prefix, "{{.Name}}")

		res := &{{.Name}}{
			rootNode:   r,
			prefix: prefix,
		}
	`, n)

	for _, c := range n.children {
		g.pt(`
		res._{{.Name}} = new{{.Type}}(r, prefix)
		`, c)
	}
	for _, l := range n.leaves {
		g.pt(`
		res._{{.Name}} = new{{.LeafType}}(r, prefix)
		`, l)
	}

	g.pt(`
		return res
	}

	func (n *{{.Name}}) Subscribe(cb func()) *Sub {
		return n.rootNode.subscribe(n.prefix, cb)
	}
	`, n)

	for _, c := range n.children {
		tmpl := struct {
			Node  node
			Child field
		}{
			Node:  n,
			Child: c,
		}
		g.pt(`
		func (n *{{.Node.Name}}) {{.Child.Name}}() *{{.Child.Type}} {
			return n._{{.Child.Name}}
		}
		`, tmpl)
	}

	for _, l := range n.leaves {
		tmpl := struct {
			Node node
			Leaf leafField
		}{
			Node: n,
			Leaf: l,
		}
		g.pt(`
		func (n *{{.Node.Name}}) {{.Leaf.Name}}() *{{.Leaf.LeafType}} {
			return n._{{.Leaf.Name}}
		}
		`, tmpl)
	}
}

func (g *gen) errorf(format string, args ...interface{}) {
	g.failed = true
	fmt.Fprintf(g.stderr, format, args...)
}

func (g *gen) expString(e interface{}) string {
	b := bytes.NewBuffer(nil)
	err := printer.Fprint(b, g.fset, e)
	if err != nil {
		panic(err)
	}

	return b.String()
}

func (g *gen) pf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *gen) pln(args ...interface{}) {
	fmt.Fprintln(g.buf, args...)
}

func (g *gen) pt(tmpl string, val interface{}) {
	// on the basis most templates are for convenience define inline
	// as raw string literals which start the ` on one line but then start
	// the template on the next (for readability) we strip the first leading
	// \n if one exists
	tmpl = strings.TrimPrefix(tmpl, "\n")

	t := template.New("tmp")

	_, err := t.Parse(tmpl)
	if err != nil {
		fatalf("unable to parse template: %v", err)
	}

	err = t.Execute(g.buf, val)
	if err != nil {
		fatalf("cannot execute template: %v", err)
	}
}

type importFinder struct {
	imports []*ast.ImportSpec
	matches map[*ast.ImportSpec]bool
}

func (i *importFinder) Visit(node ast.Node) ast.Visitor {
	switch node := node.(type) {
	case *ast.SelectorExpr:
		if x, ok := node.X.(*ast.Ident); ok {
			for _, imp := range i.imports {
				if imp.Name != nil {
					if x.Name == imp.Name.Name {
						i.matches[imp] = true
					}
				} else {
					cleanPath := strings.Trim(imp.Path.Value, "\"")
					parts := strings.Split(cleanPath, "/")
					if x.Name == parts[len(parts)-1] {
						i.matches[imp] = true
					}
				}
			}

		}
	}

	return i
}
