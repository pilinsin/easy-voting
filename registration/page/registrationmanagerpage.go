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

func LoadManPage(mCfgCid string, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	mCfg, err := rutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return rputil.ErrorPage(err), nil
	}
	m, err := rman.NewManager(mCfgCid, is)
	if err != nil {
		return rputil.ErrorPage(err), nil
	}

	rCfgCid := ipfs.ToCid(mCfg.Config().Marshal(), is)
	mCfgEntry := widget.NewEntry()
	mCfgEntry.Text = mCfgCid
	rCfgEntry := widget.NewEntry()
	rCfgEntry.Text = rCfgCid
	titleLabel := container.NewVBox(
		widget.NewLabel(mCfg.Title()),
		widget.NewLabel("manConfig:"),
		mCfgEntry,
		widget.NewLabel("rConfig:"),
		rCfgEntry,
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
