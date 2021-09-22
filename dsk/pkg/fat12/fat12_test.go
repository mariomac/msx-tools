package fat12

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBootSector(t *testing.T) {
	bs := newBootSector()
	assert.Len(t, bs, 512)
	assert.EqualValues(t, w2bs(numberOfSides), bs[26:28])
}
