// Code generated by reactGen. DO NOT EDIT.

package react

// LabelProps defines the properties for the <label> element
type LabelProps struct {
	AriaSet
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTML
	DataSet
	For string
	ID  string
	Key string

	OnChange
	OnClick

	Ref
	Role  string
	Style *CSS
}

func (l *LabelProps) assign(v *_LabelProps) {

	if l.AriaSet != nil {
		for dk, dv := range l.AriaSet {
			v.o.Set("aria-"+dk, dv)
		}
	}

	v.ClassName = l.ClassName

	v.DangerouslySetInnerHTML = l.DangerouslySetInnerHTML

	if l.DataSet != nil {
		for dk, dv := range l.DataSet {
			v.o.Set("data-"+dk, dv)
		}
	}

	v.For = l.For

	if l.ID != "" {
		v.ID = l.ID
	}

	if l.Key != "" {
		v.Key = l.Key
	}

	if l.OnChange != nil {
		v.o.Set("onChange", l.OnChange.OnChange)
	}

	if l.OnClick != nil {
		v.o.Set("onClick", l.OnClick.OnClick)
	}

	if l.Ref != nil {
		v.o.Set("ref", l.Ref.Ref)
	}

	v.Role = l.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = l.Style.hack()

}
