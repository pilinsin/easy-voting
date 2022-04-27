package votingpage

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	vt "github.com/pilinsin/easy-voting/voting"
	gutil "github.com/pilinsin/easy-voting/gui/util"
)

func LoadPage(ctx context.Context, vCfgAddr, idStr string) (fyne.CanvasObject, func()) {
	v, err := vt.NewVoting(ctx, vCfgAddr, idStr)
	if err != nil{return nil, nil}
	closer := func(){v.Close()}

	vCfg := v.Config()

	addrLabel := gutil.NewCopyButton(vCfgAddr)
	titleLabel := widget.NewLabel(vCfg.Title)
	noteLabel := widget.NewLabel("")

	vCfg.ShuffleCandidates()
	contents := CandCards(vCfg.Candidates)

	vbtn := voteBtn(v, vCfg.CandNameGroups(), noteLabel)
	cmvbtn := checkMyVoteBtn(v, noteLabel)
	rbtn := resultBtn(v, noteLabel)

	titles := container.NewVBox(addrLabel.Render(), titleLabel)
	page := container.NewVBox(contents, vbtn, cmvbtn, rbtn, noteLabel)
	page = container.NewBorder(titles, nil, nil, nil, page)
	return page, closer
}

