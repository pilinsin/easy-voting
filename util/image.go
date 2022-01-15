package util

import(
	"io"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/vp8l"
	_ "golang.org/x/image/webp"
)

func LoadImage(r io.Reader) (image.Image, error){
	img, _, err := image.Decode(r)
	return img, err
}
func LoadImageConfig(r io.Reader) (image.Config, error){
	cfg, _, err := image.DecodeConfig(r)
	return cfg, err
}
