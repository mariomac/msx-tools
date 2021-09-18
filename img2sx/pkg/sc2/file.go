package sc2

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
)

var signature = []uint8{0xfe, 0x00, 0x00, 0xff, 0x37, 0x00, 0x00}

func init() {
	image.RegisterFormat("sc2", string(signature), decode, decodeConfig)
}

const minFileLength = 14343

func decode(r io.Reader) (image.Image, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if len(content) > minFileLength {
		return nil, fmt.Errorf("file is too short (%v bytes, min length: %v)",
			len(content), minFileLength)
	}
	if bytes.Compare(signature, content[:len(signature)]) != 0 {
		return nil, fmt.Errorf("invalid signature: %v", content[:len(signature)])
	}
	content = content[len(signature):]
	// read pattern generator tables
	ts := TileSet{}
	i := 0
	for t := range ts.Table {
		for tileNum := 0 ; tileNum < 256 ; tileNum++ {
			tl := Tile{}
			for p := 0; p < 8; p++ {
				tl[p].Bitmap = content[i]
				i++
			}
			ts.Table[t] = append(ts.Table[t], tl)
		}
	}
	content = content[i:]
	// read pattern name tables
	ts.Names[0] = content[:256]
	ts.Names[1] = content[256:512]
	ts.Names[2] = content[512:768]
	content = content[768:]
	// discard sprites and palette
	content = content[color1-spriteAttrs:]

	// read color tables
	i = 0
	for t := range ts.Table {
		for tl := range ts.Table[t] {
			for p := range ts.Table[t][tl] {
				ts.Table[t][tl][p].Color = content[i]
				i++
			}
		}
	}
	return &ts, nil
}

func decodeConfig(_ io.Reader) (image.Config, error) {
	return image.Config{
		ColorModel: Palette,
		Width: 256,
		Height: 192,
	}, nil
}

func Encode(out io.Writer, i image.Image) error {
	s, ok := i.(*TileSet)
	if !ok {
		s = FromImage(i)
	}

	// Write file signature
	if _, err := out.Write(signature); err != nil {
		return err
	}
	// Write pattern generator tables
	bget := func(p Pattern) uint8 {
		return p.Bitmap
	}
	for _, table := range s.Table {
		if _, err := out.Write(acquireBytes(table, true, bget)); err != nil {
			return err
		}
	}
	// Write pattern name tables
	for _, names := range s.Names {
		if err := writeNames(out, names, true); err != nil {
			return err
		}
	}
	// Fill sprite attributes & palette with zeroes
	if _, err := out.Write(make([]uint8, color1-spriteAttrs)); err != nil {
		return err
	}
	// Fill color tables
	cget := func(p Pattern) uint8 {
		return p.Color
	}
	for _, table := range s.Table {
		if _, err := out.Write(acquireBytes(table, true, cget)); err != nil {
			return err
		}
	}
	return nil
}

// if fillZeroes == true, it fills the array with zeros until
// its size
func acquireBytes(tiles []Tile, fillZeroes bool, getter func(Pattern) uint8) []uint8 {
	bytes := make([]uint8, tablePatterns)
	if len(tiles) > tableTiles {
		tiles = tiles[:tableTiles]
	}
	for tn, tile := range tiles {
		for pn, pattern := range tile {
			bytes[tn*tilePatterns+pn] = getter(pattern)
		}
	}
	if fillZeroes {
		return bytes
	} else {
		return bytes[:len(tiles)*tilePatterns]
	}
}

func writeNames(out io.Writer, names []uint8, fillZeroes bool) error {
	if _, err := out.Write(names); err != nil {
		return err
	}
	if fillZeroes && len(names) < tableTiles {
		if _, err := out.Write(make([]uint8, tableTiles-len(names))); err != nil {
			return err
		}
	}
	return nil
}
