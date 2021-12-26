package registrationpage

import (
	"encoding/csv"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
)

type setupPage struct {
	fyne.CanvasObject
	is             *ipfs.IPFS
	userDataset    chan []string
	userDataLabels []string
	err            error
}

func NewSetupPage(a fyne.App, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")

	titleEntry := widget.NewEntry()

	userDataset := make(chan []string)
	var userDataLabels []string
	var err error
	//load icon
	loadCsvBtn := widget.NewButtonWithIcon(
		"load userDataset csv",
		theme.UploadIcon(),
		csvDialog(a, userDataset, &userDataLabels, err),
	)

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", titleEntry))
	form.Items = append(form.Items, widget.NewFormItem("csv", loadCsvBtn))
	form.OnSubmit = func() {
		if err != nil {
			noteLabel.Text = "invalid csv file"
			return
		}

		noteLabel.Text = "processing..."
		mCfg, rCfg := rutil.NewConfigs(titleEntry.Text, userDataset, userDataLabels, is)
		mCfgCid := ipfs.ToCidWithAdd(mCfg.Marshal(), is)
		ipfs.ToCidWithAdd(rCfg.Marshal(), is)
		noteLabel.Text = "registration manager cid: " + mCfgCid
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel)
	return page //&setupPage{page, is, userDataset, userDataLabels, err}
}

/*
	{"Label 1, Label 2, ..., Label M"}
	{"data1 11, data1 12, ..., data 1M"}
	...
	{"data N1, data N2, ..., data NM"}
*/
func csvDialog(a fyne.App, csvMat chan<- []string, csvLabels *[]string, csvErr error) func() {
	return func() {
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if err != nil {
				csvErr = err
				return
			}
			reader := csv.NewReader(rc)
			labels, err := reader.Read()
			if err != nil {
				csvErr = err
				return
			}
			csvLabels = &labels
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
		fOpenWin := a.NewWindow("Open a csv file")
		dialog.NewFileOpen(onSelected, fOpenWin)
		fOpenWin.ShowAndRun()
	}
}
