package registrationpage

import (
	"context"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	rgst "github.com/pilinsin/easy-voting/registration"
	riface "github.com/pilinsin/easy-voting/registration/interface"
	rutil "github.com/pilinsin/easy-voting/registration/util"
)

func NewSetupPage(w fyne.Window, rs map[string]riface.IRegistration) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	addrLabel := gutil.NewCopyButton("registration config cid")
	if r, exist := rs["setup"]; exist {
		noteLabel.SetText("registration config is already generated")
		addrs := strings.Split(r.Address(), "/")
		addr := strings.Join(addrs[1:], "/")
		addrLabel.SetText(addr)
	}

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
		addrLabel.SetText("registration config cid")
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

		mapKey := "setup"
		if _, exist := rs[mapKey]; exist {
			rs[mapKey].Close()
			rs[mapKey] = nil
		}
		rCfgAddr := bAddrEntry.Text + "/" + cid
		r, err := rgst.NewRegistration(context.Background(), rCfgAddr, baseDir)
		if err != nil {
			noteLabel.SetText("new rConfig error: " + err.Error())
			return
		}
		noteLabel.SetText("done")
		addrLabel.SetText(cid)
		rs[mapKey] = r
		//form.Hide()
	}
	form.ExtendBaseWidget(form)

	page := container.NewVBox(form, noteLabel, addrLabel.Render())
	return page
}
