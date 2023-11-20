package types

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

type StyleList struct {
	FillStyles []*FillStyleRecord
	LineStyles []*LineStyleRecord
}

func (l *StyleList) GetFillStyle(i int) *FillStyleRecord {
	if len(l.FillStyles) > i {
		return l.FillStyles[i]
	}
	return nil
}

func (l *StyleList) GetLineStyle(i int) *LineStyleRecord {
	if len(l.LineStyles) > i {
		return l.LineStyles[i]
	}
	return nil
}

func StyleListFromSWFItems(fillStyles subtypes.FILLSTYLEARRAY, lineStyles subtypes.LINESTYLEARRAY) (r StyleList) {
	for _, s := range fillStyles.FillStyles {
		r.FillStyles = append(r.FillStyles, FillStyleRecordFromSWFFILLSTYLE(s))
	}

	if len(lineStyles.LineStyles) > 0 {
		for _, s := range lineStyles.LineStyles {
			r.LineStyles = append(r.LineStyles, &LineStyleRecord{
				//TODO: any reason for  max(types.TwipFactor)?
				Width: max(types.Twip(s.Width), types.TwipFactor),
				Color: math.Color{
					R:     s.Color.R(),
					G:     s.Color.G(),
					B:     s.Color.B(),
					Alpha: s.Color.A(),
				},
			})
		}
	} else if len(lineStyles.LineStyles2) > 0 {
		for _, s := range lineStyles.LineStyles2 {
			if !s.Flag.HasFill {
				r.LineStyles = append(r.LineStyles, &LineStyleRecord{
					//TODO: any reason for  max(types.TwipFactor)?
					Width: max(types.Twip(s.Width), types.TwipFactor),
					Color: math.Color{
						R:     s.Color.R(),
						G:     s.Color.G(),
						B:     s.Color.B(),
						Alpha: s.Color.A(),
					},
				})
			} else {
				fill := FillStyleRecordFromSWFFILLSTYLE(s.FillType)
				switch fillEntry := fill.Fill.(type) {
				case types.Color:
					r.LineStyles = append(r.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: max(types.Twip(s.Width), types.TwipFactor),
						Color: math.Color{
							R:     fillEntry.R(),
							G:     fillEntry.G(),
							B:     fillEntry.B(),
							Alpha: fillEntry.A(),
						},
					})
				case Gradient:
					//TODO: gradient fill of lines
					color := fillEntry.GetItems()[0].Color
					r.LineStyles = append(r.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: max(types.Twip(s.Width), types.TwipFactor),
						Color: color,
					})
				}
			}
		}
	}

	return r
}
