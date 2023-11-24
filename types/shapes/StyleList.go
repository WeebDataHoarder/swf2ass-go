package shapes

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

func StyleListFromSWFItems(collection ObjectCollection, fillStyles subtypes.FILLSTYLEARRAY, lineStyles subtypes.LINESTYLEARRAY) (r StyleList) {
	for _, s := range fillStyles.FillStyles {
		r.FillStyles = append(r.FillStyles, FillStyleRecordFromSWF(collection, s.FillStyleType, s.Color, s.Gradient, s.GradientMatrix, s.BitmapMatrix, s.BitmapId))
	}

	if len(lineStyles.LineStyles) > 0 {
		for _, s := range lineStyles.LineStyles {
			r.LineStyles = append(r.LineStyles, &LineStyleRecord{
				//TODO: any reason for  max(types.TwipFactor)?
				Width: types.Twip(s.Width).Float64(),
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
					Width: types.Twip(s.Width).Float64(),
					Color: math.Color{
						R:     s.Color.R(),
						G:     s.Color.G(),
						B:     s.Color.B(),
						Alpha: s.Color.A(),
					},
				})
			} else {
				fill := FillStyleRecordFromSWF(collection, s.FillType.FillStyleType, s.FillType.Color, s.FillType.Gradient, s.FillType.GradientMatrix, s.FillType.BitmapMatrix, s.FillType.BitmapId)
				switch fillEntry := fill.Fill.(type) {
				case types.Color:
					r.LineStyles = append(r.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: types.Twip(s.Width).Float64(),
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
						Width: types.Twip(s.Width).Float64(),
						Color: color,
					})
				}
			}
		}
	}

	return r
}

func StyleListFromSWFMorphItems(collection ObjectCollection, fillStyles subtypes.MORPHFILLSTYLEARRAY, lineStyles subtypes.MORPHLINESTYLEARRAY) (start, end StyleList) {
	for _, s := range fillStyles.FillStyles {
		start.FillStyles = append(start.FillStyles, FillStyleRecordFromSWFMORPHFILLSTYLEStart(collection, s))
		end.FillStyles = append(end.FillStyles, FillStyleRecordFromSWFMORPHFILLSTYLEEnd(collection, s))
	}

	if len(lineStyles.LineStyles) > 0 {
		for _, s := range lineStyles.LineStyles {
			start.LineStyles = append(start.LineStyles, &LineStyleRecord{
				//TODO: any reason for  max(types.TwipFactor)?
				Width: types.Twip(s.StartWidth).Float64(),
				Color: math.Color{
					R:     s.StartColor.R(),
					G:     s.StartColor.G(),
					B:     s.StartColor.B(),
					Alpha: s.StartColor.A(),
				},
			})

			end.LineStyles = append(end.LineStyles, &LineStyleRecord{
				//TODO: any reason for  max(types.TwipFactor)?
				Width: types.Twip(s.EndWidth).Float64(),
				Color: math.Color{
					R:     s.EndColor.R(),
					G:     s.EndColor.G(),
					B:     s.EndColor.B(),
					Alpha: s.EndColor.A(),
				},
			})
		}
	} else if len(lineStyles.LineStyles2) > 0 {
		for _, s := range lineStyles.LineStyles2 {
			if !s.Flag.HasFill {
				start.LineStyles = append(start.LineStyles, &LineStyleRecord{
					//TODO: any reason for  max(types.TwipFactor)?
					Width: types.Twip(s.StartWidth).Float64(),
					Color: math.Color{
						R:     s.StartColor.R(),
						G:     s.StartColor.G(),
						B:     s.StartColor.B(),
						Alpha: s.StartColor.A(),
					},
				})
				end.LineStyles = append(end.LineStyles, &LineStyleRecord{
					//TODO: any reason for  max(types.TwipFactor)?
					Width: types.Twip(s.EndWidth).Float64(),
					Color: math.Color{
						R:     s.EndColor.R(),
						G:     s.EndColor.G(),
						B:     s.EndColor.B(),
						Alpha: s.EndColor.A(),
					},
				})
			} else {
				fillStart := FillStyleRecordFromSWFMORPHFILLSTYLEStart(collection, s.FillType)
				fillEnd := FillStyleRecordFromSWFMORPHFILLSTYLEEnd(collection, s.FillType)
				switch fillEntry := fillStart.Fill.(type) {
				case types.Color:
					start.LineStyles = append(start.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: types.Twip(s.StartWidth).Float64(),
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
					start.LineStyles = append(start.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: types.Twip(s.StartWidth).Float64(),
						Color: color,
					})
				}
				switch fillEntry := fillEnd.Fill.(type) {
				case types.Color:
					end.LineStyles = append(end.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: types.Twip(s.EndWidth).Float64(),
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
					end.LineStyles = append(end.LineStyles, &LineStyleRecord{
						//TODO: any reason for  max(types.TwipFactor)?
						Width: types.Twip(s.EndWidth).Float64(),
						Color: color,
					})
				}
			}
		}
	}

	return start, end
}
