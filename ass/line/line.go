package line

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/tag"
	asstime "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EventLine struct {
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

	Tags []tag.Tag

	cachedEncode *string
}

const StyleFill = "fill"
const StyleLine = "line"

func (l *EventLine) GetStart() int64 {
	return l.Start
}

func (l *EventLine) GetEnd() int64 {
	return l.End
}

func (l *EventLine) Transition(frameInfo types.FrameInformation, object *types.RenderedObject) *EventLine {
	line := *l
	line.End = frameInfo.GetFrameNumber()
	line.Tags = make([]tag.Tag, 0, len(l.Tags))
	//TODO: clip?

	if object.GetDepth().Equals(l.Layer) && object.ObjectId == l.ObjectId {
		if len(object.DrawPathList) <= line.ShapeIndex {
			return nil
		}
		command := object.DrawPathList[line.ShapeIndex]
		for _, t := range l.Tags {
			if positioningTag, ok := t.(tag.PositioningTag); ok {
				t = positioningTag.TransitionMatrixTransform(&line, object.MatrixTransform)
				if t == nil {
					return nil
				}
			}
			if colorTag, ok := t.(tag.ColorTag); ok {
				t = colorTag.TransitionColor(&line, object.ColorTransform)
				if t == nil {
					return nil
				}
			}
			if colorTag, ok := t.(tag.PathTag); ok {
				t = colorTag.TransitionShape(&line, command.Commands)
				if t == nil {
					return nil
				}
			}
			if colorTag, ok := t.(tag.ClipPathTag); ok {
				t = colorTag.TransitionClipPath(&line, object.Clip)
				if t == nil {
					return nil
				}
			}
			line.Tags = append(line.Tags, t)
		}
	}
	line.DropCache()
	return &line
}

func (l *EventLine) Encode(frameDuration time.Duration) string {
	if frameDuration == 1000*time.Millisecond && l.cachedEncode != nil {
		return *l.cachedEncode
	}

	eventTime := asstime.NewEventTime(l.Start, l.End-l.Start+1, frameDuration)

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

	eventTimeTags := eventTime.Encode()

	if len(eventTimeTags) > 0 {
		text = append(text, eventTimeTags)
	}

	for _, t := range l.Tags {
		text = append(text, "{"+t.Encode(eventTime)+"}")
	}

	line = append(line, strings.Join(text, ""))

	event := strings.Join(line, ",")

	if frameDuration == 1000*time.Millisecond {
		l.cachedEncode = &event
	}

	return event
}

func (l *EventLine) DropCache() {
	l.cachedEncode = nil
}

func (l *EventLine) Equalish(o *EventLine) bool {
	return l.ObjectId == o.ObjectId &&
		len(l.Tags) == len(o.Tags) &&
		l.Layer.Equals(o.Layer) &&
		l.Encode(1000*time.Millisecond) == o.Encode(1000*time.Millisecond)
}

var eventLineRegexp = regexp.MustCompile(`^(?P<Kind>[^:]+): (?P<Layer>\d+),(?P<StartTimecode>[\d:.]+),(?P<EndTimecode>[\d:.]+),(?P<Style>[^,]*),(?P<Name>[^,]*),(?P<MarginL>\d+),(?P<MarginR>\d+),(?P<MarginV>\d+),(?P<Effect>[^,]*),(?P<Text>.*)$`)

func EventLineFromString(line string) (out *EventLine) {
	var l EventLine

	matches := eventLineRegexp.FindStringSubmatch(strings.TrimSpace(line))
	if matches == nil {
		return nil
	}

	var start, end asstime.Time
	var text string
	var err error
	for i, name := range eventLineRegexp.SubexpNames() {
		val := matches[i]
		switch name {
		case "Kind":
			if val == "Dialogue" {
				l.IsComment = false
			} else if val == "Comment" {
				l.IsComment = true
			} else {
				return nil
			}
		case "Layer":
			layer, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return nil
			}
			l.Layer = types.DepthFromPackedLayer(int32(layer))
		case "StartTimecode":
			start, err = asstime.FromString(val)
			if err != nil {
				return nil
			}
		case "EndTimecode":
			end, err = asstime.FromString(val)
			if err != nil {
				return nil
			}
		case "Style":
			l.Style = val
		case "Name":
			l.Style = val
		case "MarginL":
			l.Margin.Left, err = strconv.ParseInt(val, 10, 0)
			if err != nil {
				return nil
			}
		case "MarginR":
			l.Margin.Right, err = strconv.ParseInt(val, 10, 0)
			if err != nil {
				return nil
			}
		case "MarginV":
			l.Margin.Vertical, err = strconv.ParseInt(val, 10, 0)
			if err != nil {
				return nil
			}
		case "Effect":
			l.Effect = val
		case "Text":
			text = val
		case "":
			continue
		default:
			panic("not implemented")

		}
	}

	eventTime := asstime.EventLineFromText(start, end, text)

	_ = eventTime

	return out
}

func EventLinesFromRenderObject(frameInfo types.FrameInformation, object *types.RenderedObject, bakeMatrixTransforms bool) (out []*EventLine) {
	out = make([]*EventLine, 0, len(object.DrawPathList))
	for i, drawPath := range object.DrawPathList {
		style := ""
		if _, ok := drawPath.Style.(*shapes.FillStyleRecord); ok {
			style = StyleFill
		} else if _, ok = drawPath.Style.(*shapes.LineStyleRecord); ok {
			style = StyleLine
		} else {
			panic("unsupported")
		}
		out = append(out, &EventLine{
			Layer:      object.GetDepth(),
			ShapeIndex: i,
			ObjectId:   object.ObjectId,
			Start:      frameInfo.GetFrameNumber(),
			End:        frameInfo.GetFrameNumber(),
			Tags:       []tag.Tag{tag.ContainerTagFromPathEntry(drawPath, object.Clip, object.ColorTransform, object.MatrixTransform, bakeMatrixTransforms)},
			Name:       fmt.Sprintf("o:%d d:%s", object.ObjectId, object.GetDepth().String()),
			Style:      style,
		})
	}
	return out
}
