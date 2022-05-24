package votingpage

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	vt "github.com/pilinsin/easy-voting/voting"
)

func LoadPage(ctx context.Context, vCfgAddr, baseDir string) (string, fyne.CanvasObject, func()) {
	v, err := vt.NewVoting(ctx, vCfgAddr, baseDir)
	if err != nil {
		return "", nil, nil
	}
	closer := func() { v.Close() }

	vCfg := v.Config()

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("User/Man Indentity")
	idEntry.OnChanged = func(s string){v.SetIdentity(s)}

	addrLabel := gutil.NewCopyButton(vCfgAddr)
	titleLabel := widget.NewLabel(vCfg.Title)
	noteLabel := widget.NewLabel("")
	resLabel := gutil.NewCopyButton("result:")

	vCfg.ShuffleCandidates()
	contents := CandCards(vCfg.Candidates)

	vBtn := voteBtn(v, vCfg.CandNameGroups(), noteLabel)
	cmvBtn := checkMyVoteBtn(v, resLabel)
	resBtn := resultBtn(v, resLabel)

	titles := container.NewVBox(titleLabel, addrLabel.Render())
	vObjs := container.NewVBox(idEntry, contents, vBtn, noteLabel)
	resObjs := container.NewVBox(container.NewHBox(cmvBtn, resBtn), resLabel.Render())
	page := container.NewVBox(vObjs, resObjs)
	page = container.NewBorder(titles, nil, nil, nil, page)
	return vCfg.Title, page, closer
}
