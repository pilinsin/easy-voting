package votingpage

import (
	"fmt"
	"time"
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/util"
	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
	"EasyVoting/voting"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
	viface "EasyVoting/voting/interface"
)

func LoadPage(vCfgCid string, userIdentity *rutil.UserIdentity, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return nil, nil
	}
	v, err := voting.NewVoting(vCfgCid, userIdentity, is)
	if err != nil {
		v.Close()
		return nil, nil
	}

	titleLabel := widget.NewLabel(vCfg.Title() + " (" + vCfgCid + ")")
	noteLabel := widget.NewLabel("")

	vCfg.ShuffleCandidates()
	contents := vputil.CandCards(vCfg.Candidates())

	var idPage fyne.CanvasObject
	if ok := v.VerifyIdentity(); !ok {
		idPage = nil
	} else {
		voteBtn := vputil.VotingBtn(v, vCfg.CandNameGroups(), noteLabel)
		checkBtn := vputil.CheckMyVoteBtn(v, noteLabel)
		idPage = container.NewVBox(voteBtn, checkBtn)
	}

	verfMapVerifyLabel := widget.NewLabel("verifying IdVerfKeyMap...")
	ctx, cancel := util.CancelContext()
	newVerifyVerfMapGoRoutine(ctx, v, verfMapVerifyLabel)
	closer := rputil.NewPageCloser(v.Close, cancel)

	counter := vputil.CountBtn(v, noteLabel)

	page := container.NewVBox(contents, idPage, counter, noteLabel, verfMapVerifyLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, closer
}

func newVerifyVerfMapGoRoutine(ctx context.Context, v viface.IVoting, label *widget.Label) {
	go func(ctx context.Context) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				fmt.Println("process stop")
				return
			case <-ticker.C:
				if ok := v.VerifyIdVerfKeyMap(); ok {
					label.SetText("IdVerfKeyMap is verified")
				} else {
					label.SetText("invalid IdVerfKeyMap")
					return
				}
			}
		}
	}(ctx)
}

