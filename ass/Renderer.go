package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/line"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"runtime"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Renderer struct {
	Header        []string
	RunningBuffer []*line.EventLine
	Display       shapes.Rectangle[float64]

	Statistics map[uint16]rendererStatsEntry
}

type rendererStatsEntry struct {
	Frames *atomic.Uint64
	Lines  *atomic.Uint64
	Size   *atomic.Uint64
}

func NewRenderer(frameRate float64, display shapes.Rectangle[float64]) *Renderer {
	width := int64(display.Width() * settings.GlobalSettings.VideoScaleMultiplier)
	height := int64(display.Height() * settings.GlobalSettings.VideoScaleMultiplier)

	ar := float64(width) / float64(height)

	frameRate *= settings.GlobalSettings.VideoRateMultiplier

	return &Renderer{
		Statistics: make(map[uint16]rendererStatsEntry),
		Display:    display,
		Header: []string{
			"[Script Info]",
			"; Script generated by swf2ass Renderer",
			"; https://git.gammaspectra.live/WeebDataHoarder/swf2ass-go",
			"Title: swf2ass",
			"ScriptType: v4.00+",
			//TODO: WrapStyle: 0 or 2?
			"WrapStyle: 2",
			"ScaledBorderAndShadow: yes",
			"YCbCr Matrix: PC.709",
			fmt.Sprintf("PlayResX: %d", width),
			fmt.Sprintf("PlayResY: %d", height),
			"",
			"",
			"[Aegisub Project Garbage]",
			"Last Style Storage: f",
			fmt.Sprintf("Video File: ?dummy:%s:10000:%d:%d:160:160:160:c", strconv.FormatFloat(frameRate, 'f', -1, 64), width, height),
			fmt.Sprintf("Video AR Value: %.6F", ar),
			"Active Line: 0",
			"Video Zoom Percent: 2.000000",
			"",
			"[V4+ Styles]",
			"Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding",
			fmt.Sprintf("Style: %s,Arial,20,&H00000000,&H00000000,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,0,0,7,0,0,0,1", line.StyleFill),
			fmt.Sprintf("Style: %s,Arial,20,&H00000000,&H00000000,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,0,0,7,0,0,0,1", line.StyleLine),
			"",
			"[Events]",
			"Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text",
		},
	}
}

type RendererStatistics []RendererStatistic

func (s RendererStatistics) SortByFrames() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Frames - a.Frames)
	})
}

func (s RendererStatistics) SortByLines() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Lines - a.Lines)
	})
}

func (s RendererStatistics) SortBySize() {
	slices.SortStableFunc(s, func(a, b RendererStatistic) int {
		return int(b.Size - a.Size)
	})
}

func (s RendererStatistics) Reverse() {
	slices.Reverse(s)
}

func (s RendererStatistics) Strings() (r []string) {
	maxObjectId := len("Object Id")
	maxFrames := len("Frames")
	maxLines := len("Lines")
	maxSize := len("Lines")
	for _, e := range s {
		oId, frames, lines, size := e.Strings()

		maxObjectId = max(maxObjectId, len(oId))
		maxFrames = max(maxFrames, len(frames))
		maxLines = max(maxLines, len(lines))
		maxSize = max(maxSize, len(size))
	}

	maxIndex := len(strconv.FormatUint(uint64(len(s)), 10))

	for i, e := range s {
		oId, frames, lines, size := e.Strings()
		i := strconv.FormatUint(uint64(i), 10)

		r = append(r, fmt.Sprintf("| %*s | %*s | %*s | %*s | %*s |", maxIndex, i, maxObjectId, oId, maxFrames, frames, maxLines, lines, maxSize, size))
	}
	slices.Insert(r, 0, fmt.Sprintf("| %*s | %*s | %*s | %*s | %*s |", maxIndex, "#", maxObjectId, "Object Id", maxFrames, "Frames", maxLines, "Lines", maxSize, "Size"))
	return r
}

type RendererStatistic struct {
	ObjectId uint16
	Frames   uint64
	Lines    uint64
	Size     uint64
}

func (s RendererStatistic) Strings() (string, string, string, string) {
	oId := strconv.FormatUint(uint64(s.ObjectId), 10)
	frames := strconv.FormatUint(s.Frames, 10)
	lines := strconv.FormatUint(s.Lines, 10)
	size := strconv.FormatFloat(float64(s.Size)/(1<<20), 'f', 2, 64) + " MiB"
	return oId, frames, lines, size
}

func (r *Renderer) AggregateStatistics() (s RendererStatistics) {
	for objectId, stats := range r.Statistics {
		s = append(s, RendererStatistic{
			ObjectId: objectId,
			Frames:   stats.Frames.Load(),
			Lines:    stats.Lines.Load(),
			Size:     stats.Size.Load(),
		})
	}
	return s
}

