package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	"EasyVoting/voting"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
)

func LoadManPage(mCfgCid string, is *ipfs.IPFS) (fyne.CanvasObject, rputil.IPageCloser) {
	mCfg, err := vutil.ManConfigFromCid(mCfgCid, is)
	if err != nil {
		return vputil.ErrorPage(err), nil
	}
	m, err := voting.NewManager(mCfgCid, is)
	if err != nil {
		return vputil.ErrorPage(err), nil
	}

	vCfgCid := ipfs.ToCid(mCfg.Config().Marshal(), is)
	mCfgEntry := widget.NewEntry()
	mCfgEntry.Text = mCfgCid
	vCfgEntry := widget.NewEntry()
	vCfgEntry.Text = vCfgCid
	titleLabel := container.NewVBox(
		widget.NewLabel(mCfg.Title()),
		widget.NewLabel("manConfig:"),
		mCfgEntry,
		widget.NewLabel("vConfig:"),
		vCfgEntry,
	)
	noteLabel := widget.NewLabel("")

	mCfg.ShuffleCandidates()
	contents := vputil.CandCards(mCfg.Candidates())

	cuForm := vputil.CheckUserForm(mCfg.UserDataLabels(), m, noteLabel)
	getBtn := vputil.GetResultMapBtn(m, noteLabel)
	verifyBtn := vputil.VerifyResultMapBtn(m, noteLabel)

	page := container.NewVBox(contents, cuForm, getBtn, verifyBtn, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, rputil.NewPageCloser(m.Close, func() {})
}
