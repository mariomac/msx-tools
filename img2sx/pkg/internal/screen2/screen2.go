package screen2

import (
	"github.com/mariomac/msxtools/img2sx/pkg/internal/img"
	"image"
	"image/color"
)

// Screen 2 VRAM addresses
const (
	AddrSpriteAttrs = 0x1b00
	AddrColorTable1 = 0x2000
	AddrSpriteGen   = 0x3800
)

// screen size
const (
	PixelsWidth  = 256
	PixelsHeight = 192
)

const (
	// number of patterns per tile
	TilePatterns = 8
	// number of tiles for the each pattern name table
	TableTiles = 256
)

// Tile is formed by 8 patterns, representing an 8x8 pixels tile
type Tile [TilePatterns]Pattern

// Pattern represent the 8-pixel line of a tile
type Pattern struct {
	// Bitmap where each bit represents a pixel
	Bitmap uint8
	// Color in :
	// bits 3-0: Background color
	// bits 7-4: Foreground color
	Color uint8
}

var Palette = color.Palette{
	img.RGB(0x000000),
	img.RGB(0x010101),
	img.RGB(0x3eb849),
	img.RGB(0x74d07d),
	img.RGB(0x5955e0),
	img.RGB(0x8076f1),
	img.RGB(0xb95e51),
	img.RGB(0x65dbef),
	img.RGB(0xdb6559),
	img.RGB(0xff897d),
	img.RGB(0xccc35e),
	img.RGB(0xded087),
	img.RGB(0x3aa241),
	img.RGB(0xb766b5),
	img.RGB(0xcccccc),
	img.RGB(0xffffff),
}

var InversePalette = map[img.RGB]uint8{
	//0x000000: 0,
	0x010101: 1,
	0x3eb849: 2,
	0x74d07d: 3,
	0x5955e0: 4,
	0x8076f1: 5,
	0xb95e51: 6,
	0x65dbef: 7,
	0xdb6559: 8,
	0xff897d: 9,
	0xccc35e: 10,
	0xded087: 11,
	0x3aa241: 12,
	0xb766b5: 13,
	0xcccccc: 14,
	0xffffff: 15,
}

func ToRGB(c color.Color) img.RGB {
	r, g, b, a := c.RGBA()

	r *= 0xff
	r /= a

	g *= 0xff
	g /= a

	b *= 0xff
	b /= a

	return img.RGB(r<<16 | g<<8 | b)
}

func Sample(bitmap image.Image, x, y int) Pattern {
	// count the 2 most frequent colors
	var frequency [16]int
	mf, mf2 := -1, -1
	for i := 0; i < 8; i++ {
		// TODO: replace color 0 by color 1
		cn := InversePalette[ToRGB(bitmap.At(x+i, y))]
		frequency[cn]++
		if mf < 0 || frequency[cn] > frequency[mf] {
			mf2 = mf
			mf = int(cn)
		} else if int(cn) != mf && (mf2 < 0 || frequency[cn] > frequency[mf2]) {
			mf2 = int(cn)
		}
	}
	// foreground is the most frequent color, background the second most frequent color
	// build the bitmap as a function of the closer colors to background or foreground
	bmp := uint8(0)
	for i := 0; i < 8; i++ {
		bmp <<= 1
		b := uint8(1)
		px := ToRGB(bitmap.At(x+i, y))

		if mf2 >= 0 {
			dmf := px.DistanceTo(Palette[mf].(img.RGB))
			dmf2 := px.DistanceTo(Palette[mf2].(img.RGB))
			if dmf == dmf2 && i%2 == 0 || dmf2 < dmf {
				b = 0
			}
		}
		bmp |= b
	}
	return Pattern{
		Bitmap: bmp,
		Color:  uint8(mf)<<4 | uint8(mf2&0b1111),
	}
}

