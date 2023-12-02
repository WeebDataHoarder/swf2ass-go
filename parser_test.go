package main

import (
	"errors"
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf-go"
	swftag "git.gammaspectra.live/WeebDataHoarder/swf-go/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"io"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	file, err := os.Open("im_alive.swf")
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
			if errors.Is(err, swftag.ErrUnknownTag) {
				continue
			}
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

		//TODO: handle multiple sounds

		if processor.Audio != nil && frameOffset == 0 {
			if processor.Audio.Start == nil {
				fmt.Printf("Skipped frame %d: audio not started\n", frame.FrameNumber)
				continue
			}
			frameOffset = *processor.Audio.Start
		} else if processor.Audio == nil {
			//TODO: make this an option
			fmt.Printf("Skipped frame %d: no audio\n", frame.FrameNumber)
			continue
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

		assRenderer.RenderFrame(*frame, filteredRendered, settings.GlobalSettings.KeyFrameInterval)
	}

	if lastFrame == nil {
		t.Fatal("no frames generated")
	}

	assRenderer.Flush(*lastFrame)
}
