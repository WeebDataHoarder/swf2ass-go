package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/line"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type FillColorTag struct {
	Color         *math.Color
	OriginalColor *math.Color
}

func (t *FillColorTag) FromStyleRecord(record shapes.StyleRecord) StyleTag {
	if fillStyleRecord, ok := record.(*shapes.FillStyleRecord); ok {
		if color, ok := fillStyleRecord.Fill.(math.Color); ok {
			t.Color = &color
			t.OriginalColor = &color
		} else if gradient, ok := fillStyleRecord.Fill.(shapes.Gradient); ok {
			items := gradient.GetItems()
			t.Color = &items[0].Color
			t.OriginalColor = &items[0].Color
			panic("Gradient fill not supported here")
		} else {
			panic("not implemented")
		}
	} else {
		t.OriginalColor = nil
		t.Color = nil
	}
	return t
}

func (t *FillColorTag) TransitionStyleRecord(line *line.Line, record shapes.StyleRecord) StyleTag {
	t2 := &LineColorTag{}
	t2.FromStyleRecord(record)
	return t2
}

func (t *FillColorTag) ApplyColorTransform(transform math.ColorTransform) ColorTag {
	color := t.Color
	if t.OriginalColor != nil {
		color2 := transform.ApplyToColor(*t.OriginalColor)
		color = &color2
	}
	return &FillColorTag{
		Color:         color,
		OriginalColor: t.OriginalColor,
	}
}

func (t *FillColorTag) TransitionColor(line *line.Line, transform math.ColorTransform) ColorTag {
	return t.ApplyColorTransform(transform)
}

func (t *FillColorTag) Equals(tag Tag) bool {
	if o, ok := tag.(*LineColorTag); ok {
		return (t.Color == o.Color || (t.Color != nil && t.Color.Equals(*o.Color, true))) && (t.OriginalColor == o.OriginalColor || (t.OriginalColor != nil && t.OriginalColor.Equals(*o.OriginalColor, true)))
	}
	return false
}

func (t *FillColorTag) Encode(event time.EventTime) string {
	if t.Color == nil {
		return "\\1a&HFF&"
	} else {
		return fmt.Sprintf("\\1c&H%02X%02X%02X&\\1a&H%02X&", t.Color.B, t.Color.G, t.Color.R, 255-t.Color.Alpha)
	}
}
