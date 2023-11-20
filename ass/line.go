package ass

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"strconv"
	"strings"
	"time"
)

type Line struct {
	Layer      types.Depth
	ShapeIndex int
	ObjectId   uint16

	Start, End int64

	Style string
	Name  string

	Margin struct {
		Left     int64
		Right    int64
		Vertical int64
	}
	Effect string

	IsComment bool

	Tags []Tag

	cachedEncode *string
}

func (l *Line) Encode(frameDuration time.Duration) string {
	if frameDuration == 1000*time.Millisecond && l.cachedEncode != nil {
		return *l.cachedEncode
	}

	eventTime := NewEventTime(uint64(l.Start), uint64(l.End-l.Start+1), frameDuration)

	line := make([]string, 0, 10)
	if l.IsComment {
		line = append(line, fmt.Sprintf("Comment: %d", l.Layer.GetPackedLayer()))
	} else {
		line = append(line, fmt.Sprintf("Dialogue: %d", l.Layer.GetPackedLayer()))
	}

	line = append(line, eventTime.Start.String())
	line = append(line, eventTime.End.String())
	line = append(line, l.Style)
	line = append(line, l.Name)
	line = append(line, strconv.FormatInt(l.Margin.Left, 10))
	line = append(line, strconv.FormatInt(l.Margin.Right, 10))
	line = append(line, strconv.FormatInt(l.Margin.Vertical, 10))
	line = append(line, l.Effect)

	text := make([]string, 0, 1+len(l.Tags))

	if eventTime.Start.AdjustedMillisecondError != 0 || eventTime.End.AdjustedMillisecondError != 0 {
		//Adjust frame precision exactly to frame boundaries. This is necessary due to low ASS timing precision
		//TODO: Maybe use fade?
		frameStartTime := eventTime.GetDurationFromStartOffset(0).Milliseconds()
		frameEndTime := eventTime.GetDurationFromEndOffset(0).Milliseconds()
		//TODO: maybe needs to be -1?
		text = append(text, fmt.Sprintf("{\\fade(255,0,255,%d,%d,%d,%d)\\err(%d~%d,%d~%d)}", frameStartTime, frameStartTime, frameEndTime, frameEndTime, eventTime.Start.Milliseconds, eventTime.Start.AdjustedMillisecondError, eventTime.End.Milliseconds, eventTime.End.AdjustedMillisecondError))
	}

	for _, t := range l.Tags {
		text = append(text, "{"+t.Encode(eventTime)+"}")
	}

	line = append(line, strings.Join(line, ""))

	event := strings.Join(line, ",")

	if frameDuration == 1000*time.Millisecond {
		l.cachedEncode = &event
	}

	return event
}

func (l *Line) DropCache() {
	l.cachedEncode = nil
}

func (l *Line) Equalish(o *Line) bool {
	return l.ObjectId == o.ObjectId &&
		len(l.Tags) == len(o.Tags) &&
		l.Layer.Equals(o.Layer) &&
		l.Encode(1000*time.Millisecond) == o.Encode(1000*time.Millisecond)
}

func LinesFromRenderObject(frameInfo types.FrameInformation, object types.RenderedObject, bakeTransforms bool) (out []Line) {
	out = make([]Line, 0, len(object.DrawPathList))
	for i := range object.DrawPathList {
		out = append(out, Line{
			Layer:      object.GetDepth(),
			ShapeIndex: i,
			ObjectId:   object.ObjectId,
			Start:      frameInfo.GetFrameNumber(),
			End:        frameInfo.GetFrameNumber(),
			Tags:       []Tag{ContainerTagFromPathEntry(object.DrawPathList[i], object.Clip, object.ColorTransform, object.MatrixTransform, bakeTransforms)},
			Name:       fmt.Sprintf("o:%d d:%s", object.ObjectId, object.GetDepth().String()),
		})
	}
	return out
}
