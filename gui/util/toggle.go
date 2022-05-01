package guiutil

import (
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type toggleButton struct {
	*widget.Button
	onoff bool
}

func NewToggleButton(onEnabled, onDisabled func() error) *toggleButton {
	tbtn := &toggleButton{onoff: false}

	btn := &widget.Button{
		Icon: theme.RadioButtonIcon(),
	}
	btn.OnTapped = func() {
		if tbtn.onoff {
			//on -> off
			if err := onDisabled(); err != nil {
				return
			}
			tbtn.onoff = false
			btn.Icon = theme.RadioButtonIcon()
			btn.Refresh()
		} else {
			//off -> on
			if err := onEnabled(); err != nil {
				return
			}
			tbtn.onoff = true
			btn.Icon = theme.RadioButtonCheckedIcon()
			btn.Refresh()
		}
	}

	tbtn.Button = btn
	tbtn.ExtendBaseWidget(tbtn)
	return tbtn
}
