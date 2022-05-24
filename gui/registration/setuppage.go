package registrationpage

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	riface "github.com/pilinsin/easy-voting/registration/interface"
	rgst "github.com/pilinsin/easy-voting/registration"
)

func NewSetupPage(w fyne.Window) fyne.CanvasObject {
	var r riface.IRegistration

	noteLabel := widget.NewLabel("")
	addrLabel := gutil.NewCopyButton("registration config address")

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
		addrLabel.SetText("registration config address")
		labels, dataset, err := csvBtn.Read()
		if err != nil {
			noteLabel.SetText("load csv error: " + err.Error())
			return
		}
		cid, baseDir, err := rutil.NewConfig(titleEntry.Text, dataset, labels, bAddrEntry.Text)
		if err != nil {
			noteLabel.SetText("new rConfig error: " + err.Error())
			return
		}
		r, err = rgst.NewRegistration(context.Background(), cid, baseDir)
		if err != nil {
			noteLabel.SetText("new rConfig error: " + err.Error())
			return
		}
		noteLabel.SetText("done")
		addrLabel.SetText(cid)
		//form.Hide()
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel, addrLabel.Render())
	return page
}
