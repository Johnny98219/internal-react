// Code generated by reactGen. DO NOT EDIT.

package react

// PProps are the props for a <div> component
type PProps struct {
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

func (p *PProps) assign(v *_PProps) {

	v.ClassName = p.ClassName

	v.DangerouslySetInnerHTML = p.DangerouslySetInnerHTML

	if p.DataSet != nil {
		for dk, dv := range p.DataSet {
			v.o.Set("data-"+dk, dv)
		}
	}

	if p.ID != "" {
		v.ID = p.ID
	}

	if p.Key != "" {
		v.Key = p.Key
	}

	if p.OnChange != nil {
		v.o.Set("onChange", p.OnChange.OnChange)
	}

	if p.OnClick != nil {
		v.o.Set("onClick", p.OnClick.OnClick)
	}

	if p.Ref != nil {
		v.o.Set("ref", p.Ref.Ref)
	}

	v.Role = p.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = p.Style.hack()

}
