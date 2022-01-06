package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"EasyVoting/ipfs"
	rpage "EasyVoting/registration/page"
	rputil "EasyVoting/registration/page/util"
	rutil "EasyVoting/registration/util"
	vpage "EasyVoting/voting/page"
	vutil "EasyVoting/voting/util"
)

type GUI struct {
	a    fyne.App
	w    fyne.Window
	page *fyne.Container
	is   *ipfs.IPFS
}

func New(title string, is *ipfs.IPFS) *GUI {
	a := app.New()
	win := a.NewWindow(title)
	page := container.NewMax()
	return &GUI{a, win, page, is}
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
}

func (gui *GUI) loadPage(cid string) (fyne.CanvasObject, rputil.IPageCloser) {
	if _, err := rutil.ConfigFromCid(cid, gui.is); err == nil {
		fmt.Println("rCfg")
		return rpage.LoadPage(cid, gui.is)
	}
	if _, err := rutil.ManConfigFromCid(cid, gui.is); err == nil {
		fmt.Println("rmCfg")
		return rpage.LoadManPage(cid, gui.is)
	}
	if _, err := vutil.ConfigFromCid(cid, gui.is); err == nil {
		fmt.Println("vCfg")
		return vpage.LoadPage(cid, gui.is)
	}
	if _, err := vutil.ManConfigFromCid(cid, gui.is); err == nil {
		fmt.Println("vmCfg")
		return vpage.LoadManPage(cid, gui.is)
	}
	return nil, nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	cidEntry := widget.NewEntry()
	cidEntry.PlaceHolder = "page CID (Qm...)"
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

func (gui *GUI) Close() {
	gui.is.Close()
}
