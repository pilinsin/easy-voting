package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
	"EasyVoting/voting"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
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

	counter := vputil.CountBtn(v, noteLabel)

	page := container.NewVBox(contents, idPage, counter, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, rputil.NewPageCloser(v.Close, func() {})
}

