package img

import (
	"fmt"
	"image/color"
)

// RGB color without Alpha (24 lower bits)
type RGB uint32

var _ color.Color = RGB(0)

// Red component of a color
func (c RGB) Red() uint8 {
	return uint8(c >> 16)
}

// Green component of a color
func (c RGB) Green() uint8 {
	return uint8(c >> 8)
}

// Blue component of a color
func (c RGB) Blue() uint8 {
	return uint8(c)
}

// DistanceTo another color. It calculates Manhattan distance
func (c *RGB) DistanceTo(o RGB) int {
	// efficient way to calculate an absolute value between 2 numbers
	// https://stackoverflow.com/questions/664852/which-is-the-fastest-way-to-get-the-absolute-value-of-a-number
	var abs = func(i int) int {
		temp := uint16(i >> 15)                     // make a mask of the sign bit
		return int((uint16(i) ^ temp) + (temp & 1)) // toggle the bits if value is negative
	}
	return abs(int(c.Red())-int(o.Red())) +
		abs(int(c.Green())-int(o.Green())) +
		abs(int(c.Blue())-int(o.Blue()))
}

// String representation in the form of hexadecimal #RRGGBB
func (c RGB) String() string {
	return fmt.Sprintf("#%06X", uint32(c))
}

// RGBA implements color.Color interface
// TODO: inline the same formula as color.NRGBA item
func (c RGB) RGBA() (r, g, b, a uint32) {
	return color.NRGBA{R: c.Red(), G: c.Green(), B: c.Blue(), A: 0xFF}.RGBA()
}
