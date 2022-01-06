package votingpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
)

func NewSetupPage(w fyne.Window, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	cidEntry := widget.NewEntry()
	cidEntry.SetPlaceHolder("voting manager cid")

	title := widget.NewEntry()
	begin := vputil.NewTimeSelect()
	end := vputil.NewTimeSelect()
	loc := widget.NewSelect(util.GetOsTimeZones(), nil)
	rCfgCid := widget.NewEntry()
	cands := vputil.NewCandForm()
	vParam := vputil.NewVParamEntry()
	vType := widget.NewSelect(vutil.VotingTypes(), nil)

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", title))
	form.Items = append(form.Items, widget.NewFormItem("begin", begin.Render()))
	form.Items = append(form.Items, widget.NewFormItem("end", end.Render()))
	form.Items = append(form.Items, widget.NewFormItem("location", loc))
	form.Items = append(form.Items, widget.NewFormItem("rCfgCid", rCfgCid))
	form.Items = append(form.Items, widget.NewFormItem("candidates", cands.Render(w)))
	form.Items = append(form.Items, widget.NewFormItem("voteParams", vParam.Render()))
	form.Items = append(form.Items, widget.NewFormItem("voting type", vType))
	form.OnSubmit = func() {
		noteLabel.SetText("processing...")

		vt, err := vutil.VotingTypeFromStr(vType.Selected)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}

		mCfg, vCfg, err := vutil.NewConfigs(
			title.Text,
			begin.Time(),
			end.Time(),
			loc.Selected,
			rCfgCid.Text,
			cands.Candidates(),
			vParam.VoteParams(),
			vt,
			is,
		)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
		} else {
			mCfgCid := ipfs.ToCidWithAdd(mCfg.Marshal(), is)
			ipfs.ToCidWithAdd(vCfg.Marshal(), is)
			noteLabel.SetText("voting manager cid:")
			cidEntry.SetText(mCfgCid)
		}
	}
	form.ExtendBaseWidget(form)

	return container.NewVScroll(container.NewVBox(form, noteLabel, cidEntry))
}
