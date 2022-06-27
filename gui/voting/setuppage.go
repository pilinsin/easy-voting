package votingpage

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	voting "github.com/pilinsin/easy-voting/voting"
	viface "github.com/pilinsin/easy-voting/voting/interface"
	vutil "github.com/pilinsin/easy-voting/voting/util"
	"github.com/pilinsin/util"
)

func NewSetupPage(w fyne.Window, vs map[string]viface.IVoting) fyne.CanvasObject {
	noteLabel := widget.NewLabel("")
	addrLabel := gutil.NewCopyButton("voting config cid")
	maIdLabel := gutil.NewCopyButton("voting manager address")
	if v, exist := vs["setup"]; exist {
		noteLabel.SetText("voting config is already generated")
		addrs := strings.Split(v.Address(), "/")
		addr := strings.Join(addrs[1:], "/")
		addrLabel.SetText(addr)
		maIdLabel.SetText(v.GetIdentity())
	}

	title := widget.NewEntry()
	begin := gutil.NewTimeSelect()
	end := gutil.NewTimeSelect()
	loc := widget.NewSelect(util.GetOsTimeZones(), nil)
	bAddr := widget.NewEntry()
	bAddr.SetPlaceHolder("Bootstraps Address")
	rCfgCid := widget.NewEntry()
	rCfgCid.SetPlaceHolder("Registration Config Cid")
	rCfgAddr := container.NewGridWithColumns(2, bAddr, rCfgCid)
	cands := NewCandForm()
	vParam := NewVParamEntry()
	vType := widget.NewSelect(vutil.VotingTypes(), nil)

	form := &widget.Form{}
	form.Items = append(form.Items, widget.NewFormItem("title", title))
	form.Items = append(form.Items, widget.NewFormItem("begin", begin.Render()))
	form.Items = append(form.Items, widget.NewFormItem("end", end.Render()))
	form.Items = append(form.Items, widget.NewFormItem("location", loc))
	form.Items = append(form.Items, widget.NewFormItem("rCfgAddr", rCfgAddr))
	form.Items = append(form.Items, widget.NewFormItem("candidates", cands.Render(w)))
	form.Items = append(form.Items, widget.NewFormItem("voteParams", vParam.Render()))
	form.Items = append(form.Items, widget.NewFormItem("voting type", vType))
	form.OnSubmit = func() {
		noteLabel.SetText("processing...")
		addrLabel.SetText("voting config cid")
		maIdLabel.SetText("voting manager address")

		tInfo, err := util.NewTimeInfo(begin.Time(), end.Time(), loc.Selected)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}

		vt, err := vutil.VotingTypeFromStr(vType.Selected)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}
		if loc.Selected == "" {
			noteLabel.SetText("location is empty")
			return
		}
		candidates := cands.Candidates()
		if len(candidates) == 0 {
			noteLabel.SetText("there are no candidates")
			return
		}

		rCfgAddr := bAddr.Text + "/" + rCfgCid.Text
		cid, mid, vStores, err := vutil.NewConfig(
			title.Text,
			rCfgAddr,
			tInfo,
			candidates,
			vParam.VoteParams(),
			vt,
		)
		if err != nil {
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}

		mapKey := "setup"
		if _, exist := vs[mapKey]; exist {
			vs[mapKey].Close()
			vs[mapKey] = nil
		}
		vCfgAddr := bAddr.Text + "/" + cid
		v, err := voting.NewVotingWithStores(vCfgAddr, vStores.Is, vStores.Hkm, vStores.Ivm)
		if err != nil {
			vStores.Close()
			noteLabel.SetText(fmt.Sprintln(err))
			return
		}
		v.SetIdentity(mid)

		noteLabel.SetText("done")
		addrLabel.SetText(cid)
		maIdLabel.SetText(mid)
		vs[mapKey] = v
	}
	form.ExtendBaseWidget(form)

	page := container.NewVScroll(container.NewVBox(form, noteLabel, addrLabel.Render(), maIdLabel.Render()))
	return page
}
