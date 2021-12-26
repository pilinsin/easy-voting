package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	"EasyVoting/voting"
	viface "EasyVoting/voting/interface"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
)

type managerPage struct {
	fyne.CanvasObject
	m viface.IManager
}

func LoadManPage(mCfgCid string, is *ipfs.IPFS) fyne.CanvasObject {
	mCfg, err := vutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return vputil.ErrorPage(err)
	}
	m, err := voting.NewManager(mCfgCid, is)
	if err != nil {
		return vputil.ErrorPage(err)
	}

	vCfgCid := ipfs.ToCid(mCfg.Config().Marshal(), is)
	titleLabel := container.NewVBox(
		widget.NewLabel(mCfg.Title()),
		widget.NewLabel("manConfig:("+mCfgCid),
		widget.NewLabel("vConfig:("+vCfgCid),
	)
	noteLabel := widget.NewLabel("")

	mCfg.ShuffleCandidates()
	contents := vputil.CandCards(mCfg.Candidates())

	cuForm := vputil.CheckUserForm(mCfg.UserDataLabels(), m, noteLabel)
	getBtn := vputil.GetResultMapBtn(m, noteLabel)
	verifyBtn := vputil.VerifyResultMapBtn(m, noteLabel)

	page := container.NewVBox(contents, cuForm, getBtn, verifyBtn, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return &managerPage{page, m}
}
