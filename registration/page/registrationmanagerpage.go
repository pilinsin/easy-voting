package registrationpage

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rman "EasyVoting/registration"
	rutil "EasyVoting/registration/util"
)

type managerPage struct {
	fyne.CanvasObject
	m rman.IManager
}

func LoadManPage(mCfgCid string, is *ipfs.IPFS) fyne.CanvasObject {
	mCfg, err := rutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return errorPage(err)
	}
	r, err := rman.NewManager(mCfgCid, is)
	if err != nil {
		return errorPage(err)
	}

	rCfgCid := ipfs.ToCid(mCfg.Config().Marshal(), is)
	titleLabel := container.NewVBox(
		widget.NewLabel(mCfg.Title()),
		widget.NewLabel("manConfig:("+mCfgCid),
		widget.NewLabel("rConfig:("+rCfgCid),
	)
	noteLabel := widget.NewLabel("")

	check := rCheck(r, noteLabel)

	page := container.NewVBox(check, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return &managerPage{page, r}
}

func rCheck(m rman.IManager, noteLabel *widget.Label) *widget.Check {
	check := &widget.Check{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "registration manager on/off:",
	}
	check.OnChanged = func(state bool) {
		if !state {
			noteLabel.Text = ""
		} else {
			noteLabel.Text = "processing..."
			for {
				err := m.Registrate()
				if err != nil {
					noteLabel.Text = fmt.Sprintln(err)
					check.Checked = false
					check.Refresh()
					return
				} else {
					<-time.After(5 * time.Second)
				}
			}
		}
	}
	check.ExtendBaseWidget(check)
	return check
}
