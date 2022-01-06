package registrationpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
)

func NewSetupPage(w fyne.Window, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	cidEntry := widget.NewEntry()
	cidEntry.SetPlaceHolder("registration manager cid")

	titleEntry := widget.NewEntry()

	userDataset := make(chan []string)
	var err error
	//load icon
	loadCsvBtn := &widget.Button{
		Text: "upload userDataset csv",
		Icon: theme.UploadIcon(),
	}
	loadCsvBtn.OnTapped = func() {
		dialogErr := rputil.CsvDialog(w, userDataset)
		if dialogErr == nil {
			loadCsvBtn.Hide()
			noteLabel.Text = "csv file uploaded"
		} else {
			noteLabel.Text = "invalid csv file"
		}
		err = dialogErr
	}
	loadCsvBtn.ExtendBaseWidget(loadCsvBtn)

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", titleEntry))
	form.Items = append(form.Items, widget.NewFormItem("csv", loadCsvBtn))
	form.OnSubmit = func() {
		if titleEntry.Text == "" {
			noteLabel.SetText("title is empty")
			return
		}
		if err != nil {
			noteLabel.SetText("invalid csv file")
			return
		}

		noteLabel.SetText("processing...")
		userDataLabels := <-userDataset
		mCfg, rCfg := rutil.NewConfigs(titleEntry.Text, userDataset, userDataLabels, is)
		mCfgCid := ipfs.ToCidWithAdd(mCfg.Marshal(), is)
		ipfs.ToCidWithAdd(rCfg.Marshal(), is)
		noteLabel.SetText("registration manager cid: ")
		cidEntry.SetText(mCfgCid)

		form.Hide()
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel, cidEntry)
	return page
}
