// Code generated by reactGen. DO NOT EDIT.

package react

// NavProps defines the properties for the <nav> element
type NavProps struct {
	AriaSet
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

func (n *NavProps) assign(_v *_NavProps) {

	if n.AriaSet != nil {
		for dk, dv := range n.AriaSet {
			_v.o.Set("aria-"+dk, dv)
		}
	}

	_v.ClassName = n.ClassName

	_v.DangerouslySetInnerHTML = n.DangerouslySetInnerHTML

	if n.DataSet != nil {
		for dk, dv := range n.DataSet {
			_v.o.Set("data-"+dk, dv)
		}
	}

	if n.ID != "" {
		_v.ID = n.ID
	}

	if n.Key != "" {
		_v.Key = n.Key
	}

	if n.OnChange != nil {
		_v.o.Set("onChange", n.OnChange.OnChange)
	}

	if n.OnClick != nil {
		_v.o.Set("onClick", n.OnClick.OnClick)
	}

	if n.Ref != nil {
		_v.o.Set("ref", n.Ref.Ref)
	}

	_v.Role = n.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	_v.Style = n.Style.hack()

}
