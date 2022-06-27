package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	bpage "github.com/pilinsin/easy-voting/gui/bootstrap"
	rpage "github.com/pilinsin/easy-voting/gui/registration"
	vpage "github.com/pilinsin/easy-voting/gui/voting"

	riface "github.com/pilinsin/easy-voting/registration/interface"
	viface "github.com/pilinsin/easy-voting/voting/interface"

	rgst "github.com/pilinsin/easy-voting/registration"
	vt "github.com/pilinsin/easy-voting/voting"

	evutil "github.com/pilinsin/easy-voting/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
)

func init() {
	evutil.Init()
}

func pageToTabItem(title string, page fyne.CanvasObject) *container.TabItem {
	return container.NewTabItem(title, page)
}

type GUI struct {
	rt   *i2p.I2pRouter
	bs   map[string]pv.IBootstrap
	rs   map[string]riface.IRegistration
	vs   map[string]viface.IVoting
	w    fyne.Window
	size fyne.Size
	page *fyne.Container
	tabs *container.AppTabs
}

func New(title string, width, height float32) *GUI {
	rt := i2p.NewI2pRouter()
	bs := make(map[string]pv.IBootstrap)
	rs := make(map[string]riface.IRegistration)
	vs := make(map[string]viface.IVoting)

	size := fyne.NewSize(width, height)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	win := a.NewWindow(title)
	win.Resize(size)
	page := container.NewMax()
	tabs := container.NewAppTabs()
	return &GUI{rt, bs, rs, vs, win, size, page, tabs}
}

func (gui *GUI) withRemove(page fyne.CanvasObject) fyne.CanvasObject {
	rmvBtn := widget.NewButtonWithIcon("", theme.ContentClearIcon(), func() {
		gui.tabs.Remove(gui.tabs.Selected())
	})
	return container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, page)
}

func (gui *GUI) loadPage(bAddr, cid string) (string, fyne.CanvasObject) {
	addrs := strings.Split(cid, "/")
	if len(addrs) != 2 {
		return "", nil
	}
	mode, stAddr := addrs[0], addrs[1]

	var err error
	if mode == "r" {
		baseDir := evutil.BaseDir("registration", stAddr)
		r, exist := gui.rs[baseDir]
		if !exist {
			r, err = rgst.NewRegistration(bAddr+"/"+cid, baseDir)
			if err != nil {
				return "", nil
			}
			gui.rs[baseDir] = r
		}
		return rpage.LoadPage(bAddr, cid, r)
	}
	if mode == "v" {
		baseDir := evutil.BaseDir("voting", stAddr)
		v, exist := gui.vs[baseDir]
		if !exist {
			v, err = vt.NewVoting(bAddr+"/"+cid, baseDir)
			if err != nil {
				return "", nil
			}
			gui.vs[baseDir] = v
		}
		return vpage.LoadPage(bAddr, cid, v)
	}
	return "", nil
}
func (gui *GUI) loadPageForm() fyne.CanvasObject {
	bAddrEntry := widget.NewEntry()
	bAddrEntry.SetPlaceHolder("Bootstraps Address")
	cidEntry := widget.NewEntry()
	cidEntry.SetPlaceHolder("Registration/Voting Config Cid")
	addrEntry := container.NewGridWithColumns(2, bAddrEntry, cidEntry)

	onTapped := func() {
		title, loadPage := gui.loadPage(bAddrEntry.Text, cidEntry.Text)
		bAddrEntry.SetText("")
		cidEntry.SetText("")
		if loadPage == nil {
			return
		}
		page := container.NewVScroll(loadPage)
		withRmvPage := gui.withRemove(page)
		withRmvTab := pageToTabItem(title, withRmvPage)
		gui.tabs.Append(withRmvTab)
		gui.tabs.Select(withRmvTab)
	}
	loadBtn := widget.NewButtonWithIcon("", theme.MailForwardIcon(), onTapped)

	return container.NewBorder(nil, nil, nil, loadBtn, addrEntry)
}

func (gui *GUI) newPageForm() fyne.CanvasObject {
	var setup fyne.CanvasObject
	chmod := &widget.Select{
		Options:  []string{"bootstrap", "registration", "voting"},
		Selected: "bootstrap",
	}
	chmod.OnChanged = func(mode string) {
		if mode == "registration" {
			setup = rpage.NewSetupPage(gui.w, gui.rs)
		} else if mode == "voting" {
			setup = vpage.NewSetupPage(gui.w, gui.vs)
		} else {
			setup = bpage.NewSetupPage(gui.bs)
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
		if err := gui.rt.Start(); err == nil {
			i2pNote.SetText("i2p router on")
		} else {
			gui.initErrorPage()
		}
	}()
}
func (gui *GUI) Close() {
	for _, v := range gui.vs {
		v.Close()
	}
	for _, r := range gui.rs {
		r.Close()
	}
	for _, b := range gui.bs {
		b.Close()
	}
	gui.rt.Stop()
}

func (gui *GUI) Run() {
	i2pNote := widget.NewLabel("i2p router setup...")
	gui.i2pStart(i2pNote)

	gui.tabs.Append(gui.defaultPage())
	loadForm := gui.loadPageForm()
	gui.page.Add(container.NewBorder(loadForm, i2pNote, nil, nil, gui.tabs))

	gui.w.SetContent(gui.page)
	gui.w.SetOnClosed(gui.Close)
	gui.w.ShowAndRun()
}
