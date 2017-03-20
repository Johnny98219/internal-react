// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package examples

import (
	r "github.com/myitcv/gopherjs/react"
	"github.com/myitcv/gopherjs/react/examples/hellomessage"
	"github.com/myitcv/gopherjs/react/examples/immtodoapp"
	"github.com/myitcv/gopherjs/react/examples/markdowneditor"
	"github.com/myitcv/gopherjs/react/examples/timer"
	"github.com/myitcv/gopherjs/react/examples/todoapp"
	"honnef.co/go/js/xhr"
)

//go:generate reactGen
//go:generate immutableGen

// ExamplesDef is the definition of the Examples component
type ExamplesDef struct {
	r.ComponentDef
}

type tab int

const (
	tabGo tab = iota
	tabJsx
)

// Examples creates instances of the Examples component
func Examples() *ExamplesDef {
	res := new(ExamplesDef)
	r.BlessElement(res, nil)
	return res
}

type (
	_Imm_tabS map[exampleKey]tab
)

// ExamplesState is the state type for the Examples component
type ExamplesState struct {
	examples     *exampleSource
	selectedTabs *tabS
}

// ComponentWillMount is a React lifecycle method for the Examples component
func (p *ExamplesDef) ComponentWillMount() {
	if !fetchStarted {
		for i, e := range sources.Range() {
			go func(i exampleKey, e *source) {
				req := xhr.NewRequest("GET", "https://raw.githubusercontent.com/myitcv/gopherjs/master/react/examples/"+e.file())
				err := req.Send(nil)
				if err != nil {
					panic(err)
				}

				sources = sources.Set(i, e.setSrc(req.ResponseText))

				newSt := p.State()
				newSt.examples = sources
				p.SetState(newSt)
			}(i, e)
		}

		fetchStarted = true
	}
}

// GetInitialState returns in the initial state for the Examples component
func (p *ExamplesDef) GetInitialState() ExamplesState {
	return ExamplesState{
		examples:     sources,
		selectedTabs: newTabS(),
	}
}

// Render renders the Examples component
func (p *ExamplesDef) Render() r.Element {
	toRender := []r.Element{
		r.H3(nil, r.S("Reference")),
		r.P(nil, r.S("This entire page is a React application. An outer "), r.Code(nil, r.S("Examples")), r.S(" component contains a number of inner components.")),
		r.P(nil,
			r.S("For the source code, raising issues, questions etc, please see "),
			r.A(
				r.AProps(func(ap *r.APropsDef) {
					ap.Href = "https://github.com/myitcv/gopherjs/tree/master/react/examples"
					ap.Target = "_blank"
				}),
				r.S("the Github repo"),
			),
			r.S("."),
		),
		r.P(nil,
			r.S("Note the examples below show the GopherJS source code from "), r.Code(nil, r.S("master")),
		),

		p.renderExample(
			exampleHello,
			r.S("A Simple Example"),
			r.P(nil, r.S("The hellomessage.HelloMessage component demonstrates the simple use of a Props type.")),
			helloMessageJsx,
			hellomessage.HelloMessage(hellomessage.HelloMessageProps{Name: "Jane"}),
		),

		r.HR(nil),

		p.renderExample(
			exampleTimer,
			r.S("A Stateful Component"),
			r.P(nil, r.S("The timer.Timer component demonstrates the use of a State type.")),
			timerJsx,
			timer.Timer(),
		),

		r.HR(nil),

		p.renderExample(
			exampleTodo,
			r.S("An Application"),
			r.P(nil, r.S("The todoapp.TodoApp component demonstrates the use of state and event handling, but also the "+
				"problems of having a non-comparable state struct type.")),
			applicationJsx,
			todoapp.TodoApp(),
		),

		r.HR(nil),

		p.renderExample(
			exampleImmTodo,
			r.Span(nil, r.S("An Application using "), r.Code(nil, r.S("github.com/myitcv/immutable"))),
			r.P(nil, r.S("The immtodoapp.TodoApp component is a reimplementation of todoapp.TodoApp using immutable data structures.")),
			"n/a",
			immtodoapp.TodoApp(),
		),

		r.HR(nil),

		p.renderExample(
			exampleMarkdown,
			r.S("A Component Using External Plugins"),
			r.P(nil, r.S("The markdowneditor.MarkdownEditor component demonstrates the use of an external Javascript library.")),
			markdownEditorJsx,
			markdowneditor.MarkdownEditor(),
		),
	}

	return r.Div(
		r.DivProps(func(dp *r.DivPropsDef) {
			dp.ClassName = "container"
		}),

		toRender...,
	)
}

func (p *ExamplesDef) renderExample(key exampleKey, title, msg r.Element, jsxSrc string, elem r.Element) r.Element {

	var goSrc string
	src, _ := p.State().examples.Get(key)
	if src != nil {
		goSrc = src.src()
	}

	var code r.Element
	switch v, _ := p.State().selectedTabs.Get(key); v {
	case tabGo:
		code = r.Pre(nil, r.S(goSrc))
	case tabJsx:
		code = r.Pre(nil, r.S(jsxSrc))
	}

	return r.Div(nil,
		r.H3(nil, title),
		msg,
		r.Div(
			r.DivProps(func(dp *r.DivPropsDef) {
				dp.ClassName = "row"
			}),
			r.Div(
				r.DivProps(func(dp *r.DivPropsDef) {
					dp.ClassName = "col-md-8"
				}),
				r.Div(
					r.DivProps(func(dp *r.DivPropsDef) {
						dp.ClassName = "panel panel-default with-nav-tabs"
					}),
					r.Div(
						r.DivProps(func(dp *r.DivPropsDef) {
							dp.ClassName = "panel-heading"
						}),
						r.Ul(
							r.UlProps(func(ulp *r.UlPropsDef) {
								ulp.ClassName = "nav nav-tabs"
							}),

							p.buildExampleNavTab(key, tabGo, "GopherJS"),
							p.buildExampleNavTab(key, tabJsx, "JSX"),
						),
					),
					r.Div(
						r.DivProps(func(dp *r.DivPropsDef) {
							dp.ClassName = "panel-body"
						}),
						r.Pre(nil, code),
					),
				),
			),
			r.Div(
				r.DivProps(func(dp *r.DivPropsDef) {
					dp.ClassName = "col-md-4"
				}),
				plainPanel(elem),
			),
		),
	)
}

func (p *ExamplesDef) buildExampleNavTab(key exampleKey, t tab, title string) *r.LiDef {
	return r.Li(
		r.LiProps(func(lip *r.LiPropsDef) {
			if v, _ := p.State().selectedTabs.Get(key); v == t {
				lip.ClassName = "active"
			}
			lip.Role = "presentation"
		}),
		r.A(
			r.AProps(func(ap *r.APropsDef) {
				ap.Href = "#"
				ap.OnClick = p.handleTabChange(key, t)
			}),
			r.S(title),
		),
	)

}

func (p *ExamplesDef) handleTabChange(key exampleKey, t tab) func(*r.SyntheticMouseEvent) {
	return func(e *r.SyntheticMouseEvent) {
		cts := p.State().selectedTabs
		newSt := p.State()

		newSt.selectedTabs = cts.Set(key, t)
		p.SetState(newSt)

		e.PreventDefault()
	}
}

func plainPanel(children ...r.Element) r.Element {
	return r.Div(
		r.DivProps(func(dp *r.DivPropsDef) {
			dp.ClassName = "panel panel-default panel-body"
		}),
		children...,
	)
}
