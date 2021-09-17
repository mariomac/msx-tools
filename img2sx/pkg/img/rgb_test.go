package img

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRGBDistance(t *testing.T) {
	type testCase struct {
		c1     RGB
		c2     RGB
		expect int
	}
	cases := []testCase{
		{c1: RGB(0x010101), c2: RGB(0x000000), expect: 3},
		{c1: RGB(0x010100), c2: RGB(0x000001), expect: 3},
		{c1: RGB(0xFF0000), c2: RGB(0x0000FF), expect: 0xFF * 2},
		{c1: RGB(0x0011FF), c2: RGB(0xFF0000), expect: 0xFF*2 + 0x11},
		{c1: RGB(0x123456), c2: RGB(0x654321), expect: 83 + 15 + 53},
		{c1: RGB(0xFFFFFF), c2: RGB(0xFFFFFF), expect: 0},
		{c1: RGB(0), c2: RGB(0), expect: 0},
	}
	for _, tc := range cases {
		assert.Equalf(t, tc.expect, tc.c1.DistanceTo(tc.c2), "%v - %v", tc.c1, tc.c2)
		assert.Equalf(t, tc.expect, tc.c2.DistanceTo(tc.c1), "%v - %v", tc.c2, tc.c1)
	}
}
