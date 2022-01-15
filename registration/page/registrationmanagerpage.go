package registrationpage

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rman "EasyVoting/registration"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
)

func LoadManPage(rCfgCid string, manIdentity *rutil.ManIdentity, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	rCfg, err := rutil.ConfigFromCid(rCfgCid, is)
	if err != nil {
		return nil, nil
	}
	if ok := rCfg.IsCompatible(manIdentity); !ok{
		return nil, nil
	}
	m, err := rman.NewManager(rCfgCid, manIdentity, is)
	if err != nil {
		m.Close()
		return nil, nil
	}
	
	rCfgEntry := widget.NewEntry()
	rCfgEntry.SetText(rCfgCid)
	mIdEntry := widget.NewEntry()
	mIdEntry.SetText(util.AnyBytes64ToStr(manIdentity.Marshal()))
	titleLabel := container.NewVBox(
		widget.NewLabel(rCfg.Title()),
		widget.NewLabel("rConfig:"),
		rCfgEntry,
		widget.NewLabel("registration manager identity:"),
		mIdEntry,
	)
	noteLabel := widget.NewLabel("")

	ctx, cancel := util.CancelContext()
	closer := rputil.NewPageCloser(m.Close, cancel)
	check := rCheck(ctx, closer, m, noteLabel)

	page := container.NewVBox(check, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, closer
}

func rCheck(ctx context.Context, closer rputil.IPageCloser, m rman.IManager, noteLabel *widget.Label) *widget.Check {
	check := &widget.Check{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "registration manager on/off:",
	}
	check.OnChanged = func(state bool) {
		if !state {
			noteLabel.SetText("")
			closer.Cancel()
		} else {
			ctx, cancel := util.CancelContext()
			closer.SetCancel(cancel)
			noteLabel.SetText("processing...")
			go func(ctx context.Context) {
				ticker := time.NewTicker(5 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						fmt.Println("process stop")
						return
					case <-ticker.C:
						if err := m.Registrate(); err != nil {
							noteLabel.SetText(fmt.Sprintln(err))
							check.Checked = false
							check.Refresh()
							return
						}
					}
				}
			}(ctx)
		}
	}
	check.ExtendBaseWidget(check)
	return check
}
