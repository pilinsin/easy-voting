package registrationpage

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	rgst "github.com/pilinsin/easy-voting/registration"
	gutil "github.com/pilinsin/easy-voting/gui/util"
)

func LoadPage(ctx context.Context, rCfgAddr, idStr string) (string, fyne.CanvasObject, func()) {
	uidLabel := gutil.NewCopyButton("user identity address")

	r, err := rgst.NewRegistration(ctx, rCfgAddr, idStr)
	if err != nil{return "", nil, nil}
	closer := func(){r.Close()}

	rCfg := r.Config()
	addrLabel := gutil.NewCopyButton(rCfgAddr)
	titleLabel := widget.NewLabel(rCfg.Title)
	noteLabel := widget.NewLabel("")

	entries := make([]*widget.Entry, len(rCfg.Labels))
	rForm := &widget.Form{}
	for idx, label := range rCfg.Labels{
		entries[idx] = widget.NewEntry()
		rForm.Items = append(rForm.Items, widget.NewFormItem(label, entries[idx]))
	}
	rForm.OnSubmit = func(){
		noteLabel.SetText("processing...")
		dataset := make([]string, len(rCfg.Labels))
		for idx, entry := range entries{
			dataset[idx] = entry.Text
		}
		uidStr, err := r.Registrate(dataset...)
		if err != nil{
			noteLabel.SetText("registration error: "+err.Error())
			return
		}
		noteLabel.SetText("done")
		uidLabel.SetText(uidStr)
	}
	rForm.ExtendBaseWidget(rForm)

	titles := container.NewVBox(addrLabel.Render(), titleLabel)
	page := container.NewVBox(rForm, noteLabel, uidLabel.Render())
	page = container.NewBorder(titles, nil, nil, nil, page)
	return rCfg.Title, page, closer
}
