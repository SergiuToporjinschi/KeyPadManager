package widgets

import (
	"sync/atomic"

	"fyne.io/fyne/v2/data/binding"
)

type Binder struct {
	callback atomic.Value
}

func (binder *Binder) Bind(data binding.DataItem, f func(data binding.DataItem)) {
	listener := binding.NewDataListener(func() {
		f := binder.callback.Load()
		if fn, ok := f.(func(binding.DataItem)); ok && fn != nil {
			fn(data)
		}
	})
	data.AddListener(listener)
	binder.callback.Store(f)
}
