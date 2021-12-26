package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rpage "EasyVoting/registration/page"
	rutil "EasyVoting/registration/util"
	vpage "EasyVoting/voting/page"
	vutil "EasyVoting/voting/util"
)

type GUI struct {
	a     fyne.App
	w     fyne.Window
	page  fyne.CanvasObject
	setup fyne.CanvasObject
	is    *ipfs.IPFS
}

func New(title string, is *ipfs.IPFS) *GUI {
	a := app.New()
	win := a.NewWindow(title)
	return &GUI{a: a, w: win, is: is}
}

func (gui *GUI) withRemove(page fyne.CanvasObject) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		gui.changePage(gui.defaultPage())
	})
	return container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, page)
}

func (gui *GUI) changePage(page fyne.CanvasObject) {
	gui.page = page
}

func (gui *GUI) loadPage(cid string) fyne.CanvasObject {
	if _, err := rutil.ConfigFromCid(cid, gui.is); err == nil {
		return rpage.LoadPage(cid, gui.is)
	}
	if _, err := rutil.ManConfigFromCid(cid, gui.is); err == nil {
		return rpage.LoadManPage(cid, gui.is)
	}
	if _, err := vutil.ConfigFromCid(cid, gui.is); err == nil {
		return vpage.LoadPage(cid, gui.is)
	}
	if _, err := vutil.ManConfigFromCid(cid, gui.is); err == nil {
		return vpage.LoadManPage(cid, gui.is)
	}
	return nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	cidEntry := widget.NewEntry()
	cidEntry.PlaceHolder = "page CID (Qm...)"
	loadBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		page := container.NewVScroll(gui.loadPage(cidEntry.Text))
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		gui.changePage(gui.withRemove(page))
	})
	return container.NewBorder(nil, nil, nil, loadBtn, cidEntry)
}

func (gui *GUI) newPageForm() fyne.CanvasObject {
	chmod := &widget.Select{
		Options:     []string{"registration", "voting"},
		PlaceHolder: "registration / voting",
	}
	chmod.OnChanged = func(mode string) {
		if mode == "registration" {
			gui.setup = rpage.NewSetupPage(gui.a, gui.is) // <- 空白
		} else {
			gui.setup = vpage.NewSetupPage(gui.a, gui.is) // <- nil pointer dereference
		}
	}
	chmod.ExtendBaseWidget(chmod)
	if gui.setup == nil {
		return chmod
	} else {
		return container.NewVBox(chmod, gui.setup)
	}
}

func (gui *GUI) defaultPage() fyne.CanvasObject {
	loadForm := gui.loadPageForm()
	newForm := gui.newPageForm()
	return container.NewBorder(loadForm, nil, nil, nil, newForm)
}

func (gui *GUI) Run() {
	gui.page = container.NewMax(gui.defaultPage())
	gui.w.SetContent(gui.page)
	gui.w.ShowAndRun()
}

func (gui *GUI) Close() {
	gui.is.Close()
}
