package shapes

import (
	"fmt"
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	math2 "math"
	"slices"
)

type StyleRecord interface {
	ApplyColorTransform(transform math.ColorTransform) StyleRecord
	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) StyleRecord
}

type Fillable interface {
	Fill(shape Shape) DrawPathList
	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) Fillable
	ApplyColorTransform(transform math.ColorTransform) Fillable
}

type FillStyleRecord struct {
	// Fill can be math.Color or Fillable
	Fill      any
	Border    *LineStyleRecord
	Blur      float64
	fillCache struct {
		Shape *Shape
		List  DrawPathList
	}
}

func (r *FillStyleRecord) IsFlat() bool {
	_, ok := r.Fill.(math.Color)
	return ok
}

// Flatten Creates a fill that is only composed of FillStyleRecord with Fill being math.Color
func (r *FillStyleRecord) Flatten(s Shape) DrawPathList {
	if _, ok := r.Fill.(math.Color); ok {
		return DrawPathList{
			{
				Style: r,
				Shape: s,
			},
		}
	} else if fillable, ok := r.Fill.(Fillable); ok {
		//TODO: inherit blur/border
		if r.fillCache.Shape != nil && r.fillCache.Shape.Equals(s) {
			return r.fillCache.List
		}
		fill := fillable.Fill(s)
		r.fillCache.List = fill
		s2 := slices.Clone(s)
		r.fillCache.Shape = &s2
		return fill
	} else {
		panic("not supported")
	}
}

func (r *FillStyleRecord) ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) StyleRecord {
	fill := r.Fill
	if color, ok := fill.(math.Color); ok {
		fill = color
	} else if fillable, ok := r.Fill.(Fillable); ok {
		fill = fillable.ApplyMatrixTransform(transform, applyTranslation)
	} else {
		panic("not supported")
	}

	if r.Border != nil {
		return &FillStyleRecord{
			Fill:   fill,
			Border: r.Border.ApplyMatrixTransform(transform, applyTranslation).(*LineStyleRecord),
			Blur:   r.Blur, //TODO: scale blur?
		}
	}
	return &FillStyleRecord{
		Fill:   fill,
		Border: nil,
		Blur:   r.Blur, //TODO: scale blur?
	}
}

func (r *FillStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	fill := r.Fill
	if color, ok := fill.(math.Color); ok {
		fill = transform.ApplyToColor(color)
	} else if fillable, ok := r.Fill.(Fillable); ok {
		fill = fillable.ApplyColorTransform(transform)
	} else {
		panic("not supported")
	}
	return &FillStyleRecord{
		Border: r.Border,
		Fill:   fill,
		Blur:   r.Blur,
	}
}

func FillStyleRecordFromSWF(collection ObjectCollection, fillType swfsubtypes.FillStyleType, color swftypes.Color, gradient swfsubtypes.GRADIENT, focalGradient swfsubtypes.FOCALGRADIENT, gradientMatrix, bitmapMatrix swftypes.MATRIX, bitmapId uint16) (r *FillStyleRecord) {
	switch fillType {
	case swfsubtypes.FillStyleSolid:
		return &FillStyleRecord{
			Fill: math.Color{
				R:     color.R(),
				G:     color.G(),
				B:     color.B(),
				Alpha: color.A(),
			},
		}
	case swfsubtypes.FillStyleLinearGradient:
		return &FillStyleRecord{
			Fill: LinearGradientFromSWF(gradient.Records, gradientMatrix, gradient.SpreadMode, gradient.InterpolationMode),
		}
	case swfsubtypes.FillStyleRadialGradient:
		return &FillStyleRecord{
			Fill: RadialGradientFromSWF(gradient.Records, gradientMatrix, gradient.SpreadMode, gradient.InterpolationMode),
		}
	case swfsubtypes.FillStyleFocalRadialGradient:
		//TODO: do it properly
		return &FillStyleRecord{
			Fill: RadialGradientFromSWF(focalGradient.Records, gradientMatrix, gradient.SpreadMode, gradient.InterpolationMode),
		}
	case swfsubtypes.FillStyleClippedBitmap, swfsubtypes.FillStyleRepeatingBitmap:
		if bitmapId == math2.MaxUint16 { //Special case, TODO:???
			return &FillStyleRecord{
				Fill: math.Color{
					R:     0,
					G:     0,
					B:     0,
					Alpha: 0,
				},
			}
		}
		bitmap := collection.Get(bitmapId)
		if bitmap == nil {
			fmt.Printf("bitmap %d not found!\n", bitmapId)
			return &FillStyleRecord{
				Fill: math.Color{
					R:     0,
					G:     0,
					B:     0,
					Alpha: 0,
				},
			}
		}

		//TODO: repeating

		//TODO: what blur factor should it pick
		blurFactor := 0.1
		//TODO: extend color
		return &FillStyleRecord{
			Fill: BitmapFillFromSWF(bitmap.GetShapeList(ObjectProperties{}).ApplyFunction(func(p DrawPath) DrawPath {
				if fillStyle, ok := p.Style.(*FillStyleRecord); ok {
					return DrawPathFill(&FillStyleRecord{
						Fill:   fillStyle.Fill,
						Border: fillStyle.Border,
						Blur:   blurFactor,
					}, p.Shape)
				}
				return p
			}), bitmapMatrix),
		}
	case swfsubtypes.FillStyleNonSmoothedClippedBitmap, swfsubtypes.FillStyleNonSmoothedRepeatingBitmap:
		if bitmapId == math2.MaxUint16 { //Special case, TODO:???
			return &FillStyleRecord{
				Fill: math.Color{
					R:     0,
					G:     0,
					B:     0,
					Alpha: 0,
				},
			}
		}

		//TODO: repeating

		bitmap := collection.Get(bitmapId)
		if bitmap == nil {
			fmt.Printf("bitmap %d not found!\n", bitmapId)
			return &FillStyleRecord{
				Fill: math.Color{
					R:     0,
					G:     0,
					B:     0,
					Alpha: 0,
				},
			}
		}
		//TODO: extend color
		return &FillStyleRecord{
			Fill: BitmapFillFromSWF(bitmap.GetShapeList(ObjectProperties{}), bitmapMatrix),
		}
		//TODO other styles
	}

	return &FillStyleRecord{
		Fill: math.Color{
			R:     0,
			G:     0,
			B:     0,
			Alpha: 0,
		},
	}
}

