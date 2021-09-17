package sc2

import "io"

// first 7 bytes of an SC2 file
var signature = []uint8{0xfe, 0x00, 0x00, 0xff, 0x37, 0x00, 0x00}

func (s *TileSet) Write(out io.Writer) error {
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
