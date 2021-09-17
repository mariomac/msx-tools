package img

import (
	"errors"
	"image"
	"image/color"
	_ "image/color/palette"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
)

type Format string
const (
	FormatPNG Format = "png"
)

type Bitmap struct {
	img image.Image
	palette color.Palette
}

func Load(r io.Reader, palette color.Palette) (Bitmap, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return Bitmap{}, err
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	img2 := image.NewPaletted(image.Rect(0, 0, w, h), palette)
	for y := 0 ; y < h ; y++ {
		for x := 0 ; x < w ; x++ {
			img2.Set(x, y, img.At(x, y))
		}
	}
	return Bitmap{img: img2, palette: palette}, nil
}

func (b *Bitmap) Save(w io.Writer, format Format) error {
	if format != FormatPNG {
		return errors.New("TODO: only PNG save allowed")
	}
	return png.Encode(w, b.img)
}
