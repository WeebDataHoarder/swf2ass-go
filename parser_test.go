package main

import (
	"errors"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf"
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
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

	processor := types.NewSWFProcessor(tags, shapes.RectangleFromSWF(swfReader.Header().FrameSize), swfReader.Header().FrameRate.Float64(), int64(swfReader.Header().FrameCount), swfReader.Header().Version)

	assRenderer := ass.NewRenderer(processor.FrameRate, processor.ViewPort)

	const KeyFrameEveryNSeconds = 10

	keyframeInterval := int64(KeyFrameEveryNSeconds * processor.FrameRate)

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

		rendered := frame.Frame.Render(0, nil, types.None[math.ColorTransform](), types.None[math.MatrixTransform]())

		if frame.GetFrameNumber() == 0 {
			for _, object := range rendered {
				t.Logf("frame 0: object %d depth: %s\n", object.ObjectId, object.Depth.String())
			}
		}

		filteredRendered := make(types.RenderedFrame, 0, len(rendered))

		var drawCalls, drawItems, filteredObjects, clipCalls, clipItems int

		for _, object := range rendered {
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

		t.Logf("=== frame %d/%d ~ %d : Depth count: %d :: Object count: %d :: Paths: %d draw calls, %d items :: Filtered: %d :: Clips %d draw calls, %d items\n",
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

		assRenderer.RenderFrame(*frame, filteredRendered)

		//TODO: do this per object transition? GlobalSettings?
		if frame.GetFrameNumber() > 0 && frame.GetFrameNumber()%keyframeInterval == 0 {
			assRenderer.Flush(*frame)
		}
	}

	if lastFrame == nil {
		t.Fatal("no frames generated")
	}

	assRenderer.Flush(*lastFrame)
}
