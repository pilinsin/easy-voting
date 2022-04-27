package gui

import (
	"context"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	bpage "github.com/pilinsin/easy-voting/gui/bootstrap"
	rpage "github.com/pilinsin/easy-voting/gui/registration"
	vpage "github.com/pilinsin/easy-voting/gui/voting"

	evutil "github.com/pilinsin/easy-voting/util"
)

func init(){
	evutil.Init()
}

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

func (gui *GUI) withRemove(page fyne.CanvasObject, closer func()) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		closer()
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

func (gui *GUI) loadPage(ctx context.Context, addr, idStr string) (fyne.CanvasObject, func()) {
	if ok := strings.HasPrefix(addr, "r/"); ok{
		return rpage.LoadPage(ctx, addr, idStr)
	}
	if ok := strings.HasPrefix(addr, "v/"); ok{
		return vpage.LoadPage(ctx, addr, idStr)
	}
	return nil, nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	addrEntry := widget.NewEntry()
	addrEntry.PlaceHolder = "Registration/Voting Config Address"
	idEntry := widget.NewEntry()
	idEntry.PlaceHolder = "User/Manager Identity Address"
	
	onTapped := func(){
		loadPage, closer := gui.loadPage(context.Background(), addrEntry.Text, idEntry.Text)
		if loadPage == nil {
			addrEntry.SetText("")
			idEntry.SetText("")
			return
		}
		page := container.NewVScroll(loadPage)
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		gui.changePage(gui.withRemove(page, closer))
	}
	loadBtn := widget.NewButtonWithIcon("", theme.MailForwardIcon(), onTapped)

	entries := container.NewVBox(addrEntry, idEntry)
	return container.NewBorder(nil, nil, nil, loadBtn, entries)
}

func (gui *GUI) newPageForm() fyne.CanvasObject {
	var setup fyne.CanvasObject
	chmod := &widget.Select{
		Options:  []string{"registration", "voting", "bootstrap"},
		Selected: "registration",
	}
	chmod.OnChanged = func(mode string) {
		if mode == "registration" {
			setup = rpage.NewSetupPage(gui.w)
		} else if mode == "voting"{
			setup = vpage.NewSetupPage(gui.w)
		} else{
			setup = bpage.NewSetupPage(gui.w)
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


