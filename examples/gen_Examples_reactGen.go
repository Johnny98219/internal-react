// Code generated by reactGen. DO NOT EDIT.

package examples

import "myitcv.io/react"

type ExamplesElem struct {
	react.Element
}

func (e ExamplesDef) ShouldComponentUpdateIntf(nextProps, prevState, nextState interface{}) bool {
	res := false

	v := prevState.(ExamplesState)
	res = !v.EqualsIntf(nextState) || res
	return res
}

func buildExamples(cd react.ComponentDef) react.Component {
	return ExamplesDef{ComponentDef: cd}
}

func buildExamplesElem(children ...react.Element) *ExamplesElem {
	return &ExamplesElem{
		Element: react.CreateElement(buildExamples, nil),
	}
}

// SetState is an auto-generated proxy proxy to update the state for the
// Examples component.  SetState does not immediately mutate e.State()
// but creates a pending state transition.
func (e ExamplesDef) SetState(state ExamplesState) {
	e.ComponentDef.SetState(state)
}

// State is an auto-generated proxy to return the current state in use for the
// render of the Examples component
func (e ExamplesDef) State() ExamplesState {
	return e.ComponentDef.State().(ExamplesState)
}

// IsState is an auto-generated definition so that ExamplesState implements
// the myitcv.io/react.State interface.
func (e ExamplesState) IsState() {}

var _ react.State = ExamplesState{}

// GetInitialStateIntf is an auto-generated proxy to GetInitialState
func (e ExamplesDef) GetInitialStateIntf() react.State {
	return e.GetInitialState()
}

func (e ExamplesState) EqualsIntf(val interface{}) bool {
	return e == val.(ExamplesState)
}
