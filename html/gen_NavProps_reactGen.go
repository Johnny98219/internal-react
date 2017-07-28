// Code generated by reactGen. DO NOT EDIT.

package html

import "myitcv.io/react/dom"

// NavProps defines the properties for the <nav> element
type NavProps struct {
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTML
	ID                      string
	Key                     string
	OnChange                dom.OnChange
	OnClick                 dom.OnClick
	Role                    string
	Style                   *CSS
}

func (n *NavProps) assign(v *_NavProps) {

	v.ClassName = n.ClassName

	v.DangerouslySetInnerHTML = n.DangerouslySetInnerHTML

	if n.ID != "" {
		v.ID = n.ID
	}

	if n.Key != "" {
		v.Key = n.Key
	}

	v.OnChange = n.OnChange

	v.OnClick = n.OnClick

	v.Role = n.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = n.Style.hack()

}
