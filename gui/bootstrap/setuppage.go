package bootstrappage

import (
	peer "github.com/libp2p/go-libp2p-core/peer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	gutil "github.com/pilinsin/easy-voting/gui/util"
	i2p "github.com/pilinsin/go-libp2p-i2p"
	pv "github.com/pilinsin/p2p-verse"
)

func NewSetupPage(bs map[string]pv.IBootstrap) fyne.CanvasObject {
	baddrsLabel := gutil.NewCopyButton("bootstrap list address")

	form := NewBootstrapsForm()
	addrsBtn := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		var self pv.IBootstrap
		var err error
		mapKey := "setup"
		
		baddrsLabel.SetText("processing...")
		baddrs := form.AddrInfos()
		if b, exist := bs[mapKey]; exist{
			self = b
		}
		self, err = pv.NewBootstrap(i2p.NewI2pHost, baddrs...)
		if err != nil{
			baddrsLabel.SetText("bootstrap list address")
			return
		}
		bs[mapKey] = self

		baddrs = append(baddrs, self.AddrInfo())
		s := pv.AddrInfosToString(baddrs...)
		if s == "" {
			baddrsLabel.SetText("bootstrap list address")
		} else {
			baddrsLabel.SetText(s)
		}
	})

	baddrs := container.NewBorder(nil, nil, addrsBtn, nil, baddrsLabel.Render())
	return container.NewVBox(form.Render(), baddrs)
}

func mapToSlice(m map[string]peer.AddrInfo) []peer.AddrInfo {
	ais := make([]peer.AddrInfo, len(m))
	idx := 0
	for _, v := range m {
		ais[idx] = v
		idx++
	}
	return ais
}
func sliceToMap(ais []peer.AddrInfo) map[string]peer.AddrInfo {
	m := make(map[string]peer.AddrInfo)
	for _, ai := range ais {
		s := pv.AddrInfoToString(ai)
		if s != "" {
			m[s] = ai
		}
	}
	return m
}

type bootstrapsForm struct {
	*gutil.RemovableEntryForm
}

func NewBootstrapsForm() *bootstrapsForm {
	ref := gutil.NewRemovableEntryForm()
	return &bootstrapsForm{ref}
}
func (bf *bootstrapsForm) AddrInfos() []peer.AddrInfo {
	txts := bf.Texts()
	aiMap := make(map[string]peer.AddrInfo)

	for _, txt := range txts {
		ai := pv.AddrInfoFromString(txt)
		if ai.ID != "" && len(ai.Addrs) > 0 {
			aiMap[txt] = ai
		} else {
			ais := pv.AddrInfosFromString(txt)
			for _, ai := range ais {
				if ai.ID == "" || len(ai.Addrs) == 0 {
					continue
				}
				s := pv.AddrInfoToString(ai)
				aiMap[s] = ai
			}
		}
	}

	return mapToSlice(aiMap)
}
