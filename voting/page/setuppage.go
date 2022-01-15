package votingpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	"EasyVoting/util"
	rutil "EasyVoting/registration/util"
	vputil "EasyVoting/voting/page/util"
	vutil "EasyVoting/voting/util"
)

func NewSetupPage(w fyne.Window, is *ipfs.IPFS) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	cidEntry := widget.NewEntry()
	cidEntry.SetPlaceHolder("voting config cid will be output here")

	title := widget.NewEntry()
	begin := vputil.NewTimeSelect()
	end := vputil.NewTimeSelect()
	loc := widget.NewSelect(util.GetOsTimeZones(), nil)
	rCfgCid := widget.NewEntry()
	cands := vputil.NewCandForm()
	vParam := vputil.NewVParamEntry()
	vType := widget.NewSelect(vutil.VotingTypes(), nil)

	kwEntry := widget.NewEntry()
	kwEntry.SetPlaceHolder("keyword of voting manager identity")

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", title))
	form.Items = append(form.Items, widget.NewFormItem("begin", begin.Render()))
	form.Items = append(form.Items, widget.NewFormItem("end", end.Render()))
	form.Items = append(form.Items, widget.NewFormItem("location", loc))
	form.Items = append(form.Items, widget.NewFormItem("rCfgCid", rCfgCid))
	form.Items = append(form.Items, widget.NewFormItem("candidates", cands.Render(w)))
	form.Items = append(form.Items, widget.NewFormItem("voteParams", vParam.Render()))
	form.Items = append(form.Items, widget.NewFormItem("voting type", vType))
	form.Items = append(form.Items, widget.NewFormItem("keyword", kwEntry))
	form.OnSubmit = func() {
		noteLabel.SetText("processing...")

		vt, err := vutil.VotingTypeFromStr(vType.Selected)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}
		if loc.Selected == ""{
			noteLabel.SetText("location is empty")
			return
		}
		candidates := cands.Candidates()
		if len(candidates) == 0{
			noteLabel.SetText("there are no candidates")
			return
		}

		manIdentity, vCfg, err := vutil.NewConfigs(
			title.Text,
			begin.Time(),
			end.Time(),
			loc.Selected,
			rCfgCid.Text,
			candidates,
			vParam.VoteParams(),
			vt,
			is,
		)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
		} else {
			vCfgCid := ipfs.ToCidWithAdd(vCfg.Marshal(), is)
			noteLabel.SetText("voting config cid:")
			cidEntry.SetText(vCfgCid)

			idStore := rutil.NewIdentityStore()
			idStore.Put(kwEntry.Text, manIdentity.Marshal())
			idStore.Close()
		}
	}
	form.ExtendBaseWidget(form)

	return container.NewVScroll(container.NewVBox(form, noteLabel, cidEntry))
}
