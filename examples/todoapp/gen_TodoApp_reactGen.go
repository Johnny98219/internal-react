// Code generated by reactGen. DO NOT EDIT.

package todoapp

import "myitcv.io/react"

type TodoAppElem struct {
	react.Element
}

func buildTodoApp(cd react.ComponentDef) react.Component {
	return TodoAppDef{ComponentDef: cd}
}

func buildTodoAppElem(children ...react.Element) *TodoAppElem {
	return &TodoAppElem{
		Element: react.CreateElement(buildTodoApp, nil, children...),
	}
}

func (t TodoAppDef) RendersElement() react.Element {
	return t.Render()
}

// SetState is an auto-generated proxy proxy to update the state for the
// TodoApp component.  SetState does not immediately mutate t.State()
// but creates a pending state transition.
func (t TodoAppDef) SetState(state TodoAppState) {
	t.ComponentDef.SetState(state)
}

// State is an auto-generated proxy to return the current state in use for the
// render of the TodoApp component
func (t TodoAppDef) State() TodoAppState {
	return t.ComponentDef.State().(TodoAppState)
}

// IsState is an auto-generated definition so that TodoAppState implements
// the myitcv.io/react.State interface.
func (t TodoAppState) IsState() {}

var _ react.State = TodoAppState{}

// GetInitialStateIntf is an auto-generated proxy to GetInitialState
func (t TodoAppDef) GetInitialStateIntf() react.State {
	return TodoAppState{}
}

func (t TodoAppState) EqualsIntf(val react.State) bool {
	return t.Equals(val.(TodoAppState))
}
