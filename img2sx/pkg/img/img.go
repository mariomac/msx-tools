package img

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	_ "image/color/palette"
	"image/draw"
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
	Img *image.NRGBA
}

func Load(r io.Reader, palette color.Palette) (Bitmap, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return Bitmap{}, err
	}
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	resized := resize.Resize(256,0, img, resize.Bilinear)
	paletted := image.NewPaletted(image.Rect(0, 0, w, h), palette)
	draw.Draw(paletted, img.Bounds(), resized, image.Pt(0, 0), draw.Src)
	colorNormalized := image.NewNRGBA(img.Bounds())
	draw.Draw(colorNormalized, img.Bounds(), paletted, image.Pt(0, 0), draw.Src)
	return Bitmap{Img: colorNormalized}, nil
}

// TODO: remove
func (b *Bitmap) Save(w io.Writer, format Format) error {
	if format != FormatPNG {
		return errors.New("TODO: only PNG save allowed")
	}
	return png.Encode(w, b.Img)
}

func (b *Bitmap) ColorModel() color.Model {
	return b.Img.ColorModel()
}

func (b *Bitmap) Bounds() image.Rectangle {
	return b.Img.Bounds()
}

func (b *Bitmap) At(x, y int) color.Color {
	// TODO: convert on fly or cache conversions?
	return b.Img.At(x, y)
}

func (b *Bitmap) RGBAt(x, y int) RGB {
	p := b.Img.At(x, y).(color.NRGBA)
	return RGB(uint32(p.R)<<16 | uint32(p.G)<<8 | uint32(p.B))
}
