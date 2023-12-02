package processing

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"io"
	"regexp"
	"slices"
	"time"
)

var dialogPrefix = []byte("Dialogue: ")
var stylePrefix = []byte("Style: ")
var formatPrefix = []byte("Format: ")

var eventsHeader = []byte("[Events]")
var sectionHeaders = [][]byte{
	[]byte("[Script Info]"),
	[]byte("[Aegisub Project Garbage]"),
	[]byte("[Fonts]"),
	[]byte("[V4+ Styles]"),
	eventsHeader,
}

var dialogRegexp = regexp.MustCompile(`^Dialogue: (?P<Layer>[\d.]+),(?P<StartTimecode>[\d:.]+),(?P<Line>.*)`)

func PostProcess(r io.ReadSeeker, w io.WriteSeeker) (err error) {
	fmt.Print("[2pass] Processing starting\n")

	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	_, err = w.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	sections := make(map[string][][]byte)

	//Support lines up to 128 MiB in length
	buf := make([]byte, 1024*1024*128)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(buf, len(buf))

	var pos int64

	var currentSectionName string
	var currentSection [][]byte

	var eventLines EventLineHeaders

	for scanner.Scan() {
		line := scanner.Bytes()

		prevPos := pos
		pos += int64(len(line)) + 1

		if len(line) == 0 {
			// empty line
			continue
		}

		if func() bool {
			for _, h := range sectionHeaders {
				if bytes.HasPrefix(line, h) {
					if currentSectionName != "" {
						sections[currentSectionName] = currentSection
					}
					currentSectionName = string(h)
					//allow re-entering a section
					currentSection = sections[currentSectionName]
					return true
				}
			}
			return false
		}() {
			//It's a header entry
			continue
		}

		if currentSectionName == "" {
			return errors.New("expected section header first")
		}

		matches := dialogRegexp.FindSubmatch(line)
		if matches == nil {
			// not a line, write directly to section
			currentSection = append(currentSection, slices.Clone(line))
			continue
		}

		layer := matches[1]
		startTimecode := matches[2]
		restOfLine := matches[3]

		var hours, minutes, seconds, milliseconds int32

		res, _ := fmt.Sscanf(string(startTimecode), "%d:%d:%d.%d", &hours, &minutes, &seconds, &milliseconds)
		if res < 4 {
			return errors.New("timecode not parsed properly")
		}

		depth, err := types.DepthFromString(string(layer))
		if err != nil {
			return err
		}

		start := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second + time.Duration(milliseconds)*10*time.Millisecond

		eventLines = append(eventLines, EventLineHeader{
			ReadOrder:       len(eventLines),
			Index:           prevPos,
			IndexFromLayer:  prevPos + int64(len(dialogPrefix)+len(layer)),
			Length:          len(line),
			LengthFromLayer: len(line) - (len(dialogPrefix) + len(layer)),
			Depth:           depth,
			Start:           start,
		})

		_ = restOfLine
	}

	if currentSectionName != "" {
		sections[currentSectionName] = currentSection
	}

	fmt.Printf("[2pass] Processed %d bytes, %d event lines\n", pos, len(eventLines))

	eventLines.Sort()

	fmt.Printf("[2pass] Sorted %d event lines, writing back\n", len(eventLines))

	for _, k := range sectionHeaders {
		section := sections[string(k)]
		if len(section) == 0 {
			continue
		}

		_, err = w.Write(k)
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{0x0a})
		if err != nil {
			return err
		}
		for _, l := range section {
			_, err = w.Write(l)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte{0x0a})
			if err != nil {
				return err
			}
		}
		if !bytes.Equal(k, eventsHeader) {
			_, err = w.Write([]byte{0x0a, 0x0a})
			if err != nil {
				return err
			}
		}
	}

	var lineBuffer []byte

	for _, eventLine := range eventLines {
		if len(lineBuffer) < eventLine.Length {
			lineBuffer = make([]byte, eventLine.Length+1024)
		}

		_, err = r.Seek(eventLine.IndexFromLayer, io.SeekStart)
		if err != nil {
			return err
		}
		_, err = io.ReadFull(r, lineBuffer[:eventLine.LengthFromLayer])
		if err != nil {
			return err
		}
		_, err = w.Write([]byte("Dialogue: 0"))
		if err != nil {
			return err
		}
		_, err = w.Write(lineBuffer[:eventLine.LengthFromLayer])
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{0x0a})
		if err != nil {
			return err
		}
	}
	return nil
}
