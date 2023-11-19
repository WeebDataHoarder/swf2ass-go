package main

import (
	"errors"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"io"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.Open("azumanga_vector.swf")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	swfReader, err := swf.NewReader(file)
	if err != nil {
		t.Fatal(err)
	}

	for {
		readTag, err := swfReader.Tag()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatal(err)
		}

		if readTag == nil {
			//not decoded
			continue
		}

		if readTag.Code() == tag.RecordEnd {
			break
		}

		switch t := readTag.(type) {

		default:
			_ = t
		}
	}
}
