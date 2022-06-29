package guiutil

import (
	"bytes"
	"encoding/csv"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func CsvDialog(w fyne.Window, lcb *loadCsvButton, noteLabel *widget.Label) func() {
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

			lcb.reader = csv.NewReader(rc)
			noteLabel.SetText("csv file uploaded")

		}
		dialog.ShowFileOpen(onSelected, w)
	}
}

type loadCsvButton struct {
	*widget.Button
	reader *csv.Reader
}

func NewLoadCsvButton(w fyne.Window, noteLabel *widget.Label) *loadCsvButton {
	lcb := &loadCsvButton{}

	onTapped := CsvDialog(w, lcb, noteLabel)
	lcb.Button = widget.NewButtonWithIcon("upload csv", theme.UploadIcon(), onTapped)
	lcb.ExtendBaseWidget(lcb)
	return lcb
}

/*
	{"Label 1, Label 2, ..., Label M"}
	{"data1 11, data1 12, ..., data 1M"}
	...
	{"data N1, data N2, ..., data NM"}
*/
func (lcb *loadCsvButton) Csv() ([]string, <-chan []string, error) {
	buf := &bytes.Buffer{}
	w := csv.NewWriter(buf)

	labels, err := lcb.reader.Read()
	if err != nil {
		return nil, nil, err
	}
	if err := w.Write(labels); err != nil {
		return nil, nil, err
	}

	ch := make(chan []string)
	go func() {
		defer close(ch)
		for {
			data, err := lcb.reader.Read()
			if err == io.EOF {
				w.Flush()
				lcb.reader = csv.NewReader(buf)
				return
			}
			if err == nil {
				ch <- data
				if err := w.Write(data); err != nil {
					return
				}
			}
		}
	}()

	return labels, ch, nil
}
