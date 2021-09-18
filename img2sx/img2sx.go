package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"

	_ "github.com/mariomac/msxtools/img2sx/pkg/sc2"
)

//type imageConverter

func main() {
	raw, err := ioutil.ReadFile("test/input/burns.sc2")
	panicOnErr(err)

	bmp, format, err := image.Decode(bytes.NewReader(raw))
	panicOnErr(err)
	fmt.Println("decoded format", format)

	out, err := os.OpenFile("output.png", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	panicOnErr(err)
	defer out.Close()
	panicOnErr(png.Encode(out, bmp))

}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
