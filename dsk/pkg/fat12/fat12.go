package fat12

/*
----------------------------------------------------------------------
|                            | 1DD, 9  | 2DD, 9  | 1DD, 8  | 2DD, 8  |
|                            | sectors | sectors | sectors | sectors |
|----------------------------+---------+---------+---------+---------|
| media ID                   |  0F8H   |  0F9H   |  0FAH   |  0FBH   |
| number of sides            |    1    |    2    |    1    |    2    |
| tracks per side            |   80    |   80    |   80    |   80    |
| sectors per track          |    9    |    9    |    8    |    8    |
| bytes per sector           |   512   |   512   |   512   |   512   |
| cluster size (in sectors)  |    2    |    2    |    2    |    2    |
| FAT size (in sectors)      |    2    |    3    |    1    |    2    |
| number of FATs             |    2    |    2    |    2    |    2    |
| number of recordable files |   112   |   112   |   112   |   112   |
----------------------------------------------------------------------
*/

// Creates a 9-sectors 2DD disk. According to the MSX technical handbook:
// https://konamiman.github.io/MSX2-Technical-Handbook/md/Chapter3.html#table-31--media-supported-by-msx-dos
const (
	mediaID             = 0xf9
	numberOfSides       = 2
	tracksPerSide       = 80
	sectorsPerTrack     = 9
	bytesPerSector      = 512
	sectorsPerCluster   = 2
	fatSizeSectors      = 3
	numberOfFats        = 2
	maxRootDirEntries   = 112
	sectorsPerFat       = 1
	mediaDescriptorByte = 0xf9
	hiddenSectors       = 0
)

var bootSectorHeader = []byte{0xeb, 0xfb, 0x90, 'M', 'a', 'c', 'i', 'a', 's', '!', '!'}
var bootSectorSignature = []byte{0x55, 0xaa}

type Disk []byte

type BootSector []byte
type RootDirectory struct{}
type DataArea struct{}

func NewDisk() Disk {
	//bs := newBootSector() // 1 sector
	//fat := newFatTable()  // to be added twice. 2 sectors per 2 fat
	// fat 2
	// root dirRegion
	// data
	return Disk{}
}

func w2bs(n int16) []byte {
	return []byte{byte(n & 7), byte(n >> 8)}
}

func newBootSector() []byte {
	bs := make([]byte, 0, bytesPerSector)
	bs = append(bs, bootSectorHeader...)
	bs = append(bs, w2bs(bytesPerSector)...)
	bs = append(bs, sectorsPerCluster)
	bs = append(bs, w2bs(sectorsPerFat)...)
	bs = append(bs, numberOfFats)
	bs = append(bs, w2bs(maxRootDirEntries)...)
	// total number of sectors in the filesystem
	bs = append(bs, w2bs(numberOfSides*tracksPerSide*sectorsPerTrack)...)
	// media descriptor byte from https://www.win.tue.nl/~aeb/linux/fs/fat/fat-1.html
	bs = append(bs, mediaDescriptorByte)
	bs = append(bs, w2bs(sectorsPerFat)...)
	bs = append(bs, w2bs(sectorsPerTrack)...)
	bs = append(bs, w2bs(numberOfSides)...)
	bs = append(bs, w2bs(hiddenSectors)...)
	// bootstrap zeroes
	bs = append(bs, make([]byte, bytesPerSector-len(bs)-len(bootSectorSignature))...)
	return append(bs, bootSectorSignature...)
}
