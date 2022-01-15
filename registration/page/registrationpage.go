package registrationpage

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	user "EasyVoting/registration"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
)

func LoadPage(rCfgCid string, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, nil
	}
	r, err := user.NewRegistration(rCfgCid, is)
	if err != nil {
		r.Close()
		return nil, nil
	}

	titleLabel := widget.NewLabel(rCfg.Title() + " (" + rCfgCid + ")")
	noteLabel := widget.NewLabel("")

	rForm := registrationForm(rCfg.UserDataLabels(), r, noteLabel)
	//idVerfBtn := identityVerifyButton(r, noteLabel)

	hnmVerifyLabel := widget.NewLabel("verifying HashNameMap...")
	ctx, cancel := util.CancelContext()
	newVerifyHnmGoRoutine(ctx, r, hnmVerifyLabel)
	closer := rputil.NewPageCloser(r.Close, cancel)

	page := container.NewVBox(rForm, noteLabel, hnmVerifyLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, closer
}

func registrationForm(labels []string, r user.IRegistration, noteLabel *widget.Label) fyne.CanvasObject {
	kwEntry := widget.NewEntry()
	kwEntry.SetPlaceHolder("iuput a userIdentity keyword")
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("user identity will be output here")

	rForm := &widget.Form{}
	entries := make([]*widget.Entry, len(labels))
	for idx, label := range labels {
		entries[idx] = widget.NewEntry()
		formItem := widget.NewFormItem(label, entries[idx])
		rForm.Items = append(rForm.Items, formItem)
	}
	rForm.Items = append(rForm.Items, widget.NewFormItem("keyword", kwEntry))
	rForm.OnSubmit = func() {
		noteLabel.SetText("processing...")
		var texts []string
		for _, entry := range entries {
			texts = append(texts, entry.Text)
		}
		if identity, err := r.Registrate(texts...); err != nil {
			noteLabel.SetText(fmt.Sprintln("registration error:", err))
		} else {
			mid := identity.Marshal()
			idStore := rutil.NewIdentityStore()
			idStore.Put(kwEntry.Text, mid)
			idStore.Close()
			idEntry.SetText(util.AnyBytes64ToStr(mid))
			noteLabel.SetText("registration is done")
		}
	}
	rForm.ExtendBaseWidget(rForm)
	return container.NewVBox(rForm, idEntry)
}
func newVerifyHnmGoRoutine(ctx context.Context, r user.IRegistration, label *widget.Label) {
	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("process stop")
				return
			case <-ticker.C:
				if ok := r.VerifyHashNameMap(); ok {
					label.SetText("HashNameMap is verified")
				} else {
					label.SetText("invalid HashNameMap")
					return
				}
			}
		}
	}(ctx)
}
