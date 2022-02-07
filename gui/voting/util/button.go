package votingpageutil

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	viface "github.com/pilinsin/easy-voting/voting/interface"
)

func VotingBtn(v viface.IVoting, candNameGroups []string, label *widget.Label) fyne.CanvasObject {
	votingForm := v.NewVotingForm(candNameGroups)
	voteBtn := widget.NewButtonWithIcon("", theme.UploadIcon(), func() {
		label.SetText("processing...")
		err := v.Vote(votingForm.VoteInt())
		if err != nil {
			label.SetText(fmt.Sprintln(err))
			return
		}
		label.SetText("voting has been done.")
	})
	return container.NewVBox(votingForm, voteBtn)
}
func CheckMyVoteBtn(v viface.IVoting, e *widget.Entry, label *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("check my vote", theme.CheckButtonIcon(), func() {
		label.SetText("processing...")
		cid, err := v.GetMyVote()
		if err != nil {
			label.SetText(fmt.Sprintln(err))
			return
		}
		e.SetText(fmt.Sprintln("my vote cid: ", cid))
	})
}

func CountBtn(v viface.IVoting, e *widget.Entry, label *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("count", theme.DocumentIcon(), func() {
		label.SetText("processing...")
		if ok, err := v.VerifyResultMap(); err != nil {
			label.SetText(fmt.Sprintln(err))
		} else if !ok {
			label.SetText("invalid resultMap")
		} else {
			if cid, err := v.Count(); err != nil {
				e.SetText(fmt.Sprintln(err))
			} else {
				e.SetText(fmt.Sprintln("result cid: ", cid))
			}
		}
	})
}

func CheckUserForm(labels []string, m viface.IManager, noteLabel *widget.Label) fyne.CanvasObject {
	cuForm := &widget.Form{}
	entries := make([]*widget.Entry, len(labels))
	for idx, label := range labels {
		entries[idx] = widget.NewEntry()
		formItem := widget.NewFormItem(label, entries[idx])
		cuForm.Items = append(cuForm.Items, formItem)
	}
	cuForm.OnSubmit = func() {
		noteLabel.SetText("processing...")
		var texts []string
		for _, entry := range entries {
			texts = append(texts, entry.Text)
		}
		if ok := m.IsValidUser(texts...); ok {
			noteLabel.SetText("valid")
		} else {
			noteLabel.SetText("NOT valid")
		}
	}
	cuForm.ExtendBaseWidget(cuForm)
	return cuForm
}

func GetResultMapBtn(m viface.IManager, label *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("get ResultMap", theme.DownloadIcon(), func() {
		label.SetText("processing...")
		if err := m.GetResultMap(); err != nil {
			label.SetText(fmt.Sprintln(err))
		} else {
			label.SetText("GetResultMap has finished")
		}
	})
}

func VerifyResultMapBtn(m viface.IManager, label *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("verify ResultMap", theme.CheckButtonCheckedIcon(), func() {
		label.SetText("processing...")
		if ok, err := m.VerifyResultMap(); err != nil {
			label.SetText(fmt.Sprintln(err))
		} else if !ok {
			label.SetText("invalid resultMap")
		} else {
			label.SetText("valid resultMap")
		}
	})
}
