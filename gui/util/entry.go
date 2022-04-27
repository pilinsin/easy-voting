package guiutil

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type IntEntry struct{
	*widget.Entry
	num string
}
func NewIntEntry() *IntEntry{
	ie := &IntEntry{}

	e := &widget.Entry{
		Wrapping: fyne.TextTruncate,
		PlaceHolder:     "0",
	}
	e.OnChanged = func(val string) {
		if val == ""{
			ie.num = val
			return
		}
		
		_, err := strconv.Atoi(val)
		if err == nil {
			ie.num = val
		}else{
			e.SetText(ie.num)
		}
	}
	
	ie.Entry = e
	ie.ExtendBaseWidget(ie)
	return ie
}
func (ie *IntEntry) Num() int{
	n, err := strconv.Atoi(ie.Text)
	if err != nil{return 0}
	return n
}