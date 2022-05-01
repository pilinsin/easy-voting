package guiutil

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	pv "github.com/pilinsin/p2p-verse"
)

type RemovableEntryForm struct {
	entries map[string]*widget.Entry
}

func NewRemovableEntryForm() *RemovableEntryForm {
	es := make(map[string]*widget.Entry)
	return &RemovableEntryForm{es}
}
func (ref *RemovableEntryForm) Render() fyne.CanvasObject {
	contents := container.NewVBox()
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		addrEntry := widget.NewEntry()
		id := pv.RandString(16)
		ref.entries[id] = addrEntry

		rmvBtn := &widget.Button{Icon: theme.ContentClearIcon()}
		withRmvBtn := container.NewBorder(nil, nil, nil, rmvBtn, addrEntry)
		rmvBtn.OnTapped = func() {
			contents.Remove(withRmvBtn)
			delete(ref.entries, id)
		}
		rmvBtn.ExtendBaseWidget(rmvBtn)
		contents.Add(withRmvBtn)
	})
	addBtnObj := container.NewBorder(nil, nil, addBtn, nil)
	return container.NewBorder(addBtnObj, nil, nil, nil, contents)
}

func (ref *RemovableEntryForm) Texts() []string {
	txts := make([]string, len(ref.entries))
	idx := 0
	for _, entry := range ref.entries {
		txts[idx] = entry.Text
		idx++
	}
	return txts
}
