package registrationpageutil

import (
	"encoding/csv"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

)


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
func CsvDialog(w fyne.Window, csvMat chan<- []string, hideBtn fyne.CanvasObject, noteLabel *widget.Label) func(){
	return func(){
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if rc == nil || err != nil{
				noteLabel.SetText("no file is selected")
				return
			}
			if rc.URI().Extension() != ".csv" {
				noteLabel.SetText("invalid file is selected")
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
			hideBtn.Hide()
			noteLabel.SetText("csv file uploaded")
		}
		dialog.ShowFileOpen(onSelected, w)
	}
}
