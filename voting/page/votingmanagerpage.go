package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/util"
	"EasyVoting/ipfs"
	rputil "EasyVoting/registration/page/util"
	"EasyVoting/voting"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
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
	getBtn := vputil.GetResultMapBtn(m, noteLabel)
	verifyBtn := vputil.VerifyResultMapBtn(m, noteLabel)

	page := container.NewVBox(contents, cuForm, getBtn, verifyBtn, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return page, rputil.NewPageCloser(m.Close, func() {})
}