func FillStyleRecordFromSWFMORPHFILLSTYLE(collection ObjectCollection, fillStyle swfsubtypes.MORPHFILLSTYLE) (start, end *FillStyleRecord) {
	return FillStyleRecordFromSWF(collection, fillStyle.FillStyleType, fillStyle.StartColor, fillStyle.Gradient.StartGradient(), swfsubtypes.FOCALGRADIENT{}, fillStyle.StartGradientMatrix, fillStyle.StartBitmapMatrix, fillStyle.BitmapId),
		FillStyleRecordFromSWF(collection, fillStyle.FillStyleType, fillStyle.EndColor, fillStyle.Gradient.EndGradient(), swfsubtypes.FOCALGRADIENT{}, fillStyle.EndGradientMatrix, fillStyle.EndBitmapMatrix, fillStyle.BitmapId)
}

func LineStyleRecordFromSWF(width uint16, blur float64, hasFill bool, c swftypes.Color, fill *FillStyleRecord) (r *LineStyleRecord) {
	if hasFill {
		//TODO: do this properly
		switch fillType := fill.Fill.(type) {
		case math.Color:
			return &LineStyleRecord{
				Width: swftypes.Twip(width).Float64(),
				Color: fillType,
				Blur:  blur,
			}
		case Gradient:
			//TODO: gradient fill of lines
			return &LineStyleRecord{
				Width: swftypes.Twip(width).Float64(),
				Color: fillType.GetItems()[0].Color,
				Blur:  blur,
			}
			//TODO: other types, maybe generalize as a Fillable
		}
	}
	return &LineStyleRecord{
		Width: swftypes.Twip(width).Float64(),
		Color: math.Color{
			R:     c.R(),
			G:     c.G(),
			B:     c.B(),
			Alpha: c.A(),
		},
		Blur: blur,
	}
}

func LineStyleRecordFromSWFMORPHLINESTYLE(lineStyle swfsubtypes.MORPHLINESTYLE) (start, end *LineStyleRecord) {
	return LineStyleRecordFromSWF(lineStyle.StartWidth, 0, false, lineStyle.StartColor, nil),
		LineStyleRecordFromSWF(lineStyle.EndWidth, 0, false, lineStyle.EndColor, nil)
}

func LineStyleRecordFromSWFMORPHLINESTYLE2(collection ObjectCollection, lineStyle swfsubtypes.MORPHLINESTYLE2) (start, end *LineStyleRecord) {
	if lineStyle.Flag.HasFill {
		startFill, endFill := FillStyleRecordFromSWFMORPHFILLSTYLE(collection, lineStyle.FillType)
		return LineStyleRecordFromSWF(lineStyle.StartWidth, 0, lineStyle.Flag.HasFill, lineStyle.StartColor, startFill),
			LineStyleRecordFromSWF(lineStyle.EndWidth, 0, lineStyle.Flag.HasFill, lineStyle.EndColor, endFill)
	}
	return LineStyleRecordFromSWF(lineStyle.StartWidth, 0, false, lineStyle.StartColor, nil),
		LineStyleRecordFromSWF(lineStyle.EndWidth, 0, false, lineStyle.EndColor, nil)
}
