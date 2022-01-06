package registrationpageutil

import (
	"encoding/csv"
	"fmt"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/util"
)

func ErrorPage(err error) fyne.CanvasObject {
	return widget.NewLabel(fmt.Sprintln(err))
}

type IPageCloser interface {
	//only cancel goroutine in a page
	Cancel()
	SetCancel(func())
	//cancel & (registration/manager).Close()
	Close()
}
type pageCloser struct {
	close, cancel func()
}

func NewPageCloser(close, cancel func()) *pageCloser {
	return &pageCloser{close, cancel}
}
func (c *pageCloser) SetCancel(cancel func()) {
	c.cancel = cancel
}
func (c *pageCloser) Cancel() {
	c.cancel()
}
func (c *pageCloser) Close() {
	c.cancel()
	c.close()
}

/*
	{"Label 1, Label 2, ..., Label M"}
	{"data1 11, data1 12, ..., data 1M"}
	...
	{"data N1, data N2, ..., data NM"}
*/
func CsvDialog(w fyne.Window, csvMat chan<- []string) error {
	var dialogErr error
	onSelected := func(rc fyne.URIReadCloser, err error) {
		if rc == nil {
			dialogErr = util.NewError("no file is selected")
			return
		}
		if err != nil {
			dialogErr = err
			return
		}

		reader := csv.NewReader(rc)
		go func() {
			defer close(csvMat)
			for {
				data, err := reader.Read()
				if err == io.EOF {
					return
				}
				if err == nil {
					csvMat <- data
				}
			}
		}()
	}
	dialog.ShowFileOpen(onSelected, w)
	return dialogErr
}
