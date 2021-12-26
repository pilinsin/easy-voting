package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rutil "EasyVoting/registration/util"
	"EasyVoting/util"
	"EasyVoting/voting"
	viface "EasyVoting/voting/interface"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
)

type votingPage struct {
	fyne.CanvasObject
	v viface.IVoting
}

func LoadPage(vCfgCid string, is *ipfs.IPFS) fyne.CanvasObject {
	vCfg, err := vutil.ConfigFromCid(vCfgCid, is)
	if err != nil {
		return vputil.ErrorPage(err)
	}
	v, err := voting.NewVoting(vCfgCid, &rutil.UserIdentity{}, is)
	if err != nil {
		return vputil.ErrorPage(err)
	}

	titleLabel := widget.NewLabel(vCfg.Title() + " (" + vCfgCid + ")")
	noteLabel := widget.NewLabel("")

	vCfg.ShuffleCandidates()
	contents := vputil.CandCards(vCfg.Candidates())

	var idPage fyne.CanvasObject
	idEntry := identityEntry(vCfgCid, vCfg.CandNameGroups(), is, v, idPage, noteLabel)

	counter := vputil.CountBtn(v, noteLabel)

	page := container.NewVBox(contents, idEntry, idPage, counter, noteLabel)
	page = container.NewBorder(titleLabel, nil, nil, nil, page)
	return &votingPage{page, v}
}

func identityEntry(vCfgCid string, candNameGroups []string, is *ipfs.IPFS, v viface.IVoting, idPage fyne.CanvasObject, label *widget.Label) fyne.CanvasObject {
	e := &widget.Entry{Wrapping: fyne.TextTruncate}
	e.OnSubmitted = func(text string) {
		m := util.StrToBytes64(text)
		identity := &rutil.UserIdentity{}
		if err := identity.Unmarshal(m); err != nil {
			v, _ = voting.NewVoting(vCfgCid, &rutil.UserIdentity{}, is)
			idPage = nil
		} else {
			v, _ = voting.NewVoting(vCfgCid, identity, is)
			if ok := v.VerifyIdentity(); !ok {
				idPage = nil
			} else {
				voteBtn := vputil.VotingBtn(v, candNameGroups, label)
				checkBtn := vputil.CheckMyVoteBtn(v, label)
				idPage = container.NewVBox(voteBtn, checkBtn)
			}
		}
	}
	e.ExtendBaseWidget(e)
	return e
}