func (r *Renderer) RenderFrame(frameInfo types.FrameInformation, frame types.RenderedFrame) (result []string) {
	if len(r.Header) != 0 {
		result = append(result, r.Header...)
		r.Header = nil
	}

	slices.SortStableFunc(frame, RenderedObjectDepthSort)

	var runningBuffer []*line.EventLine

	scale := math.ScaleTransform(math.NewVector2(settings.GlobalSettings.VideoScaleMultiplier, settings.GlobalSettings.VideoScaleMultiplier))

	animated := 0

	for _, object := range frame {
		obEntry := *BakeRenderedObjectsFillables(object)
		object = &obEntry

		if _, ok := r.Statistics[object.ObjectId]; !ok {
			r.Statistics[object.ObjectId] = rendererStatsEntry{
				Frames: &atomic.Uint64{},
				Lines:  &atomic.Uint64{},
				Size:   &atomic.Uint64{},
			}
		}

		object.MatrixTransform = scale.Multiply(object.MatrixTransform) //TODO: order?

		depth := object.GetDepth()

		var tagsToTransition []*line.EventLine

		for i := len(r.RunningBuffer) - 1; i >= 0; i-- {
			tag := r.RunningBuffer[i]
			if depth.Equals(tag.Layer) && object.ObjectId == tag.ObjectId {
				tagsToTransition = append(tagsToTransition, tag)
				r.RunningBuffer = slices.Delete(r.RunningBuffer, i, i+1)
			}
		}
		slices.Reverse(tagsToTransition)

		canTransition := true
		var transitionedTags []*line.EventLine

		for _, tag := range tagsToTransition {
			tag = tag.Transition(frameInfo, object)
			if tag != nil {
				transitionedTags = append(transitionedTags, tag)
				tag.DropCache()
			} else {
				canTransition = false
				break
			}
		}

		r.Statistics[object.ObjectId].Frames.Add(1)

		if canTransition && len(transitionedTags) > 0 {
			animated += len(transitionedTags)
			runningBuffer = append(runningBuffer, transitionedTags...)
		} else {
			r.RunningBuffer = append(r.RunningBuffer, tagsToTransition...)

			for _, l := range line.EventLinesFromRenderObject(frameInfo, object, settings.GlobalSettings.BakeMatrixTransforms) {
				l.DropCache()
				runningBuffer = append(runningBuffer, l)
			}
		}
	}

	fmt.Printf("[ASS] Total %d objects, %d flush, %d buffer, %d animated tags.\n", len(frame), len(r.RunningBuffer), len(runningBuffer), animated)

	//Flush non dupes
	result = append(result, r.Flush(frameInfo)...)
	r.RunningBuffer = runningBuffer

	return result
}

func threadedRenderer(stats map[uint16]rendererStatsEntry, buf []*line.EventLine, duration time.Duration) []string {
	if len(buf) == 0 {
		return nil
	}
	results := make([]string, len(buf))
	var cnt atomic.Uint64
	var wg sync.WaitGroup
	for i := 0; i < min(len(buf), runtime.NumCPU()); i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()

			for {
				i := cnt.Add(1) - 1
				if i >= uint64(len(buf)) {
					break
				}
				l := buf[i]
				stats[l.ObjectId].Lines.Add(1)

				l.Name += fmt.Sprintf(" f:%d>%d~%d", l.Start, l.End, l.End-l.Start+1)
				l.DropCache()
				encode := l.Encode(duration)
				results[i] = encode
				stats[l.ObjectId].Size.Add(uint64(len(encode)))
			}
		}(i)
	}
	wg.Wait()
	return results
}

func (r *Renderer) Flush(frameInfo types.FrameInformation) (result []string) {
	result = threadedRenderer(r.Statistics, r.RunningBuffer, frameInfo.GetFrameDuration())
	r.RunningBuffer = r.RunningBuffer[:0]
	return result
}

func BakeRenderedObjectsFillables(o *types.RenderedObject) *types.RenderedObject {
	var baked bool

	drawPathList := make(shapes.DrawPathList, 0, len(o.DrawPathList))

	for _, command := range o.DrawPathList {
		if fillStyleRecord, ok := command.Style.(*shapes.FillStyleRecord); ok && !fillStyleRecord.IsFlat() {
			baked = true
			flattened := fillStyleRecord.Flatten(command.Commands)
			drawPathList = append(drawPathList, flattened...)
		} else {
			drawPathList = append(drawPathList, command)
		}
	}

	if baked {
		return &types.RenderedObject{
			Depth:           o.Depth,
			ObjectId:        o.ObjectId,
			DrawPathList:    drawPathList,
			Clip:            o.Clip,
			ColorTransform:  o.ColorTransform,
			MatrixTransform: o.MatrixTransform,
		}
	} else {
		return o
	}
}

func RenderedObjectDepthSort(a, b *types.RenderedObject) int {
	if len(b.Depth) > len(a.Depth) {
		for i, depth := range b.Depth {
			var otherDepth uint16
			if i < len(a.Depth) {
				otherDepth = a.Depth[i]
			}

			if depth != otherDepth {
				return int(otherDepth) - int(depth)
			}
		}
	} else {
		for i, depth := range a.Depth {
			var otherDepth uint16
			if i < len(b.Depth) {
				otherDepth = b.Depth[i]
			}

			if depth != otherDepth {
				return int(depth) - int(otherDepth)
			}
		}
	}

	return 0
}
