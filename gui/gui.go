package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/pilinsin/easy-voting/util"
	rpage "github.com/pilinsin/easy-voting/registration/page"
	rputil "github.com/pilinsin/easy-voting/registration/page/util"
	rutil "github.com/pilinsin/easy-voting/registration/util"
	vpage "github.com/pilinsin/easy-voting/voting/page"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

type GUI struct {
	w    fyne.Window
	size fyne.Size
	page *fyne.Container
}

func New(title string, width, height float32) *GUI {
	size := fyne.NewSize(width, height)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	win := a.NewWindow(title)
	win.Resize(size)
	page := container.NewMax()
	return &GUI{win, size, page}
}

func (gui *GUI) withRemove(page fyne.CanvasObject, closer rputil.IPageCloser) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		closer.Close()
		gui.changePage(gui.defaultPage())
	})
	return container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, page)
}

func (gui *GUI) changePage(page fyne.CanvasObject) {
	for _, obj := range gui.page.Objects {
		gui.page.Remove(obj)
	}
	gui.page.Add(page)
	gui.page.Refresh()
	gui.w.Resize(gui.size)
}

func (gui *GUI) loadPage(addr string) (fyne.CanvasObject, rputil.IPageCloser) {
	if ok := strings.HasPrefix(addr, "r/"); ok{
		return rpage.LoadPage(strings.TrimPrefix(addr, "r/"))
	}
	if ok := strings.HasPrefix(addr, "v/"); ok{
		return vpage.LoadPage(strings.TrimPrefix(addr, "v/"))
	}
	return nil, nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	cidEntry := widget.NewEntry()
	cidEntry.PlaceHolder = "registration/voting Config Address"
	loadBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		loadPage, closer := gui.loadPage(cidEntry.Text)
		if loadPage == nil {
			cidEntry.SetText("")
			return
		}
		page := container.NewVScroll(loadPage)
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		gui.changePage(gui.withRemove(page, closer))
	})
	return container.NewBorder(nil, nil, nil, loadBtn, cidEntry)
}

func (gui *GUI) newPageForm() fyne.CanvasObject {
	var setup fyne.CanvasObject
	chmod := &widget.Select{
		Options:  []string{"registration", "voting"},
		Selected: "registration",
	}
	chmod.OnChanged = func(mode string) {
		if mode == "registration" {
			setup = rpage.NewSetupPage(gui.w, gui.is)
		} else {
			setup = vpage.NewSetupPage(gui.w, gui.is)
		}
		newForm := container.NewBorder(chmod, nil, nil, nil, setup)
		defPage := container.NewBorder(gui.loadPageForm(), nil, nil, nil, newForm)
		gui.changePage(defPage)
	}
	chmod.ExtendBaseWidget(chmod)

	return container.NewBorder(chmod, nil, nil, nil)
}

func (gui *GUI) defaultPage() fyne.CanvasObject {
	loadForm := gui.loadPageForm()
	newForm := gui.newPageForm()
	return container.NewBorder(loadForm, nil, nil, nil, newForm)
}

func (gui *GUI) Run() {
	gui.page.Add(gui.defaultPage())
	gui.w.SetContent(gui.page)
	gui.w.ShowAndRun()
}


