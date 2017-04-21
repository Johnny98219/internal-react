// Code generated by reactGen. DO NOT EDIT.

package react

// TextAreaProps defines the properties for the <textarea> element
type TextAreaProps struct {
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTMLDef
	DefaultValue            string
	ID                      string
	Key                     string
	OnChange                func(e *SyntheticEvent)
	OnClick                 func(e *SyntheticMouseEvent)
	Placeholder             string
	Role                    string
	Style                   *CSS
	Value                   string
}

func (t *TextAreaProps) assign(v *_TextAreaProps) {

	v.ClassName = t.ClassName

	v.DangerouslySetInnerHTML = t.DangerouslySetInnerHTML

	if t.DefaultValue != "" {
		v.DefaultValue = t.DefaultValue
	}

	if t.ID != "" {
		v.ID = t.ID
	}

	if t.Key != "" {
		v.Key = t.Key
	}

	v.OnChange = t.OnChange

	v.OnClick = t.OnClick

	v.Placeholder = t.Placeholder

	v.Role = t.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = t.Style.hack()

	v.Value = t.Value

}
