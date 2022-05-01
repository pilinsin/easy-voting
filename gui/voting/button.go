package votingpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	viface "github.com/pilinsin/easy-voting/voting/interface"
)

func voteBtn(v viface.IVoting, candNameGroups []string, label *widget.Label) fyne.CanvasObject {
	votingForm := v.NewVotingForm(candNameGroups)
	voteBtn := widget.NewButtonWithIcon("vote", theme.MailSendIcon(), func() {
		label.SetText("processing...")
		vi := votingForm.VoteInt()
		err := v.Vote(vi)
		if err != nil {
			label.SetText(fmt.Sprintln(err))
			return
		}
		label.SetText(fmt.Sprintln("voted:", vi))
	})
	return container.NewVBox(votingForm, voteBtn)
}
func checkMyVoteBtn(v viface.IVoting, copy *gutil.CopyButton) fyne.CanvasObject {
	return widget.NewButtonWithIcon("check my vote", theme.DocumentIcon(), func() {
		copy.SetText("processing...")
		vi, err := v.GetMyVote()
		if err != nil {
			copy.SetText(fmt.Sprintln(err))
			return
		}
		copy.SetText(fmt.Sprintln("my vote: ", *vi))
	})
}

func resultBtn(v viface.IVoting, copy *gutil.CopyButton) fyne.CanvasObject {
	return widget.NewButtonWithIcon("result", theme.FolderOpenIcon(), func() {
		copy.SetText("processing...")

		if vr, nVoters, nVoted, err := v.GetResult(); err != nil {
			copy.SetText(fmt.Sprintln(err))
		} else {
			copy.SetText(fmt.Sprintln("turnout:", nVoted, "/", nVoters, ", result: ", *vr))
		}
	})
}
