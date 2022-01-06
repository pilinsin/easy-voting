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
		return rputil.ErrorPage(err), nil
	}
	r, err := user.NewRegistration(rCfgCid, is)
	if err != nil {
		return rputil.ErrorPage(err), nil
	}

	titleLabel := widget.NewLabel(rCfg.Title() + " (" + rCfgCid + ")")
	noteLabel := widget.NewLabel("")

	rForm := registrationForm(rCfg.UserDataLabels(), r, noteLabel)
	idVerfBtn := identityVerifyButton(r, noteLabel)

	hnmVerifyLabel := widget.NewLabel("verifying HashNameMap...")
	ctx, cancel := util.CancelContext()
	newVerifyHnmGoRoutine(ctx, r, hnmVerifyLabel)
	closer := rputil.NewPageCloser(r.Close, cancel)

	page := container.NewVBox(rForm, idVerfBtn, noteLabel, hnmVerifyLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, closer
}

func registrationForm(labels []string, r user.IRegistration, noteLabel *widget.Label) fyne.CanvasObject {
	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("userIdentity will be output here")

	rForm := &widget.Form{}
	entries := make([]*widget.Entry, len(labels))
	for idx, label := range labels {
		entries[idx] = widget.NewEntry()
		formItem := widget.NewFormItem(label, entries[idx])
		rForm.Items = append(rForm.Items, formItem)
	}
	rForm.OnSubmit = func() {
		noteLabel.SetText("processing...")
		var texts []string
		for _, entry := range entries {
			texts = append(texts, entry.Text)
		}
		if identity, err := r.Registrate(texts...); err != nil {
			noteLabel.SetText(fmt.Sprintln("registration error:", err))
		} else {
			idStr := util.AnyBytes64ToStr(identity.Marshal())
			noteLabel.SetText("registration is done")
			idEntry.SetText(idStr)
		}
	}
	rForm.ExtendBaseWidget(rForm)
	return container.NewVBox(rForm, idEntry)
}

func identityVerifyButton(r user.IRegistration, noteLabel *widget.Label) fyne.CanvasObject {
	e := widget.NewEntry()
	e.SetPlaceHolder("input userIdentity")
	f := &widget.Form{}
	f.Items = append(f.Items, widget.NewFormItem("verify registration", e))
	f.OnSubmit = func() {
		noteLabel.SetText("processing...")
		str := e.Text
		m := util.StrToAnyBytes64(str)
		identity := &rutil.UserIdentity{}
		if err := identity.Unmarshal(m); err != nil {
			noteLabel.SetText("this identity is invalid")
			return
		}
		if ok := r.VerifyUserIdentity(identity); ok {
			noteLabel.SetText("this identity is registrated")
		} else {
			noteLabel.SetText("this identity is not registrated")
		}
	}
	f.ExtendBaseWidget(f)
	return f
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
