// Code generated by reactGen; DO NOT EDIT.

package main

import "github.com/myitcv/gopherjs/react"

func (a *AppDef) ShouldComponentUpdateIntf(nextProps, prevState, nextState interface{}) bool {
	res := false

	v := prevState.(AppState)
	res = !v.EqualsIntf(nextState) || res
	return res
}

// SetState is an auto-generated proxy proxy to update the state for the
// App component.  SetState does not immediately mutate a.State()
// but creates a pending state transition.
func (a *AppDef) SetState(s AppState) {
	a.ComponentDef.SetState(s)
}

// State is an auto-generated proxy to return the current state in use for the
// render of the App component
func (a *AppDef) State() AppState {
	return a.ComponentDef.State().(AppState)
}

// IsState is an auto-generated definition so that AppState implements
// the github.com/myitcv/gopherjs/react.State interface.
func (a AppState) IsState() {}

var _ react.State = AppState{}

// GetInitialStateIntf is an auto-generated proxy to GetInitialState
func (a *AppDef) GetInitialStateIntf() react.State {
	return AppState{}
}

func (a AppState) EqualsIntf(v interface{}) bool {
	return a == v.(AppState)
}
