package sc2

import (
	"image/color"
)

var Palette = color.Palette{
	int2nrgba(0x000000),
	int2nrgba(0x010101),
	int2nrgba(0x3eb849),
	int2nrgba(0x74d07d),
	int2nrgba(0x5955e0),
	int2nrgba(0x8076f1),
	int2nrgba(0xb95e51),
	int2nrgba(0x65dbef),
	int2nrgba(0xdb6559),
	int2nrgba(0xff897d),
	int2nrgba(0xccc35e),
	int2nrgba(0xded087),
	int2nrgba(0x3aa241),
	int2nrgba(0xb766b5),
	int2nrgba(0xcccccc),
	int2nrgba(0xffffff),
}

func int2nrgba(i uint32) color.NRGBA {
	return color.NRGBA{R: uint8(i >> 16), G: uint8(i >> 8), B: uint8(i), A: 0xFF}
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

// Pattern represent the 8-pixel line of a tile
type Pattern struct {
	// Bitmap where each bit represents a pixel
	Bitmap uint8
	// Color in :
	// bits 3-0: Background color
	// bits 7-4: Foreground color
	Color uint8
}

// Tile is formed by 8 patterns, representing an 8x8 pixels tile
type Tile [tilePatterns]Pattern

// Image keeps the data of a whole Screen 2 image
type Image struct {
	Table1 []Tile
	Table2 []Tile
	Table3 []Tile
	Names1 []uint8
	Names2 []uint8
	Names3 []uint8
}
