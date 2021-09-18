package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mariomac/msxtools/img2sx/pkg/sc2"
)

func main() {
	inFile, outFile, options := getConfig()

	outFileType := path.Ext(outFile)

	raw, err := ioutil.ReadFile(inFile)
	exitOnErr(err)

	img, format, err := image.Decode(bytes.NewReader(raw))
	exitOnErr(err)
	fmt.Println("decoded format", format)

	out, err := os.OpenFile(outFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	exitOnErr(err)
	defer out.Close()

	switch strings.ToLower(outFileType) {
	case ".sc2":
		e := sc2.Encoder{Opt: sc2.ConvertOpt(options)}
		exitOnErr(e.Encode(out, img))
	case ".png":
		exitOnErr(png.Encode(out, img))
	case ".jpg", ".jpeg":
		// todo: allow jpeg options
		exitOnErr(jpeg.Encode(out, img, nil))
	case ".gif":
		// todo: allow gif options
		exitOnErr(gif.Encode(out, img, nil))
	default:
		exitOnErr(fmt.Errorf("unknown output format: %s", outFileType))
	}
	fmt.Println("Done!")
}

func getConfig() (inFile, outFile, options string) {
	flag.Usage = func() {
		fmt.Println("Usage of", os.Args[0])
		flag.PrintDefaults()
	}
	out := flag.String("out", "", "Output file")
	in := flag.String("in", "", "Input file")
	opt := flag.String("o", "", "Options. For SC2: crop, stretch, keepaspect (default)")
	help := flag.Bool("h", false, "Print this help")
	flag.Parse()
	if help != nil && *help {
		flag.Usage()
		os.Exit(0)
	}
	errStr := ""
	if out == nil || *out == "" {
		errStr = "missing output file"
	} else if in == nil || *in == "" {
		errStr = "missing input file"
	} else if opt != nil && *opt != "" {
		validOptions := map[string]struct{}{
			"crop":       {},
			"stretch":    {},
			"keepaspect": {},
		}
		if _, ok := validOptions[*opt]; !ok {
			errStr = "invalid option: " + *opt
		}
	}
	if errStr != "" {
		fmt.Fprintf(os.Stderr, "ERROR: %s", errStr)
		os.Exit(-1)
	}
	return *in, *out, *opt
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(-1)
	}
}
