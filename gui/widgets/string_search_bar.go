package widgets

import (
	resources "main/assets"
	"main/txt"

	"fyne.io/fyne/v2/widget"
)

type StringSearchBar[T any] struct {
	*widget.Entry
	*widget.Button
	searchList  []T
	foundIdexes []int
	cumpFunc    func(T, string) bool
	onFound     func(poz int)
	lastSearch  string
	lastRetPoz  int
}

func NewStringSearchBar[T any](searchList []T, compFunc func(T, string) bool, onFound func(poz int)) *StringSearchBar[T] {
	entr := widget.NewEntry()
	entr.SetPlaceHolder(txt.GetLabel("btn.search"))

	inst := &StringSearchBar[T]{
		Entry:      entr,
		cumpFunc:   compFunc,
		onFound:    onFound,
		searchList: searchList,
	}
	inst.Button = widget.NewButtonWithIcon(
		txt.GetLabel("btn.search"),
		resources.ResSearchPng,
		inst.onBtnPush,
	)

	return inst
}

func (sb *StringSearchBar[T]) GetControls() (*widget.Entry, *widget.Button) {
	return sb.Entry, sb.Button
}

func (sb *StringSearchBar[T]) onBtnPush() {
	if sb.Entry.Text == "" {
		sb.foundIdexes = sb.foundIdexes[:0]
		sb.lastRetPoz = 0
		return
	}
	if sb.Entry.Text == sb.lastSearch {
		sb.lastRetPoz++
		if sb.lastRetPoz >= len(sb.foundIdexes) {
			sb.lastRetPoz = 0
		}
		sb.onFound(sb.foundIdexes[sb.lastRetPoz])
		return
	}
	for i, val := range sb.searchList {
		if sb.cumpFunc(val, sb.Entry.Text) {
			sb.foundIdexes = append(sb.foundIdexes, i)
		}
	}
	if len(sb.foundIdexes) > 0 {
		sb.lastRetPoz = 0
		sb.onFound(sb.foundIdexes[sb.lastRetPoz])
		sb.lastSearch = sb.Entry.Text
	}
}

func (sb *StringSearchBar[T]) SetSearchList(searchList []T) {
	sb.searchList = searchList
	sb.Reset()
}

func (sb *StringSearchBar[T]) Reset() {
	sb.foundIdexes = sb.foundIdexes[:0]
	sb.lastRetPoz = 0
	sb.lastSearch = ""
	sb.Entry.SetText("")
}
