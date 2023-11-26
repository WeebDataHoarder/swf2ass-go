package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

func LerpFillStyle(start, end *FillStyleRecord, ratio float64) *FillStyleRecord {
	if start == nil || end == nil {
		return nil
	}

	return &FillStyleRecord{
		Fill:   LerpFillable(start.Fill, end.Fill, ratio),
		Border: LerpLineStyle(start.Border, end.Border, ratio),
		Blur:   math.Lerp(start.Blur, end.Blur, ratio),
	}
}

func LerpLineStyle(start, end *LineStyleRecord, ratio float64) *LineStyleRecord {
	if start == nil || end == nil {
		return nil
	}

	return &LineStyleRecord{
		Width: math.Lerp(start.Width, end.Width, ratio),
		Color: math.LerpColor(start.Color, end.Color, ratio),
		Blur:  math.Lerp(start.Blur, end.Blur, ratio),
	}
}

func LerpFillable(start, end any, ratio float64) any {
	switch s := start.(type) {
	case math.Color:
		return math.LerpColor(s, end.(math.Color), ratio)
	case Bitmap:
		return Bitmap{
			List:      s.List,
			Transform: math.LerpMatrix(s.Transform, end.(Bitmap).Transform, ratio),
		}
	case Gradient:
		return LerpGradient(s, end.(Gradient), ratio)
	case DrawPathList:
		return start
	//TODO: focal gradient
	default:
		panic("not supported")
	}
}

func LerpGradient(start, end Gradient, ratio float64) Gradient {

	startRecords := start.GetItems()
	endRecords := end.GetItems()

	if len(startRecords) != len(endRecords) {
		panic("not supported")
	}

	records := make([]GradientItem, 0, len(startRecords))
	for i := range startRecords {
		records = append(records, GradientItem{
			Ratio: math.Lerp(startRecords[i].Ratio, endRecords[i].Ratio, ratio),
			Color: math.LerpColor(startRecords[i].Color, endRecords[i].Color, ratio),
		})
	}

	return Gradient{
		Records:           records,
		Transform:         math.LerpMatrix(start.Transform, end.Transform, ratio),
		SpreadMode:        start.SpreadMode,
		InterpolationMode: start.InterpolationMode,
		Interpolation:     start.Interpolation,
	}
}
