package guiutil

import(
	"bytes"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
)


func SetUrl(text string, urlStr string) fyne.CanvasObject {
	if urlStr != "" {
		parsedUrl, err := url.Parse(urlStr)
		if err == nil {
			return widget.NewHyperlink(text, parsedUrl)
		}
	}
	//When invalid urlStr, err is not raised.
	return widget.NewLabel("")
}

func DefaultIcon() fyne.Resource{
	return theme.FyneLogo()
}
func ResourceEqual(selected, def fyne.Resource) bool {
	name := selected.Name() == def.Name()
	content := bytes.Equal(selected.Content(), def.Content())
	return name && content
}