package img

import (
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/color/palette"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

type Format string

const (
	FormatPNG Format = "png"
)

func Load(r io.Reader, palette color.Palette) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	resized := resize.Resize(256,0, img, resize.Bilinear)
	paletted := image.NewPaletted(image.Rect(0, 0, w, h), palette)
	draw.Draw(paletted, img.Bounds(), resized, image.Pt(0, 0), draw.Src)
	colorNormalized := image.NewNRGBA(img.Bounds())
	draw.Draw(colorNormalized, img.Bounds(), paletted, image.Pt(0, 0), draw.Src)
	return colorNormalized, nil
}

