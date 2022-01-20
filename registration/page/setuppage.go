package registrationpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/pilinsin/easy-voting/ipfs"
	rputil "github.com/pilinsin/easy-voting/registration/page/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

func NewSetupPage(w fyne.Window, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	cidEntry := widget.NewEntry()
	cidEntry.SetPlaceHolder("registration config cid")

	titleEntry := widget.NewEntry()

	userDataset := make(chan []string)
	//load icon
	loadCsvBtn := &widget.Button{
		Text: "upload userDataset csv",
		Icon: theme.UploadIcon(),
	}
	loadCsvBtn.OnTapped = rputil.CsvDialog(w, userDataset, loadCsvBtn, noteLabel)
	loadCsvBtn.ExtendBaseWidget(loadCsvBtn)

	kwEntry := widget.NewEntry()
	kwEntry.SetPlaceHolder("keyword of registration manager identity")

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", titleEntry))
	form.Items = append(form.Items, widget.NewFormItem("csv", loadCsvBtn))
	form.Items = append(form.Items, widget.NewFormItem("keyword", kwEntry))
	form.OnSubmit = func() {
		if titleEntry.Text == "" {
			noteLabel.SetText("title is empty")
			return
		}
		if loadCsvBtn.Visible(){
			noteLabel.SetText("csv is empty")
			return
		}

		noteLabel.SetText("processing...")
		userDataLabels := <-userDataset
		mId, rCfg := rutil.NewConfigs(titleEntry.Text, userDataset, userDataLabels, is)
		rCfgCid := ipfs.ToCidWithAdd(rCfg.Marshal(), is)
		noteLabel.SetText("registration config cid: ")
		cidEntry.SetText(rCfgCid)

		idStore := rutil.NewIdentityStore()
		idStore.Put(kwEntry.Text, mId.Marshal())
		idStore.Close()

		form.Hide()
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel, cidEntry)
	return page
}
