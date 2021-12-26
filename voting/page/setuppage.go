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

type setupPage struct {
	fyne.CanvasObject
	is *ipfs.IPFS
}

func NewSetupPage(a fyne.App, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")

	title := vputil.PlaceHolderEntry("title")
	begin := vputil.NewTimeSelect()
	end := vputil.NewTimeSelect()
	loc := vputil.PlaceHolderSelect(util.GetOsTimeZones(), "location", nil)
	rCfgCid := vputil.PlaceHolderEntry("rCfgCid")
	cands := vputil.NewCandForm(a)
	vParam := vputil.NewVParamEntry()
	vType := vputil.PlaceHolderSelect(vutil.VotingTypes(), "voting type", nil)

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", title))
	form.Items = append(form.Items, widget.NewFormItem("begin", begin))
	form.Items = append(form.Items, widget.NewFormItem("end", end))
	form.Items = append(form.Items, widget.NewFormItem("location", loc))
	form.Items = append(form.Items, widget.NewFormItem("rCfgCid", rCfgCid))
	form.Items = append(form.Items, widget.NewFormItem("candidates", cands))
	form.Items = append(form.Items, widget.NewFormItem("voteParams", vParam))
	form.Items = append(form.Items, widget.NewFormItem("voting type", vType))
	form.OnSubmit = func() {
		noteLabel.Text = "processing..."

		vt, err := vutil.VotingTypeFromStr(vType.Selected)
		if err != nil {
			noteLabel.Text = fmt.Sprintln(err)
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
			noteLabel.Text = fmt.Sprintln(err)
		} else {
			mCfgCid := ipfs.ToCidWithAdd(mCfg.Marshal(), is)
			ipfs.ToCidWithAdd(vCfg.Marshal(), is)
			noteLabel.Text = "voting manager cid: " + mCfgCid
		}
	}
	form.ExtendBaseWidget(form)
	return &setupPage{container.NewVBox(form, noteLabel), is}
}
