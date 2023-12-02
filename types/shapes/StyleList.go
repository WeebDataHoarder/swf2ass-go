package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"
)

type StyleList struct {
	FillStyles []*FillStyleRecord
	LineStyles []*LineStyleRecord
}

func (l StyleList) GetFillStyle(i int) *FillStyleRecord {
	if len(l.FillStyles) > i {
		return l.FillStyles[i]
	}
	return nil
}

func (l StyleList) GetLineStyle(i int) *LineStyleRecord {
	if len(l.LineStyles) > i {
		return l.LineStyles[i]
	}
	return nil
}

func StyleListFromSWFItems(collection ObjectCollection, fillStyles subtypes.FILLSTYLEARRAY, lineStyles subtypes.LINESTYLEARRAY) (r StyleList) {
	for _, s := range fillStyles.FillStyles {
		r.FillStyles = append(r.FillStyles, FillStyleRecordFromSWF(collection, s.FillStyleType, s.Color, s.Gradient, s.FocalGradient, s.GradientMatrix, s.BitmapMatrix, s.BitmapId))
	}

	if len(lineStyles.LineStyles) > 0 {
		for _, s := range lineStyles.LineStyles {
			r.LineStyles = append(r.LineStyles, LineStyleRecordFromSWF(s.Width, 0, false, s.Color, nil))
		}
	} else if len(lineStyles.LineStyles2) > 0 {
		for _, s := range lineStyles.LineStyles2 {
			if s.Flag.HasFill {
				r.LineStyles = append(r.LineStyles, LineStyleRecordFromSWF(s.Width, 0,
					s.Flag.HasFill,
					s.Color,
					FillStyleRecordFromSWF(collection, s.FillType.FillStyleType, s.FillType.Color, s.FillType.Gradient, s.FillType.FocalGradient, s.FillType.GradientMatrix, s.FillType.BitmapMatrix, s.FillType.BitmapId),
				))
			} else {
				r.LineStyles = append(r.LineStyles, LineStyleRecordFromSWF(s.Width, 0, false, s.Color, nil))
			}
		}
	}

	return r
}

func StyleListFromSWFMorphItems(collection ObjectCollection, fillStyles subtypes.MORPHFILLSTYLEARRAY, lineStyles subtypes.MORPHLINESTYLEARRAY) (start, end StyleList) {
	for _, s := range fillStyles.FillStyles {
		startStyle, endStyle := FillStyleRecordFromSWFMORPHFILLSTYLE(collection, s)
		start.FillStyles = append(start.FillStyles, startStyle)
		end.FillStyles = append(end.FillStyles, endStyle)
	}

	if len(lineStyles.LineStyles) > 0 {
		for _, s := range lineStyles.LineStyles {
			startStyle, endStyle := LineStyleRecordFromSWFMORPHLINESTYLE(s)
			start.LineStyles = append(start.LineStyles, startStyle)
			end.LineStyles = append(end.LineStyles, endStyle)
		}
	} else if len(lineStyles.LineStyles2) > 0 {
		for _, s := range lineStyles.LineStyles2 {
			startStyle, endStyle := LineStyleRecordFromSWFMORPHLINESTYLE2(collection, s)
			start.LineStyles = append(start.LineStyles, startStyle)
			end.LineStyles = append(end.LineStyles, endStyle)
		}
	}

	return start, end
}
