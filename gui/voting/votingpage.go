package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
)

func LoadPage(bAddr, vCfgCid string, v viface.IVoting) (string, fyne.CanvasObject) {
	vCfg := v.Config()

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("User/Man Indentity")
	idEntry.OnChanged = func(s string) { v.SetIdentity(s) }

	bAddrLabel := gutil.NewCopyButton(bAddr)
	cfgLabel := gutil.NewCopyButton(vCfgCid)
	addrLabel := container.NewVBox(bAddrLabel.Render(), cfgLabel.Render())
	titleLabel := widget.NewLabel(vCfg.Title)
	noteLabel := widget.NewLabel("")
	resLabel := gutil.NewCopyButton("result:")

	vCfg.ShuffleCandidates()
	contents := CandCards(vCfg.Candidates)

	vBtn := voteBtn(v, vCfg.CandNameGroups(), noteLabel)
	cmvBtn := checkMyVoteBtn(v, resLabel)
	resBtn := resultBtn(v, resLabel)

	titles := container.NewVBox(titleLabel, addrLabel)
	vObjs := container.NewVBox(idEntry, contents, vBtn, noteLabel)
	resObjs := container.NewVBox(container.NewHBox(cmvBtn, resBtn), resLabel.Render())
	page := container.NewVBox(vObjs, resObjs)
	page = container.NewBorder(titles, nil, nil, nil, page)
	return vCfg.Title, page
}
