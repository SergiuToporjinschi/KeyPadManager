package widgets

import (
	"main/types"

	"fyne.io/fyne/v2/widget"
)

type ComboBox[T types.Displayable] struct {
	widget.SelectEntry
	values             []T
	selectedValue      T
	OnSelectionChanged func(T)
}

func NewCombobox[T types.Displayable](values []T) *ComboBox[T] {
	inst := &ComboBox[T]{
		values: values,
	}
	inst.ExtendBaseWidget(inst)
	inst.SetPlaceHolder("Select")
	inst.SetOptions(inst.getDisplayTexts())
	inst.OnChanged = inst.onChanged
	return inst
}

func (c *ComboBox[T]) GetSelectedValue() T {
	return c.selectedValue
}

func (c *ComboBox[T]) onChanged(value string) {
	for _, v := range c.values {
		if value == v.GetDisplayText() {
			c.selectedValue = v
			if c.OnSelectionChanged != nil {
				c.OnSelectionChanged(v)
			}
			break
		}
	}
}

func (c *ComboBox[T]) getDisplayTexts() []string {
	var displayTexts []string
	for _, v := range c.values {
		displayTexts = append(displayTexts, v.GetDisplayText())
	}
	return displayTexts
}
