package guiutil

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CsvDialog(w fyne.Window, rc0 chan *csv.Reader, noteLabel *widget.Label) func() {
	return func() {
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if rc == nil || err != nil {
				noteLabel.SetText("no file is selected")
				return
			}
			if rc.URI().Extension() != ".csv" {
				noteLabel.SetText("invalid file is selected")
				return
			}

			go func() {
				rc0 <- csv.NewReader(rc)
				noteLabel.SetText("csv file uploaded")
			}()
		}
		dialog.ShowFileOpen(onSelected, w)
	}
}

type loadCsvButton struct {
	*widget.Button
	reader chan *csv.Reader
}

func NewLoadCsvButton(w fyne.Window, noteLabel *widget.Label) *loadCsvButton {
	r := make(chan *csv.Reader, 1)

	onTapped := CsvDialog(w, r, noteLabel)
	btn := widget.NewButtonWithIcon("upload csv", theme.UploadIcon(), onTapped)

	lcb := &loadCsvButton{btn, r}
	lcb.ExtendBaseWidget(lcb)
	return lcb
}

/*
	{"Label 1, Label 2, ..., Label M"}
	{"data1 11, data1 12, ..., data 1M"}
	...
	{"data N1, data N2, ..., data NM"}
*/
func (lcb *loadCsvButton) Read() ([]string, <-chan []string, error) {
	var r *csv.Reader
	var isReaderGet bool
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	select {
	case r = <-lcb.reader:
		isReaderGet = true
	case <-ctx.Done():
		//close(ce.thumbnail)
	}
	if !isReaderGet {
		return nil, nil, errors.New("reader is nil")
	}

	labels, err := r.Read()
	if err == io.EOF {
		return labels, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan []string)
	go func() {
		defer close(ch)
		for {
			data, err := r.Read()
			if err == io.EOF {
				return
			}
			if err == nil {
				ch <- data
			}
		}
	}()
	return labels, ch, nil
}
