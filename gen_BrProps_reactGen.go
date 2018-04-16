// Code generated by reactGen. DO NOT EDIT.

package react

// BrProps defines the properties for the <br> element
type BrProps struct {
	AriaExpanded            bool
	AriaHasPopup            bool
	AriaLabelledBy          string
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

func (b *BrProps) assign(_v *_BrProps) {

	_v.AriaExpanded = b.AriaExpanded

	_v.AriaHasPopup = b.AriaHasPopup

	_v.AriaLabelledBy = b.AriaLabelledBy

	_v.ClassName = b.ClassName

	_v.DangerouslySetInnerHTML = b.DangerouslySetInnerHTML

	if b.DataSet != nil {
		for dk, dv := range b.DataSet {
			_v.o.Set("data-"+dk, dv)
		}
	}

	if b.ID != "" {
		_v.ID = b.ID
	}

	if b.Key != "" {
		_v.Key = b.Key
	}

	if b.OnChange != nil {
		_v.o.Set("onChange", b.OnChange.OnChange)
	}

	if b.OnClick != nil {
		_v.o.Set("onClick", b.OnClick.OnClick)
	}

	if b.Ref != nil {
		_v.o.Set("ref", b.Ref.Ref)
	}

	_v.Role = b.Role

	// TODO: until we have a resolution on
	// https://github.com/gopherjs/gopherjs/issues/236
	_v.Style = b.Style.hack()

}
