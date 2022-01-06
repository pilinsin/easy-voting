package votingpageutil

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/util"
	vutil "EasyVoting/voting/util"
)

func ErrorPage(err error) fyne.CanvasObject {
	return widget.NewLabel(fmt.Sprintln(err))
}

func PlaceHolderEntry(ph string) *widget.Entry {
	e := &widget.Entry{
		Wrapping:    fyne.TextTruncate,
		PlaceHolder: ph,
	}
	e.ExtendBaseWidget(e)
	return e
}

func intEntry() *widget.Entry {
	e := &widget.Entry{
		Wrapping: fyne.TextTruncate,
		Text:     "0",
	}
	e.OnChanged = func(val string) {
		if val == "" {
			return
		}
		if _, err := strconv.Atoi(val); err != nil {
			e.Text = "0"
			e.Refresh()
		}
	}
	e.ExtendBaseWidget(e)
	return e
}

func PlaceHolderSelect(opts []string, ph string, ch func(string)) *widget.Select {
	s := &widget.Select{
		OnChanged:   ch,
		Options:     opts,
		PlaceHolder: ph,
	}
	s.ExtendBaseWidget(s)
	return s
}

func setUrl(text string, urlStr string) fyne.CanvasObject {
	if urlStr != "" {
		parsedUrl, err := url.Parse(urlStr)
		if err == nil {
			return widget.NewHyperlink(text, parsedUrl)
		}
	}
	//When invalid urlStr, err is not raised.
	return widget.NewLabel("")
}

func resourceEqual(selected, def fyne.Resource) bool {
	name := selected.Name() == def.Name()
	content := util.ConstTimeBytesEqual(selected.Content(), def.Content())
	return name && content
}

func imageDialog(w fyne.Window, res fyne.Resource) func() {
	return func() {
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if rc == nil || err != nil {
				return
			}
			img := canvas.NewImageFromURI(rc.URI())
			img.FillMode = canvas.ImageFillContain
			res = img.Resource
		}
		dialog.ShowFileOpen(onSelected, w)
	}
}

type TimeSelect struct {
	Y *widget.Select
	M *widget.Select
	D *widget.Select
	h *widget.Select
	m *widget.Select
}

func lastDay(year, month string) int {
	t, _ := time.Parse(util.Layout, year+"-"+month+"-"+"1 0:00")
	return t.AddDate(0, 1, -1).Day()
}
func NewTimeSelect() *TimeSelect {
	now := time.Now()
	year := &widget.Select{
		Options:  util.ArangeStr(now.Year(), now.Year()+100, 1),
		Selected: strconv.Itoa(now.Year()),
	}
	month := &widget.Select{
		Options:  util.ArangeStr(1, 13, 1),
		Selected: strconv.Itoa(int(now.Month())),
	}
	day := &widget.Select{
		Options:  util.ArangeStr(1, 32, 1),
		Selected: strconv.Itoa(now.Day()),
	}
	year.OnChanged = func(y string) {
		d := lastDay(y, month.Selected)
		day.Options = util.ArangeStr(1, d+1, 1)
		selected, _ := strconv.Atoi(day.Selected)
		if selected > d {
			day.Selected = strconv.Itoa(d)
		}
		day.Refresh()
	}
	month.OnChanged = func(mth string) {
		d := lastDay(year.Selected, mth)
		day.Options = util.ArangeStr(1, d+1, 1)
		selected, _ := strconv.Atoi(day.Selected)
		if selected > d {
			day.Selected = strconv.Itoa(d)
		}
		day.Refresh()
	}
	year.ExtendBaseWidget(year)
	month.ExtendBaseWidget(month)
	day.ExtendBaseWidget(day)

	hour := &widget.Select{
		Options:  util.ArangeStr(0, 24, 1),
		Selected: "0",
	}
	hour.ExtendBaseWidget(hour)
	min := &widget.Select{
		Options:  util.ArangeStr(0, 60, 1),
		Selected: "0",
	}
	min.ExtendBaseWidget(min)

	return &TimeSelect{
		Y: year,
		M: month,
		D: day,
		h: hour,
		m: min,
	}
}
func (ts *TimeSelect) Render() fyne.CanvasObject {
	return container.NewHBox(ts.Y, ts.M, ts.D, ts.h, ts.m)
}
func (ts *TimeSelect) Time() string {
	Y := ts.Y.Selected
	M := ts.M.Selected
	D := ts.D.Selected
	h := ts.h.Selected
	m := ts.m.Selected
	//2006-1-2 15:4
	return Y + "-" + M + "-" + D + " " + h + ":" + m
}

type VParamEntry struct {
	min   *widget.Entry
	max   *widget.Entry
	total *widget.Entry
}

func NewVParamEntry() *VParamEntry {
	return &VParamEntry{
		min:   intEntry(),
		max:   intEntry(),
		total: intEntry(),
	}
}
func (vpe *VParamEntry) Render() fyne.CanvasObject {
	form := widget.NewForm()
	form.Append("min", vpe.min)
	form.Append("max", vpe.max)
	form.Append("total", vpe.total)
	return form
}
func (vpe *VParamEntry) VoteParams() vutil.VoteParams {
	min, _ := strconv.Atoi(vpe.min.Text)
	max, _ := strconv.Atoi(vpe.max.Text)
	total, _ := strconv.Atoi(vpe.total.Text)
	return vutil.VoteParams{min, max, total}
}
