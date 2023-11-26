package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf"
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"io"
	math2 "math"
	"os"
	"path"
)

type KnownSignatures map[string]KnownSignature

type KnownSignature struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Remove      []RemovalEntry `json:"remove"`
}

func (s KnownSignature) Filter(object *types.RenderedObject) bool {
	for i := range s.Remove {
		if s.Remove[i].Equals(object) {
			return true
		}
	}
	return false
}

type RemovalEntry struct {
	ObjectId        *uint16     `json:"objectId"`
	ObjectIdComment *uint16     `json:"_objectId"`
	Depth           types.Depth `json:"depth"`
}

func (e RemovalEntry) Equals(object *types.RenderedObject) bool {
	return (e.ObjectId == nil || *e.ObjectId == object.ObjectId) && (len(e.Depth) == 0 || (len(object.Depth) >= len(e.Depth)) && object.Depth[:len(e.Depth)].Equals(e.Depth))
}

func main() {
	inputFile := flag.String("input", "", "Input SWF")
	outputFile := flag.String("output", "", "Output ASS")
	outputAudio := flag.String("audio", "", "Output Audio")
	removalSignatures := flag.String("signatures", "signatures.json", "JSON file containing parameters for signature removal")
	fromFrame := flag.Int64("from", 0, "Frame to start at")
	toFrame := flag.Int64("to", math2.MaxInt64, "Frame to end at")
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

	var knownSignatures KnownSignatures
	removalSignaturesData, err := os.ReadFile(*removalSignatures)
	if err == nil {
		_ = json.Unmarshal(removalSignaturesData, &knownSignatures)
	}

	var tags []swftag.Tag

	for {
		readTag, err := swfReader.Tag()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}

		if readTag == nil {
			//not decoded
			continue
		}

		tags = append(tags, readTag)

		if readTag.Code() == swftag.RecordEnd {
			break
		}

		switch t := readTag.(type) {

		default:
			_ = t
		}
	}

	var frameOffset int64

	processor := types.NewSWFProcessor(tags, shapes.RectangleFromSWF(swfReader.Header().FrameSize), swfReader.Header().FrameRate.Float64(), int64(swfReader.Header().FrameCount))

	assRenderer := ass.NewRenderer(processor.FrameRate, processor.ViewPort)

	const KeyFrameEveryNSeconds = 10

	keyframeInterval := int64(-1) //int64(KeyFrameEveryNSeconds * processor.FrameRate)

	var ks KnownSignature

	for _, s := range knownSignatures {
		if s.Name == path.Base(*inputFile) {
			ks = s
			break
		}
	}

	output, err := os.Create(*outputFile)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	outputLines := func(lines ...string) {
		for _, line := range lines {
			_, err = output.Write([]byte(line))
			if err != nil {
				panic(err)
			}
			_, err = output.Write([]byte("\n"))
			if err != nil {
				panic(err)
			}
		}
	}

	var lastFrame *types.FrameInformation
	for {
		frame := processor.NextFrameOutput()
		if frame == nil {
			break
		}
		lastFrame = frame
		if !processor.Playing || processor.Loops > 0 {
			break
		}

		if processor.Audio != nil && frameOffset == 0 {
			if processor.Audio.Start == nil {
				continue
			}
			frameOffset = *processor.Audio.Start
		}

		frame.FrameOffset = frameOffset

		rendered := frame.Frame.Render(0, nil, nil, nil)

		if frame.GetFrameNumber() == 0 {
			for _, object := range rendered {
				fmt.Printf("frame 0: object %d depth: %s\n", object.ObjectId, object.Depth.String())
			}
		}

		filteredRendered := make(types.RenderedFrame, 0, len(rendered))

		var drawCalls, drawItems, filteredObjects, clipCalls, clipItems int

		for _, object := range rendered {
			if ks.Filter(object) {
				filteredObjects++
				continue
			}
			if object.Clip != nil {
				clipCalls++
				clipItems += len(object.Clip.GetShape())
			}
			for _, p := range object.DrawPathList {
				drawCalls++
				drawItems += len(p.Shape)
			}
			filteredRendered = append(filteredRendered, object)
		}

		fmt.Printf("=== frame %d/%d ~ %d : Depth count: %d :: Object count: %d :: Paths: %d draw calls, %d items :: Filtered: %d :: Clips %d draw calls, %d items\n",
			frame.GetFrameNumber(),
			processor.ExpectedFrameCount,
			frameOffset,
			len(frame.Frame.DepthMap),
			len(filteredRendered),
			drawCalls,
			drawItems,
			filteredObjects,
			clipCalls,
			clipItems,
		)

		if *fromFrame > 0 {
			if frame.GetFrameNumber() < *fromFrame {
				continue
			} /*else {
				for _, object := range rendered {
					var count int

					for i, command := range object.DrawPathList {
						for j, record := range command.Shape.Edges {

						}
					}
				}
			}*/
		}

		outputLines(assRenderer.RenderFrame(*frame, filteredRendered)...)

		//TODO: do this per object transition? GlobalSettings?
		if frame.GetFrameNumber() > 0 && keyframeInterval != -1 && frame.GetFrameNumber()%keyframeInterval == 0 {
			outputLines(assRenderer.Flush(*frame)...)
		}

		if *toFrame != math2.MaxInt64 && frame.GetFrameNumber() >= *toFrame {
			break
		}
	}

	if lastFrame == nil {
		panic("no frames generated")
	}

	outputLines(assRenderer.Flush(*lastFrame)...)

	if *outputAudio != "" && processor.Audio != nil && processor.Audio.Format == swftag.SoundFormatMP3 {
		_ = os.WriteFile(*outputAudio, processor.Audio.Data, 0664)
	}

	stats := assRenderer.AggregateStatistics()
	stats.SortBySize()
	for _, s := range stats.Strings() {
		print(s + "\n")
	}

}
