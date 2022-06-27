package votingpage

import (
	"context"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	vutil "github.com/pilinsin/easy-voting/voting/util"
)

func CandCards(cands []*vutil.Candidate) fyne.CanvasObject {
	var candList [](fyne.CanvasObject)
	for _, cand := range cands {
		card := newCandCard(cand)
		candList = append(candList, card.Render())
	}
	return container.NewAdaptiveGrid(4, candList...)
}

type candCard struct {
	name      *widget.Label
	group     *widget.Label
	url       fyne.CanvasObject
	imgCanvas *imageCanvas
}

func newCandCard(cand *vutil.Candidate) *candCard {
	res := fyne.NewStaticResource(cand.ImgName, cand.Image)
	img := newImageCanvas(res)
	if gutil.ResourceEqual(res, gutil.DefaultIcon()) {
		img.SetImage(nil)
	}

	card := &candCard{
		name:      widget.NewLabel(cand.Name),
		group:     widget.NewLabel(cand.Group),
		url:       gutil.SetUrl("URL", cand.Url),
		imgCanvas: img,
	}
	return card
}
func (cc *candCard) Render() fyne.CanvasObject {
	return container.NewVBox(cc.imgCanvas.Render(), cc.name, cc.group, cc.url)
}

type CandForm struct {
	cands []*candEntry
}

func NewCandForm() *CandForm {
	return &CandForm{}
}
func (cf *CandForm) Render(w fyne.Window) fyne.CanvasObject {
	contents := container.NewAdaptiveGrid(4)
	//AddButton (CandEntry with RemoveButton)
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		cand := newCandEntry(w)
		cf.cands = append(cf.cands, cand)
		rmvBtn := &widget.Button{Icon: theme.ContentClearIcon()}
		withRmvBtn := container.NewBorder(container.NewBorder(nil, nil, nil, rmvBtn), nil, nil, nil, cand.Render())
		rmvBtn.OnTapped = func() {
			contents.Remove(withRmvBtn)
			cf.cands = removeCandEntry(cf.cands, cand)
		}
		rmvBtn.ExtendBaseWidget(rmvBtn)
		contents.Add(withRmvBtn)
	})
	return container.NewBorder(container.NewBorder(nil, nil, addBtn, nil), nil, nil, nil, contents)
}
func (cf *CandForm) Candidates() []*vutil.Candidate {
	candidates := make([]*vutil.Candidate, len(cf.cands))
	for idx, cand := range cf.cands {
		candidates[idx] = cand.Candidate()
	}
	return candidates
}
func removeCandEntry(cands []*candEntry, cand *candEntry) []*candEntry {
	newCandEntries := []*candEntry{}
	for _, item := range cands {
		//different pointer
		if item != cand {
			newCandEntries = append(newCandEntries, item)
		}
	}
	return newCandEntries
}

type candEntry struct {
	name      *widget.Entry
	group     *widget.Entry
	url       *widget.Entry
	imgBtn    *widget.Button
	imgCanvas *imageCanvas
	thumbnail chan fyne.Resource
}

func newCandEntry(w fyne.Window) *candEntry {
	imgCanvas := newImageCanvas(gutil.DefaultIcon())
	imgCanvas.Hide()

	thumb := make(chan fyne.Resource)
	imgBtn := &widget.Button{Icon: theme.ContentAddIcon()}
	imgBtn.OnTapped = func() {
		imageDialog(w, thumb, imgBtn, imgCanvas)()
		//fmt.Println(thumb.Name())
		//fmt.Println(resourceEqual(thumb, theme.ContentAddIcon()))
	}
	imgBtn.ExtendBaseWidget(imgBtn)

	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Name")
	groupEntry := widget.NewEntry()
	groupEntry.SetPlaceHolder("Group")
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("URL")
	cand := &candEntry{
		name:      nameEntry,
		group:     groupEntry,
		url:       urlEntry,
		imgBtn:    imgBtn,
		imgCanvas: imgCanvas,
		thumbnail: thumb,
	}
	return cand
}
func (ce *candEntry) Render() fyne.CanvasObject {
	return container.NewVBox(ce.imgCanvas.Render(), ce.imgBtn, ce.name, ce.group, ce.url)
}
func (ce *candEntry) Candidate() *vutil.Candidate {
	var res fyne.Resource
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	select {
	case res = <-ce.thumbnail:
	case <-ctx.Done():
		res = gutil.DefaultIcon()
		close(ce.thumbnail)
	}

	return &vutil.Candidate{
		Name:    ce.name.Text,
		Group:   ce.group.Text,
		Url:     ce.url.Text,
		Image:   res.Content(),
		ImgName: res.Name(),
	}
}

func imageDialog(w fyne.Window, res chan<- fyne.Resource, hideObj fyne.CanvasObject, resCanvas *imageCanvas) func() {
	return func() {
		onSelected := func(rc fyne.URIReadCloser, err error) {
			if rc == nil || err != nil {
				return
			}
			if ok := isImageExtension(rc.URI().Extension()); !ok {
				return
			}

			data, err := io.ReadAll(rc)
			if err != nil {
				return
			}
			loadRes := &fyne.StaticResource{
				StaticName:    rc.URI().Name(),
				StaticContent: data,
			}
			go func() {
				res <- loadRes
				close(res)
			}()
			hideObj.Hide()
			resCanvas.Show()
			resCanvas.SetImage(loadRes)
		}
		dialog.ShowFileOpen(onSelected, w)
	}
}
func isImageExtension(ext string) bool {
	switch ext {
	case ".bmp", ".png", ".jpeg", ".jpg", ".gif", ".tiff", ".vp8l", ".webp", ".svg":
		return true
	default:
		return false
	}
}

type imageCanvas struct {
	imgCanvas *fyne.Container
}

func newImageCanvas(res fyne.Resource) *imageCanvas {
	imgCanvas := canvas.NewImageFromResource(res)
	imgCanvas.FillMode = canvas.ImageFillContain
	imgGridCanvas := container.NewGridWrap(fyne.NewSize(169, 239.27), imgCanvas)
	return &imageCanvas{imgGridCanvas}
}
func (iCanvas *imageCanvas) Render() fyne.CanvasObject {
	return iCanvas.imgCanvas
}
func (iCanvas *imageCanvas) Hide() {
	iCanvas.imgCanvas.Hide()
}
func (iCanvas *imageCanvas) Show() {
	iCanvas.imgCanvas.Show()
}
func (iCanvas *imageCanvas) SetImage(res fyne.Resource) {
	imgCanvas, _ := iCanvas.imgCanvas.Objects[0].(*canvas.Image)
	imgCanvas.Resource = res
	imgCanvas.Refresh()
}
