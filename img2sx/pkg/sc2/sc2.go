package sc2

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/nfnt/resize"

	"github.com/mariomac/msxtools/img2sx/pkg/img"
)

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

var inversePalette = map[img.RGB]uint8{
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

// Screen 2 VRAM addresses
const (
	patternGen1  = 0x0000
	patternGen2  = 0x0800
	patternGen3  = 0x1000
	patternName1 = 0x1800
	patternName2 = 0x1900
	patternName3 = 0x1a00
	spriteAttrs  = 0x1b00
	palette      = 0x1b80
	color1       = 0x2000
	color2       = 0x2800
	color3       = 0x3000
	spriteGen    = 0x3800
)

// screen size
const (
	PixelsWidth  = 256
	PixelsHeight = 192
)

const (
	// number of patterns per tile
	tilePatterns = 8
	// number of tiles for the each pattern name table
	tableTiles = 256
	// number of patterns (actual bytes) for the pattern generator & color tables
	tablePatterns = 0x800
)

func toRGB(c color.Color) img.RGB {
	r, g, b, a := c.RGBA()

	r *= 0xff
	r /= a

	g *= 0xff
	g /= a

	b *= 0xff
	b /= a

	return img.RGB(r<<16 | g<<8 | b)
}

// Pattern represent the 8-pixel line of a tile
type Pattern struct {
	// Bitmap where each bit represents a pixel
	Bitmap uint8
	// Color in :
	// bits 3-0: Background color
	// bits 7-4: Foreground color
	Color uint8
}

func sample(bitmap image.Image, x, y int) Pattern {
	// count the 2 most frequent colors
	var frequency [16]int
	mf, mf2 := -1, -1
	for i := 0; i < 8; i++ {
		// TODO: replace color 0 by color 1
		cn := inversePalette[toRGB(bitmap.At(x+i, y))]
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
		px := toRGB(bitmap.At(x+i, y))

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

// Tile is formed by 8 patterns, representing an 8x8 pixels tile
type Tile [tilePatterns]Pattern

// TileSet keeps the data of a whole Screen 2 image
type TileSet struct {
	// todo: make private
	Table [3][]Tile
	Names [3][]uint8
}

func Convert(img image.Image, opt ConvertOpt) *TileSet {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	if w != PixelsWidth || h != PixelsHeight {
		switch opt {
		case Crop:
			// do nothing
		case Stretch:
			// TODO: allow an option to set resize type
			img = resize.Resize(PixelsWidth, PixelsHeight, img, resize.Bilinear)
		case KeepAspect:
			if float64(w)/float64(h) > float64(PixelsWidth)/float64(PixelsHeight) {
				img = resize.Resize(PixelsWidth, 0, img, resize.Bilinear)
			} else {
				img = resize.Resize(0, PixelsHeight, img, resize.Bilinear)
			}
		}
	}
	// migrating to 16 color without alpha
	// intentionally omitting zero color
	sc2Bounds := image.Rect(0, 0, PixelsWidth, PixelsHeight)
	paletted := image.NewPaletted(sc2Bounds, Palette[1:])
	draw.Draw(paletted, sc2Bounds, img, image.Pt(0, 0), draw.Src)
	colorNormalized := image.NewNRGBA(sc2Bounds)
	draw.Draw(colorNormalized, sc2Bounds, paletted, image.Pt(0, 0), draw.Src)

	ts := TileSet{}
	// TODO: avoid repeating tiles
	for table := 0; table < 3; table++ {
		name := 0
		for y := table * 64; y < table*64+64; y += 8 {
			for x := 0; x < 255; x += 8 {
				t := Tile{}
				for dy := 0; dy < 8; dy++ {
					t[dy] = sample(colorNormalized, x, y+dy)
				}
				ts.Names[table] = append(ts.Names[table], uint8(name))
				ts.Table[table] = append(ts.Table[table], t)
				name++
			}
		}
	}
	return &ts
}

func (s *TileSet) ColorModel() color.Model {
	return Palette
}

func (s *TileSet) Bounds() image.Rectangle {
	// todo: global constant
	return image.Rect(0, 0, 256, 192)
}

func (s *TileSet) At(x, y int) color.Color {
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
		return Palette[pattern.Color&0b1111]
	} else {
		// return foreground color
		return Palette[pattern.Color>>4]
	}
}
