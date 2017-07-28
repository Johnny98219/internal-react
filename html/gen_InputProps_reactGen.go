// Code generated by reactGen. DO NOT EDIT.

package html

import "myitcv.io/react/dom"

// InputProps defines the properties for the <input> element
type InputProps struct {
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTML
	DefaultValue            string
	ID                      string
	Key                     string
	OnChange                dom.OnChange
	OnClick                 dom.OnClick
	Placeholder             string
	Role                    string
	Style                   *CSS
	Type                    string
	Value                   string
}

func (i *InputProps) assign(v *_InputProps) {

	v.ClassName = i.ClassName

	v.DangerouslySetInnerHTML = i.DangerouslySetInnerHTML

	if i.DefaultValue != "" {
		v.DefaultValue = i.DefaultValue
	}

	if i.ID != "" {
		v.ID = i.ID
	}

	if i.Key != "" {
		v.Key = i.Key
	}

	v.OnChange = i.OnChange

	v.OnClick = i.OnClick

	v.Placeholder = i.Placeholder

	v.Role = i.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = i.Style.hack()

	v.Type = i.Type

	v.Value = i.Value

}
