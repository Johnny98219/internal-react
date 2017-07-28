// Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
// Use of this document is governed by a license found in the LICENSE document.

package html

import (
	"myitcv.io/react"
)

// SpanElem is the React element definition corresponding to the HTML <p> element
type SpanElem struct {
	react.Element
}

// _SpanProps defines the properties for the <p> element
type _SpanProps struct {
	*BasicHTMLElement
}

// Span creates a new instance of a <p> element with the provided props and
// children
func Span(props *SpanProps, children ...react.Element) *SpanElem {

	rProps := &_SpanProps{
		BasicHTMLElement: newBasicHTMLElement(),
	}

	if props != nil {
		props.assign(rProps)
	}

	return &SpanElem{
		Element: react.InternalCreateElement("span", rProps, children...),
	}
}
