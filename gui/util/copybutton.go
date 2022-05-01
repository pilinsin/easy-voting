package guiutil

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	cp "github.com/atotto/clipboard"
)

type CopyButton struct {
	label  *widget.Label
	button *widget.Button
}

func NewCopyButton(text string) *CopyButton {
	label := &widget.Label{
		Text:     text,
		Wrapping: fyne.TextTruncate,
	}
	label.ExtendBaseWidget(label)

	icon := theme.ContentCopyIcon()
	onTapped := func() { cp.WriteAll(label.Text) }
	btn := widget.NewButtonWithIcon("", icon, onTapped)

	return &CopyButton{label, btn}
}
func (cb *CopyButton) Render() fyne.CanvasObject {
	return container.NewBorder(nil, nil, nil, cb.button, cb.label)
}
func (cb *CopyButton) SetText(text string) {
	cb.label.SetText(text)
}
