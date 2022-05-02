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
	i2p "github.com/pilinsin/go-libp2p-i2p"
)

func init() {
	evutil.Init()
}

func pageToTabItem(title string, page fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(title, page)
}

type GUI struct {
	w    fyne.Window
	size fyne.Size
	page *fyne.Container
	tabs *container.AppTabs
}

func New(title string, width, height float32) *GUI {
	size := fyne.NewSize(width, height)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	win := a.NewWindow(title)
	win.Resize(size)
	page := container.NewMax()
	tabs := container.NewAppTabs()
	return &GUI{win, size, page, tabs}
}

func (gui *GUI) withRemove(page fyne.CanvasObject, closer func()) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		closer()
		gui.tabs.Remove(gui.tabs.Selected())
	})
	return container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, page)
}

func (gui *GUI) loadPage(ctx context.Context, addr, idStr string) (string, fyne.CanvasObject, func()) {
	if ok := strings.HasPrefix(addr, "r/"); ok {
		return rpage.LoadPage(ctx, addr, idStr)
	}
	if ok := strings.HasPrefix(addr, "v/"); ok {
		return vpage.LoadPage(ctx, addr, idStr)
	}
	return "", nil, nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	addrEntry := widget.NewEntry()
	addrEntry.PlaceHolder = "Registration/Voting Config Address"
	idEntry := widget.NewEntry()
	idEntry.PlaceHolder = "User/Manager Identity Address"

	onTapped := func() {
		title, loadPage, closer := gui.loadPage(context.Background(), addrEntry.Text, idEntry.Text)
		addrEntry.SetText("")
		idEntry.SetText("")
		if loadPage == nil {
			return
		}
		page := container.NewVScroll(loadPage)
		//page.SetMinSize(fyne,NewSize(101.1,201.2))
		withRmvPage := gui.withRemove(page, closer)
		withRmvTab := pageToTabItem(title, withRmvPage)
		gui.tabs.Append(withRmvTab)
		gui.tabs.Select(withRmvTab)
	}
	loadBtn := widget.NewButtonWithIcon("", theme.MailForwardIcon(), onTapped)

	entries := container.NewVBox(addrEntry, idEntry)
	return container.NewBorder(nil, nil, nil, loadBtn, entries)
}

func (gui *GUI) newPageForm() fyne.CanvasObject {
	var setup fyne.CanvasObject
	chmod := &widget.Select{
		Options:  []string{"bootstrap", "registration", "voting"},
		Selected: "bootstrap",
	}
	chmod.OnChanged = func(mode string) {
		if mode == "registration" {
			setup = rpage.NewSetupPage(gui.w)
		} else if mode == "voting" {
			setup = vpage.NewSetupPage(gui.w)
		} else {
			setup = bpage.NewSetupPage(gui.w)
		}
		newForm := container.NewBorder(chmod, nil, nil, nil, setup)
		newTab := pageToTabItem("setup", newForm)
		gui.tabs.Items[0] = newTab
		gui.tabs.Refresh()
	}
	chmod.ExtendBaseWidget(chmod)

	return container.NewBorder(chmod, nil, nil, nil)
}

func (gui *GUI) defaultPage() *container.TabItem {
	newForm := gui.newPageForm()
	return pageToTabItem("setup", newForm)
}

func (gui *GUI) initErrorPage() {
	for _, obj := range gui.page.Objects {
		gui.page.Remove(obj)
	}
	failed := widget.NewLabel("i2p router failed to start. please try again later.")
	gui.page.Add(failed)
	gui.page.Refresh()
}
func (gui *GUI) i2pStart(i2pNote *widget.Label) {
	go func() {
		if err := i2p.StartI2pRouter(); err == nil {
			i2pNote.SetText("i2p router on")
		} else {
			gui.initErrorPage()
		}
	}()
}

func (gui *GUI) Run() {
	i2pNote := widget.NewLabel("i2p router setup...")
	gui.i2pStart(i2pNote)

	gui.tabs.Append(gui.defaultPage())
	loadForm := gui.loadPageForm()
	gui.page.Add(container.NewBorder(loadForm, i2pNote, nil, nil, gui.tabs))

	gui.w.SetContent(gui.page)
	gui.w.SetOnClosed(i2p.StopI2pRouter)
	gui.w.ShowAndRun()
}
