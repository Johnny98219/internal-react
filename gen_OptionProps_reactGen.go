// Code generated by reactGen. DO NOT EDIT.

package react

// OptionProps defines the properties for the <option> element
type OptionProps struct {
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
	Value string
}

func (o *OptionProps) assign(v *_OptionProps) {

	v.ClassName = o.ClassName

	v.DangerouslySetInnerHTML = o.DangerouslySetInnerHTML

	if o.DataSet != nil {
		for dk, dv := range o.DataSet {
			v.o.Set("data-"+dk, dv)
		}
	}

	if o.ID != "" {
		v.ID = o.ID
	}

	if o.Key != "" {
		v.Key = o.Key
	}

	if o.OnChange != nil {
		v.o.Set("onChange", o.OnChange.OnChange)
	}

	if o.OnClick != nil {
		v.o.Set("onClick", o.OnClick.OnClick)
	}

	if o.Ref != nil {
		v.o.Set("ref", o.Ref.Ref)
	}

	v.Role = o.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	v.Style = o.Style.hack()

	v.Value = o.Value

}
