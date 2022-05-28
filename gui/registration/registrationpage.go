package registrationpage

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	riface "github.com/pilinsin/easy-voting/registration/interface"
)

func LoadPage(ctx context.Context, bAddr, rCfgCid string, r riface.IRegistration) (string, fyne.CanvasObject) {
	uidLabel := gutil.NewCopyButton("user identity address")

	rCfg := r.Config()
	bAddrLabel := gutil.NewCopyButton(bAddr)
	cfgLabel := gutil.NewCopyButton(rCfgCid)
	addrLabel := container.NewVBox(bAddrLabel.Render(), cfgLabel.Render())
	titleLabel := widget.NewLabel(rCfg.Title)
	noteLabel := widget.NewLabel("")

	entries := make([]*widget.Entry, len(rCfg.Labels))
	rForm := &widget.Form{}
	for idx, label := range rCfg.Labels {
		entries[idx] = widget.NewEntry()
		rForm.Items = append(rForm.Items, widget.NewFormItem(label, entries[idx]))
	}
	rForm.OnSubmit = func() {
		noteLabel.SetText("processing...")
		dataset := make([]string, len(rCfg.Labels))
		for idx, entry := range entries {
			dataset[idx] = entry.Text
		}
		uidStr, err := r.Registrate(dataset...)
		if err != nil {
			noteLabel.SetText("registration error: " + err.Error())
			return
		}
		noteLabel.SetText("done")
		uidLabel.SetText(uidStr)
	}
	rForm.ExtendBaseWidget(rForm)

	titles := container.NewVBox(titleLabel, addrLabel)
	page := container.NewVBox(rForm, noteLabel, uidLabel.Render())
	page = container.NewBorder(titles, nil, nil, nil, page)
	return rCfg.Title, page
}
