package shapes

import (
	"fmt"
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	math2 "math"
	"slices"
)

type StyleRecord interface {
	ApplyColorTransform(transform math.ColorTransform) StyleRecord
}

type LineStyleRecord struct {
	Width float64
	Color math.Color
	Blur  float64
}

func (r *LineStyleRecord) ApplyColorTransform(transform math.ColorTransform) StyleRecord {
	return &LineStyleRecord{
		Width: r.Width,
		Color: transform.ApplyToColor(r.Color),
		Blur:  r.Blur,
	}
}

type Fillable interface {
	Fill(shape *Shape) DrawPathList
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
func (r *FillStyleRecord) Flatten(s *Shape) DrawPathList {
	if _, ok := r.Fill.(math.Color); ok {
		return DrawPathList{
			{
				Style:    r,
				Commands: s,
			},
		}
	} else if fillable, ok := r.Fill.(Fillable); ok {
		//TODO: inherit blur/border
		if r.fillCache.Shape != nil && r.fillCache.Shape.Equals(s) {
			return r.fillCache.List
		}
		fill := fillable.Fill(s)
		r.fillCache.List = fill
		r.fillCache.Shape = &Shape{
			Edges:  slices.Clone(s.Edges),
			IsFlat: s.IsFlat,
		}
		return fill
	} else {
		panic("not supported")
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

func FillStyleRecordFromSWF(collection ObjectCollection, fillType swfsubtypes.FillStyleType, color swftypes.Color, gradient swfsubtypes.GRADIENT, gradientMatrix, bitmapMatrix swftypes.MATRIX, bitmapId uint16) (r *FillStyleRecord) {
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
			Fill: math.Color{
				R:     gradient.Records[0].Color.R(),
				G:     gradient.Records[0].Color.G(),
				B:     gradient.Records[0].Color.B(),
				Alpha: gradient.Records[0].Color.A(),
			},
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

		//TODO: what blur factor should it pick
		blurFactor := 1.0
		//TODO: extend color
		return &FillStyleRecord{
			Fill: DrawPathListFillFromSWF(bitmap.GetShapeList(0).ApplyFunction(func(p DrawPath) DrawPath {
				if fillStyle, ok := p.Style.(*FillStyleRecord); ok {
					return DrawPathFill(&FillStyleRecord{
						Fill:   fillStyle.Fill,
						Border: fillStyle.Border,
						Blur:   blurFactor,
					}, p.Commands, p.Clip)
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
			Fill: DrawPathListFillFromSWF(bitmap.GetShapeList(0), bitmapMatrix),
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

func FillStyleRecordFromSWFMORPHFILLSTYLEStart(collection ObjectCollection, fillStyle swfsubtypes.MORPHFILLSTYLE) (r *FillStyleRecord) {
	return FillStyleRecordFromSWF(collection, fillStyle.FillStyleType, fillStyle.StartColor, fillStyle.Gradient.StartGradient(), fillStyle.StartGradientMatrix, fillStyle.StartBitmapMatrix, fillStyle.BitmapId)
}

func FillStyleRecordFromSWFMORPHFILLSTYLEEnd(collection ObjectCollection, fillStyle swfsubtypes.MORPHFILLSTYLE) (r *FillStyleRecord) {
	return FillStyleRecordFromSWF(collection, fillStyle.FillStyleType, fillStyle.EndColor, fillStyle.Gradient.EndGradient(), fillStyle.EndGradientMatrix, fillStyle.EndBitmapMatrix, fillStyle.BitmapId)
}
