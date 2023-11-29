package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type LineColorTag struct {
	Color         *math.Color
	OriginalColor *math.Color
}

func (t *LineColorTag) FromStyleRecord(record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	if lineStyleRecord, ok := record.(*shapes.LineStyleRecord); ok {
		t.Color = &lineStyleRecord.Color
		t.OriginalColor = &lineStyleRecord.Color
	} else if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok && fillStyleRecord.Border != nil {
		t.Color = &fillStyleRecord.Border.Color
		t.OriginalColor = &fillStyleRecord.Border.Color
	} else {
		t.OriginalColor = nil
		t.Color = nil
	}
	return t
}

func (t *LineColorTag) TransitionStyleRecord(event Event, record shapes.StyleRecord, transform math.MatrixTransform) StyleTag {
	t2 := &LineColorTag{}
	t2.FromStyleRecord(record, transform)
	return t2
}

func (t *LineColorTag) ApplyColorTransform(transform math.ColorTransform) ColorTag {
	color := t.Color
	if t.OriginalColor != nil {
		color2 := transform.ApplyToColor(*t.OriginalColor)
		color = &color2
	}
	return &LineColorTag{
		Color:         color,
		OriginalColor: t.OriginalColor,
	}
}

func (t *LineColorTag) TransitionColor(event Event, transform math.ColorTransform) ColorTag {
	return t.ApplyColorTransform(transform)
}

func (t *LineColorTag) Equals(tag Tag) bool {
	if o, ok := tag.(*LineColorTag); ok {
		return (t.Color == o.Color || (t.Color != nil && o.Color != nil && t.Color.Equals(*o.Color, true))) && (t.OriginalColor == o.OriginalColor || (t.OriginalColor != nil && o.OriginalColor != nil && t.OriginalColor.Equals(*o.OriginalColor, true)))
	}
	return false
}

func (t *LineColorTag) Encode(event time.EventTime) string {
	if t.Color == nil {
		return "\\3a&HFF&"
	} else {
		return fmt.Sprintf("\\3c&H%02X%02X%02X&\\3a&H%02X&", t.Color.B, t.Color.G, t.Color.R, 255-t.Color.Alpha)
	}
}
