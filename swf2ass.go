package main

import (
	"errors"
	"flag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf"
	"io"
	"os"
)

func main() {
	inputFile := flag.String("input", "", "Input SWF")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	swfReader, err := swf.NewReader(file)
	if err != nil {
		panic(err)
	}

	for {
		tag, err := swfReader.Tag()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
		_ = tag
	}
}
