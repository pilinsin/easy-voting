package votingpage

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type VParamEntry struct {
	min   *gutil.IntEntry
	max   *gutil.IntEntry
	total *gutil.IntEntry
}

func NewVParamEntry() *VParamEntry {
	return &VParamEntry{
		min:   gutil.NewIntEntry(),
		max:   gutil.NewIntEntry(),
		total: gutil.NewIntEntry(),
	}
}
func (vpe *VParamEntry) Render() fyne.CanvasObject {
	form := widget.NewForm()
	form.Append("min", vpe.min)
	form.Append("max", vpe.max)
	form.Append("total", vpe.total)
	return form
}
func (vpe *VParamEntry) VoteParams() *vutil.VoteParams {
	return &vutil.VoteParams{
		Min:   vpe.min.Num(),
		Max:   vpe.max.Num(),
		Total: vpe.total.Num(),
	}
}
