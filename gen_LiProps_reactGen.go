// Code generated by reactGen. DO NOT EDIT.

package react

// LiProps defines the properties for the <li> element
type LiProps struct {
	ClassName               string
	DangerouslySetInnerHTML *DangerousInnerHTML
	DataSet
	ID  string
	Key string

	OnChange
	OnClick

	Ref
	Role  string
	Style *CSS
}

func (l *LiProps) assign(v *_LiProps) {

	v.ClassName = l.ClassName

	v.DangerouslySetInnerHTML = l.DangerouslySetInnerHTML

	if l.DataSet != nil {
		for dk, dv := range l.DataSet {
			v.o.Set("data-"+dk, dv)
		}
	}

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
