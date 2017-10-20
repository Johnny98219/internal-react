// Code generated by reactGen. DO NOT EDIT.

package hellomessage

import "myitcv.io/react"

type HelloMessageElem struct {
	react.Element
}

func (h HelloMessageDef) ShouldComponentUpdateIntf(nextProps react.Props, prevState, nextState react.State) bool {
	res := false

	{
		res = h.Props() != nextProps.(HelloMessageProps) || res
	}
	return res
}

func buildHelloMessage(cd react.ComponentDef) react.Component {
	return HelloMessageDef{ComponentDef: cd}
}

func buildHelloMessageElem(props HelloMessageProps, children ...react.Element) *HelloMessageElem {
	return &HelloMessageElem{
		Element: react.CreateElement(buildHelloMessage, props, children...),
	}
}

// IsProps is an auto-generated definition so that HelloMessageProps implements
// the myitcv.io/react.Props interface.
func (h HelloMessageProps) IsProps() {}

// Props is an auto-generated proxy to the current props of HelloMessage
func (h HelloMessageDef) Props() HelloMessageProps {
	uprops := h.ComponentDef.Props()
	return uprops.(HelloMessageProps)
}

func (h HelloMessageProps) EqualsIntf(val react.Props) bool {
	return h == val.(HelloMessageProps)
}

var _ react.Props = HelloMessageProps{}
