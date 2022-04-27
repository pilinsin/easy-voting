package registrationpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

func NewSetupPage(w fyne.Window) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	addrLabel := gutil.NewCopyButton("registration config address")
	maIdLabel := gutil.NewCopyButton("manager identity address")

	titleEntry := widget.NewEntry()
	csvBtn := gutil.NewLoadCsvButton(w, noteLabel)
	bAddrEntry := widget.NewEntry()

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", titleEntry))
	form.Items = append(form.Items, widget.NewFormItem("csv", csvBtn))
	form.Items = append(form.Items, widget.NewFormItem("bAddr", bAddrEntry))
	form.OnSubmit = func() {
		if titleEntry.Text == "" {
			noteLabel.SetText("title is empty")
			return
		}
		if bAddrEntry.Text == "" {
			noteLabel.SetText("bAddr is empty")
			return
		}

		noteLabel.SetText("processing...")
		labels, dataset, err := csvBtn.Read()
		if err != nil{
			noteLabel.SetText("load csv error: "+err.Error())
			return
		}
		cid, mid, err := rutil.NewConfig(titleEntry.Text, dataset, labels, bAddrEntry.Text)
		if err != nil{
			noteLabel.SetText("new rConfig error: "+err.Error())
			return
		}

		noteLabel.SetText("done")
		addrLabel.SetText(cid)
		maIdLabel.SetText(mid)

		//form.Hide()
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel, addrLabel.Render(), maIdLabel.Render())
	return page
}
