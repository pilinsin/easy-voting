package guiutil

import(
	cp "github.com/atotto/clipboard"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/container"
)

type copyButton struct{
	label *widget.Label
	button *widget.Button
}
func NewCopyButton(text string) *copyButton{
	label := &widget.Label{
		Text: text,
		Wrapping: fyne.TextTruncate,
	}
	label.ExtendBaseWidget(label)

	icon := theme.ContentCopyIcon()
	onTapped := func(){cp.WriteAll(label.Text)}
	btn := widget.NewButtonWithIcon("", icon, onTapped)

	return &copyButton{label, btn}
}
func (cb *copyButton) Render() fyne.CanvasObject{
	return container.NewBorder(nil, nil, nil, cb.button, cb.label)
}
func (cb *copyButton) SetText(text string){
	cb.label.SetText(text)
}