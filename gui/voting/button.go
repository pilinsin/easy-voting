package votingpage

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	viface "github.com/pilinsin/easy-voting/voting/interface"
)

func voteBtn(v viface.IVoting, candNameGroups []string, label *widget.Label) fyne.CanvasObject {
	votingForm := v.NewVotingForm(candNameGroups)
	voteBtn := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		label.SetText("processing...")
		err := v.Vote(votingForm.VoteInt())
		if err != nil {
			label.SetText(fmt.Sprintln(err))
			return
		}
		label.SetText("voted")
	})
	return container.NewVBox(votingForm, voteBtn)
}
func checkMyVoteBtn(v viface.IVoting, note *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("check my vote", theme.DocumentIcon(), func() {
		note.SetText("processing...")
		vi, err := v.GetMyVote()
		if err != nil {
			note.SetText(fmt.Sprintln(err))
			return
		}
		note.SetText(fmt.Sprintln("my vote: ", *vi))
	})
}

func resultBtn(v viface.IVoting, note *widget.Label) fyne.CanvasObject {
	return widget.NewButtonWithIcon("result", theme.FolderOpenIcon(), func() {
		note.SetText("processing...")
		
		if vr, nVoters, nVoted, err := v.GetResult(); err != nil {
			note.SetText(fmt.Sprintln(err))
		} else {
			note.SetText(fmt.Sprintln("turnout:", nVoted, "/", nVoters, ", result: ", *vr))
		}
	})
}
