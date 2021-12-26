package registrationpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	user "EasyVoting/registration"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
)

func errorPage(err error) fyne.CanvasObject {
	return widget.NewLabel(fmt.Sprintln(err))
}

type registrationPage struct {
	fyne.CanvasObject
	r user.IRegistration
}

func LoadPage(rCfgCid string, is *ipfs.IPFS) fyne.CanvasObject {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return errorPage(err)
	}
	r, err := user.NewRegistration(rCfgCid, is)
	if err != nil {
		return errorPage(err)
	}

	titleLabel := widget.NewLabel(rCfg.Title() + " (" + rCfgCid + ")")
	noteLabel := widget.NewLabel("")

	rForm := registrationForm(rCfg.UserDataLabels(), r, noteLabel)

	//check icon
	checkBtn := widget.NewButtonWithIcon("verify HashNameMap", theme.UploadIcon(), func() {
		noteLabel.Text = "processing..."
		if ok := r.VerifyHashNameMap(); ok {
			noteLabel.Text = "HashNameMap is verified"
		} else {
			noteLabel.Text = "invalid HashNameMap"
		}
	})

	page := container.NewVBox(rForm, checkBtn, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return &registrationPage{page, r}
}

func registrationForm(labels []string, r user.IRegistration, noteLabel *widget.Label) fyne.CanvasObject {
	rForm := &widget.Form{}
	entries := make([]*widget.Entry, len(labels))
	for idx, label := range labels {
		entries[idx] = widget.NewEntry()
		formItem := widget.NewFormItem(label, entries[idx])
		rForm.Items = append(rForm.Items, formItem)
	}
	rForm.OnSubmit = func() {
		noteLabel.Text = "processing..."
		var texts []string
		for _, entry := range entries {
			texts = append(texts, entry.Text)
		}
		identity, err := r.Registrate(texts...)
		if err != nil {
			noteLabel.Text = fmt.Sprintln("registration error:", err)
		} else {
			idStr := util.Bytes64ToStr(identity.Marshal())
			noteLabel.Text = fmt.Sprintln("registration is done, userIdentity:", idStr)
		}
	}
	rForm.ExtendBaseWidget(rForm)
	return rForm
}
