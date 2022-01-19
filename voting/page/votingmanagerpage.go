package votingpage

import (
	"fmt"
	"context"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/util"
	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	"EasyVoting/voting"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
	viface "EasyVoting/voting/interface"
)

func LoadManPage(vCfgCid string, manIdentity *vutil.ManIdentity, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return nil, nil
	}
	if ok := vCfg.IsCompatible(manIdentity); !ok{
		return nil, nil
	}
	m, err := voting.NewManager(vCfgCid, manIdentity, is)
	if err != nil {
		m.Close()
		return nil, nil
	}

	vCfgEntry := widget.NewEntry()
	vCfgEntry.Text = vCfgCid
	midEntry := widget.NewEntry()
	midEntry.SetText(util.AnyBytes64ToStr(manIdentity.Marshal()))
	titleLabel := container.NewVBox(
		widget.NewLabel(vCfg.Title()),
		widget.NewLabel("vConfig:"),
		vCfgEntry,
		widget.NewLabel("voting manager identity:"),
		midEntry,
	)
	noteLabel := widget.NewLabel("")

	vCfg.ShuffleCandidates()
	contents := vputil.CandCards(vCfg.Candidates())

	cuForm := vputil.CheckUserForm(vCfg.UserDataLabels(), m, noteLabel)

	ctx, cancel := util.CancelContext()
	closer := rputil.NewPageCloser(m.Close, cancel)
	vkCheck := verfKeyCheck(ctx, closer, m, noteLabel)

	getBtn := vputil.GetResultMapBtn(m, noteLabel)
	verifyBtn := vputil.VerifyResultMapBtn(m, noteLabel)

	page := container.NewVBox(contents, cuForm, vkCheck, getBtn, verifyBtn, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, closer
}

func verfKeyCheck(ctx context.Context, closer rputil.IPageCloser, m viface.IManager, noteLabel *widget.Label) *widget.Check {
	check := &widget.Check{
		DisableableWidget: widget.DisableableWidget{},
		Text:              "verfKey check on/off:",
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

