package fat12

const maxClusters = (numberOfSides * tracksPerSide * sectorsPerTrack) / sectorsPerCluster

type fatTable []byte

// 12 bits for fat-12, encoded as 16-bits
type fatEntry uint16

const freeClusterEntry fatEntry = 0
const badClusterEntry fatEntry = 0xFF7
const lastClusterEntry fatEntry = 0xFFF

func newFatTable() fatTable {
	// 1 sector per fat
	ft := make([]byte, sectorsPerFat*bytesPerSector)
	return ft
}

// gets the nth entry of the fat table as a 12-bits uint
func (ft fatTable) getUint12(n int) uint16 {
	i := n + (n / 2)
	if n%2 == 0 {
		return (uint16(ft[i+1]&0xF) << 8) |
			uint16(ft[i])
	} else {
		return (uint16(ft[i]&0xF0) >> 4) |
			(uint16(ft[i+1]) << 4)
	}
}

func (ft fatTable) setUint12(n int, val uint16) {
	i := n + (n / 2)
	if n%2 == 0 {
		ft[i] = uint8(val)
		ft[i+1] = ft[i+1]&0xf0 | uint8(val>>8)&0xF
	} else {
		ft[i] = ft[i]&0xF | uint8(val&0xF)<<4
		ft[i+1] = uint8(val >> 4)
	}
}
