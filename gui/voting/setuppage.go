package votingpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/pilinsin/util"
	gutil "github.com/pilinsin/easy-voting/gui/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

func NewSetupPage(w fyne.Window) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	addrLabel := gutil.NewCopyButton("voting config address")
	maIdLabel := gutil.NewCopyButton("voting manager address")

	title := widget.NewEntry()
	begin := gutil.NewTimeSelect()
	end := gutil.NewTimeSelect()
	loc := widget.NewSelect(util.GetOsTimeZones(), nil)
	rCfgAddr := widget.NewEntry()
	nVerifiers := gutil.NewIntEntry()
	cands := NewCandForm()
	vParam := NewVParamEntry()
	vType := widget.NewSelect(vutil.VotingTypes(), nil)

	nVerifiers.SetPlaceHolder("1~")

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", title))
	form.Items = append(form.Items, widget.NewFormItem("begin", begin.Render()))
	form.Items = append(form.Items, widget.NewFormItem("end", end.Render()))
	form.Items = append(form.Items, widget.NewFormItem("location", loc))
	form.Items = append(form.Items, widget.NewFormItem("rCfgAddr", rCfgAddr))
	form.Items = append(form.Items, widget.NewFormItem("nVerifiers", nVerifiers))
	form.Items = append(form.Items, widget.NewFormItem("candidates", cands.Render(w)))
	form.Items = append(form.Items, widget.NewFormItem("voteParams", vParam.Render()))
	form.Items = append(form.Items, widget.NewFormItem("voting type", vType))
	form.OnSubmit = func() {
		noteLabel.SetText("processing...")

		tInfo, err := util.NewTimeInfo(begin.Time(), end.Time(), loc.Selected)
		if err != nil{
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}

		vt, err := vutil.VotingTypeFromStr(vType.Selected)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}
		if loc.Selected == ""{
			noteLabel.SetText("location is empty")
			return
		}
		if nVerifiers.Num() <= 0{
			noteLabel.SetText("nVerifiers must be positive")
			return
		}
		candidates := cands.Candidates()
		if len(candidates) == 0{
			noteLabel.SetText("there are no candidates")
			return
		}

		cid, mid, err := vutil.NewConfig(
			title.Text,
			rCfgAddr.Text,
			nVerifiers.Num(),
			tInfo,
			candidates,
			vParam.VoteParams(),
			vt,
		)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
		} else {
			noteLabel.SetText("done")
			addrLabel.SetText(cid)
			maIdLabel.SetText(mid)
		}
	}
	form.ExtendBaseWidget(form)

	return container.NewVScroll(container.NewVBox(form, noteLabel, addrLabel.Render(), maIdLabel.Render()))
}
