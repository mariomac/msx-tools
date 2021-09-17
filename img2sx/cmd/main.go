package main

import (
	"bytes"
	"github.com/mariomac/msxtools/img2sx/pkg/img"
	"github.com/mariomac/msxtools/img2sx/pkg/sc2"
	"io/ioutil"
	"os"
)

func main() {
	raw, err := ioutil.ReadFile("test/input/goku.jpg")
	panicOnErr(err)

	bmp, err := img.Load(bytes.NewReader(raw), sc2.Palette)
	panicOnErr(err)

	out, err := os.OpenFile("output.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	panicOnErr(err)

	panicOnErr(bmp.Save(out, img.FormatPNG))
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}