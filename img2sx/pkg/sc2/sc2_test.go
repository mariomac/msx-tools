package sc2

import (
	"bytes"
	"github.com/mariomac/msxtools/img2sx/pkg/img"
	"image"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

var t1 = Tile{
	Pattern{Bitmap: 0, Color: 7},
	Pattern{Bitmap: 1, Color: 6},
	Pattern{Bitmap: 2, Color: 5},
	Pattern{Bitmap: 3, Color: 4},
	Pattern{Bitmap: 4, Color: 3},
	Pattern{Bitmap: 5, Color: 2},
	Pattern{Bitmap: 6, Color: 1},
	Pattern{Bitmap: 7, Color: 0},
}
var t2 = Tile{
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
	Pattern{Bitmap: 2, Color: 2},
}

func TestEquality(t *testing.T) {
	t1bis := Tile{
		Pattern{Bitmap: 0, Color: 7},
		Pattern{Bitmap: 1, Color: 6},
		Pattern{Bitmap: 2, Color: 5},
		Pattern{Bitmap: 3, Color: 4},
		Pattern{Bitmap: 4, Color: 3},
		Pattern{Bitmap: 5, Color: 2},
		Pattern{Bitmap: 6, Color: 1},
		Pattern{Bitmap: 7, Color: 0},
	}
	assert.Equal(t, t1, t1bis)
	assert.NotEqual(t, t1, t2)
}

func TestWriteImage(t *testing.T) {
	sc := TileSet{
		Table: [3][]Tile {
			{t1, t2, t2},
			{t2, t1, t2},
			{t2, t2, t1},
		},
		Names: [3][]uint8 {
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}
	buf := bytes.Buffer{}
	require.NoError(t, sc.Write(&buf))
	out := buf.Bytes()
	assert.Len(t, out, spriteGen+len(signature))
	// the signature has been copied at the beginning
	assert.Equal(t, signature, out[:len(signature)])
	out = out[len(signature):]
	// the table 1 tiles have been copied and the rest is filled with 0
	assert.Equal(t, []uint8{
		0, 1, 2, 3, 4, 5, 6, 7,
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	// table 2 tiles have been copied and the rest is filled with 0
	out = out[tablePatterns:]
	assert.Equal(t, []uint8{
		2, 2, 2, 2, 2, 2, 2, 2,
		0, 1, 2, 3, 4, 5, 6, 7,
		2, 2, 2, 2, 2, 2, 2, 2,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	// table 3 tiles have been copied and the rest is filled with 0
	out = out[tablePatterns:]
	assert.Equal(t, []uint8{
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
		0, 1, 2, 3, 4, 5, 6, 7,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	out = out[tablePatterns:]
	// Pattern name tables are written and the rest is filled with 0
	assert.Equal(t, append([]uint8{1, 2, 3}, make([]uint8, tableTiles-3)...), out[:tableTiles], "table 1 not correct")
	out = out[tableTiles:]
	assert.Equal(t, append([]uint8{4, 5, 6}, make([]uint8, tableTiles-3)...), out[:tableTiles], "table 2 not correct")
	out = out[tableTiles:]
	assert.Equal(t, append([]uint8{7, 8, 9}, make([]uint8, tableTiles-3)...), out[:tableTiles], "table 3 not correct")
	out = out[tableTiles:]
	// sprite and palette tables are zeroes
	assert.Equal(t, make([]uint8, color1-spriteAttrs), out[:color1-spriteAttrs])
	out = out[color1-spriteAttrs:]
	// Color tables are copied and the rest are zeroes
	assert.Equal(t, []uint8{
		7, 6, 5, 4, 3, 2, 1, 0,
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	out = out[tablePatterns:]
	assert.Equal(t, []uint8{
		2, 2, 2, 2, 2, 2, 2, 2,
		7, 6, 5, 4, 3, 2, 1, 0,
		2, 2, 2, 2, 2, 2, 2, 2,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	out = out[tablePatterns:]
	assert.Equal(t, []uint8{
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
		7, 6, 5, 4, 3, 2, 1, 0,
	}, out[:24])
	assert.Equal(t, make([]uint8, tablePatterns-24), out[24:tablePatterns])
	out = out[tablePatterns:]

	// and there isn't anything else left in the output file
	assert.Empty(t, out)
}

func TestSamplePattern(t *testing.T) {
	img := img.Bitmap{Img: image.NewNRGBA(image.Rect(0, 0, 8, 1))}
	img.Img.Set(0, 0, Palette[4])
	img.Img.Set(1, 0, Palette[4])
	img.Img.Set(2, 0, Palette[7])
	img.Img.Set(3, 0, Palette[9])
	img.Img.Set(4, 0, Palette[9])
	img.Img.Set(5, 0, Palette[9])
	img.Img.Set(6, 0, Palette[9])
	img.Img.Set(7, 0, Palette[4])

	p := SamplePattern(img, 0, 0)
	assert.EqualValuesf(t, 0b00011110, p.Bitmap, "%08b", p.Bitmap)
	assert.EqualValuesf(t, 0b1001_0100, p.Color, "%08b", p.Color)
}

func TestSamplePattern_OneColor(t *testing.T) {
	img := img.Bitmap{Img: image.NewNRGBA(image.Rect(0, 0, 8, 1))}
	img.Img.Set(0, 0, Palette[4])
	img.Img.Set(1, 0, Palette[4])
	img.Img.Set(2, 0, Palette[4])
	img.Img.Set(3, 0, Palette[4])
	img.Img.Set(4, 0, Palette[4])
	img.Img.Set(5, 0, Palette[4])
	img.Img.Set(6, 0, Palette[4])
	img.Img.Set(7, 0, Palette[4])

	p := SamplePattern(img, 0, 0)
	assert.EqualValuesf(t, 0b11111111, p.Bitmap, "%08b", p.Bitmap)
	assert.EqualValuesf(t, 0b0100_1111, p.Color, "%08b", p.Color)
}

// TODO: test sc2 -> img -> sc2 conversion and see that input and output are equal
