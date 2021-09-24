package fat12

import "time"

// from: https://konamiman.github.io/MSX2-Technical-Handbook/md/Chapter3.html#figure-314--directory-construction
type FileAttrs struct {
	NameExtension [11]uint8
	FileAttribute uint8
	SpaceNotUsed [10]uint8  // todo: remove
	TimeCreation uint16
	DateCreation uint16
	TopCluster uint16
	FileSize uint32
}

func encodeTime(t time.Time) (date, time uint16) {
	time = uint16((t.Second() / 2) & 0b11111 |
		(t.Minute() & 0b111111) << 5 |
		(t.Hour() & 0b11111) << 11)
	date = uint16(t.Day() & 0b11111 |
		(int(t.Month()) & 0b1111) << 5 |
		((t.Year() - 1980) & 0b1111111) << 9)
	return
}
