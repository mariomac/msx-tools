package sc2

import (
	"github.com/mariomac/msxtools/img2sx/pkg/internal/screen2"
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"
)

const (
	// number of patterns (actual bytes) for the pattern generator & color tables
	tablePatterns = 0x800
)

// Image keeps the data of a whole Screen 2 image
type Image struct {
	// todo: make private
	Table [3][]screen2.Tile
	Names [3][]uint8
}

func convert(img image.Image, opt ConvertOpt) *Image {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w != screen2.PixelsWidth || h != screen2.PixelsHeight {
		switch opt {
		case Crop:
			// do nothing
		case Stretch:
			// TODO: allow an option to set resize type
			img = resize.Resize(screen2.PixelsWidth, screen2.PixelsHeight, img, resize.Bilinear)
		case KeepAspect:
			if float64(w)/float64(h) > float64(screen2.PixelsWidth)/float64(screen2.PixelsHeight) {
				img = resize.Resize(screen2.PixelsWidth, 0, img, resize.Bilinear)
			} else {
				img = resize.Resize(0, screen2.PixelsHeight, img, resize.Bilinear)
			}
		}
	}
	// migrating to 16 color without alpha
	// intentionally omitting zero color
	sc2Bounds := image.Rect(0, 0, screen2.PixelsWidth, screen2.PixelsHeight)
	paletted := image.NewPaletted(sc2Bounds, screen2.Palette[1:])
	draw.Draw(paletted, sc2Bounds, img, image.Pt(0, 0), draw.Src)
	colorNormalized := image.NewNRGBA(sc2Bounds)
	draw.Draw(colorNormalized, sc2Bounds, paletted, image.Pt(0, 0), draw.Src)

	ts := Image{}
	// TODO: avoid repeating tiles
	for table := 0; table < 3; table++ {
		name := 0
		for y := table * 64; y < table*64+64; y += 8 {
			for x := 0; x < 255; x += 8 {
				t := screen2.Tile{}
				for dy := 0; dy < 8; dy++ {
					t[dy] = screen2.Sample(colorNormalized, x, y+dy)
				}
				ts.Names[table] = append(ts.Names[table], uint8(name))
				ts.Table[table] = append(ts.Table[table], t)
				name++
			}
		}
	}
	return &ts
}

func (s *Image) ColorModel() color.Model {
	return screen2.Palette
}

func (s *Image) Bounds() image.Rectangle {
	// todo: global constant
	return image.Rect(0, 0, 256, 192)
}

func (s *Image) At(x, y int) color.Color {
	// select absolute tile number from x, y
	screenTile := (y/8)*32 + x/8
	// select position in tile names table
	namesTable := 0
	// use logical operations to avoid loop
	for screenTile >= 256 {
		namesTable++
		screenTile -= 256
	}
	pattern := s.Table[namesTable][screenTile][y%8]
	pixelMask := pattern.Bitmap & (1 << (7 - x%8))
	if pixelMask == 0 {
		// return background color
		return screen2.Palette[pattern.Color&0b1111]
	} else {
		// return foreground color
		return screen2.Palette[pattern.Color>>4]
	}
}
