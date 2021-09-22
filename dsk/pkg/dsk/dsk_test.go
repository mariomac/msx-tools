package dsk

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeDate(t *testing.T) {
	dat, tim := encodeTime(time.Date(2021, 9, 22, 20, 9, 33, 0, time.UTC))
	assert.EqualValues(t, 0b_10100_001001_10000, tim)
	assert.EqualValues(t, 0b_101001_1001_10110, dat)
}