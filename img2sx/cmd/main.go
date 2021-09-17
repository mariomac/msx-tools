package main

import (
	"bytes"
	"github.com/mariomac/msxtools/img2sx/pkg/img"
	"github.com/mariomac/msxtools/img2sx/pkg/sc2"
	"image/png"
	"io/ioutil"
	"os"
)

func main() {
	raw, err := ioutil.ReadFile("test/input/goku.jpg")
	panicOnErr(err)

	bmp, err := img.Load(bytes.NewReader(raw), sc2.Palette)
	panicOnErr(err)

	sc2img := sc2.FromImage(bmp)

	out, err := os.OpenFile("output.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	panicOnErr(err)
	defer out.Close()
	panicOnErr(png.Encode(out, &sc2img))

	out2, err := os.OpenFile("output.sc2", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	panicOnErr(err)
	defer out2.Close()
	sc2img.Write(out2)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}