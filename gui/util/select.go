package guiutil

import(
	"time"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"

	"github.com/pilinsin/util"
)

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
