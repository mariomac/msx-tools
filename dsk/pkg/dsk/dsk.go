package dsk

import "time"

const allowedFile = "$%'-_@~`!(){}^#&"
type fileName string

func (fn *fileName) normalize() [11]byte {
	// extension
	return [11]byte{}
}

type File struct {
	name fileName
	creation time.Time
	lastAccess time.Time
	lastWrite time.Time
	data []byte
}

type Disk struct {
	//key: filename
	files map[fileName]File
}