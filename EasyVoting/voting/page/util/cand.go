package votingpageutil

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	vutil "EasyVoting/voting/util"
)

func CandCards(cands []vutil.Candidate) fyne.CanvasObject {
	var candList [](fyne.CanvasObject)
	for _, cand := range cands {
		card := candCard(cand)
		content := container.NewGridWrap(fyne.NewSize(120.0, 268.0), card)
		candList = append(candList, content)
	}
	return container.NewAdaptiveGrid(4, candList...)
}
func candCard(cand vutil.Candidate) *widget.Card {
	img := canvas.NewImageFromResource(cand.Image)
	img.FillMode = canvas.ImageFillContain

	card := &widget.Card{
		Title:    cand.Name,
		Subtitle: cand.Group,
		Content:  setUrl("URL", cand.Url),
	}
	card.ExtendBaseWidget(card)
	card.SetImage(img)

	return card
}

type CandForm struct {
	fyne.CanvasObject
	cands []*candEntry
}

func NewCandForm(a fyne.App) *CandForm {
	cf := &CandForm{}
	contents := container.NewAdaptiveGrid(4, nil)
	//AddButton (CandEntry with RemoveButton)
	addBtn := widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		cand := newCandEntry(a)
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
	cf.CanvasObject = container.NewBorder(container.NewBorder(nil, nil, addBtn, nil), nil, nil, nil, contents)
	return cf
}
func (cf *CandForm) Candidates() []vutil.Candidate {
	candidates := make([]vutil.Candidate, len(cf.cands))
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
	name   *widget.Entry
	group  *widget.Entry
	url    *widget.Entry
	imgBtn *widget.Button
}

func newCandEntry(a fyne.App) *candEntry {
	imgBtn := &widget.Button{Icon: theme.ContentAddIcon()}
	imgBtn.OnTapped = imageDialog(a, imgBtn.Icon)
	imgBtn.ExtendBaseWidget(imgBtn)

	cand := &candEntry{
		name:   PlaceHolderEntry("Name"),
		group:  PlaceHolderEntry("Group"),
		url:    PlaceHolderEntry("URL"),
		imgBtn: imgBtn,
	}
	return cand
}
func (ce *candEntry) Render() fyne.CanvasObject {
	return container.NewVBox(ce.imgBtn, ce.name, ce.group, ce.url)
}
func (ce *candEntry) Candidate() vutil.Candidate {
	return vutil.Candidate{
		Name:  ce.name.Text,
		Group: ce.group.Text,
		Url:   ce.url.Text,
		Image: ce.imgBtn.Icon,
	}
}
