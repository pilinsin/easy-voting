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
	res := fyne.NewStaticResource(cand.Name, cand.Image)
	img := canvas.NewImageFromResource(res)
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

func newCandEntry(w fyne.Window) *candEntry {
	imgBtn := &widget.Button{Icon: theme.ContentAddIcon()}
	imgBtn.OnTapped = imageDialog(w, imgBtn.Icon)
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
	var icon fyne.Resource
	defIcon := theme.ContentAddIcon()
	selected := ce.imgBtn.Icon
	if resourceEqual(selected, defIcon) {
		icon = theme.DeleteIcon()
	} else {
		icon = selected
	}
	return vutil.Candidate{
		Name:  ce.name.Text,
		Group: ce.group.Text,
		Url:   ce.url.Text,
		Image: icon.Content(),
	}
}
