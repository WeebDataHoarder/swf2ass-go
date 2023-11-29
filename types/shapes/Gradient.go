package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
	"slices"
)

type Gradient struct {
	Records []GradientItem

	Transform         math2.MatrixTransform
	SpreadMode        swfsubtypes.GradientSpreadMode
	InterpolationMode swfsubtypes.GradientInterpolationMode

	Interpolation func(self Gradient, overlap, blur float64, gradientSlices int, bb Rectangle[float64]) DrawPathList
}

func (g Gradient) GetItems() []GradientItem {
	return g.Records
}

func (g Gradient) GetInterpolatedDrawPaths(overlap, blur float64, gradientSlices int, bb Rectangle[float64]) DrawPathList {
	return g.Interpolation(g, overlap, blur, gradientSlices, bb).ApplyMatrixTransform(g.Transform, true).(DrawPathList)
}

func (g Gradient) ApplyMatrixTransform(transform math2.MatrixTransform, applyTranslation bool) Fillable {
	if transform.IsIdentity() {
		return g
	}
	g2 := g
	if !applyTranslation {
		panic("not supported")
	}
	g2.Transform = transform.Combine(g2.Transform)
	return g2
}

func (g Gradient) ApplyColorTransform(transform math2.ColorTransform) Fillable {
	g2 := g
	g2.Records = slices.Clone(g2.Records)
	for i, g := range g2.Records {
		g2.Records[i] = GradientItem{
			Ratio: g.Ratio,
			Color: transform.ApplyToColor(g.Color),
		}
	}
	return g2
}

func (g Gradient) Fill(shape Shape) DrawPathList {
	bb := Rectangle[float64]{}
	if inverse := g.Transform.Inverse(); inverse != nil {
		bb = shape.ApplyMatrixTransform(*inverse, true).BoundingBox()
	}
	return g.GetInterpolatedDrawPaths(settings.GlobalSettings.GradientOverlap, settings.GlobalSettings.GradientBlur, settings.GlobalSettings.GradientSlices, bb).Fill(shape)
}

type GradientItem struct {
	Ratio uint8
	Color math2.Color
}

type GradientSlice struct {
	Start, End float64
	Color      math2.Color
}

const GradientBoundsMin swftypes.Twip = math.MinInt16 / 2
const GradientBoundsMax = -GradientBoundsMin

var GradientBounds = Rectangle[float64]{
	TopLeft:     math2.NewVector2[float64](GradientBoundsMin.Float64(), GradientBoundsMin.Float64()),
	BottomRight: math2.NewVector2[float64](GradientBoundsMax.Float64(), GradientBoundsMax.Float64()),
}

const GradientRatioDivisor = math.MaxUint8

func InterpolateGradient(gradient Gradient, gradientSlices int) (result []GradientSlice) {
	items := gradient.GetItems()
	//TODO: spread modes

	interpolationMode := gradient.InterpolationMode

	first := items[0]
	last := items[len(items)-1]
	if first.Ratio != 0 {
		first = GradientItem{
			Ratio: 0,
			Color: first.Color,
		}
		items = slices.Insert(items, 0, first)
	}

	if last.Ratio != GradientRatioDivisor {
		last = GradientItem{
			Ratio: GradientRatioDivisor,
			Color: last.Color,
		}
		items = append(items, last)
	}

	prevItem := items[0]
	for _, item := range items[1:] {
		prevColor := prevItem.Color
		currentColor := item.Color
		if interpolationMode == swfsubtypes.GradientInterpolationLinearRGB {
			prevColor = prevColor.ToLinearRGB()
			currentColor = currentColor.ToLinearRGB()
		}

		maxColorDistance := max(math.Abs(float64(prevColor.R)-float64(currentColor.R)), math.Abs(float64(prevColor.G)-float64(currentColor.G)), math.Abs(float64(prevColor.B)-float64(currentColor.B)), math.Abs(float64(prevColor.Alpha)-float64(currentColor.Alpha)))

		prevPosition := float64(prevItem.Ratio)
		currentPosition := float64(item.Ratio)
		distance := math.Abs(currentPosition - prevPosition)

		var partitions int
		if maxColorDistance < math.SmallestNonzeroFloat64 {
			partitions = 1
		} else if gradientSlices == settings.GradientAutoSlices {
			partitions = max(1, int(math.Ceil(min(GradientRatioDivisor/float64(len(items)+1), max(1, math.Ceil(maxColorDistance))))))
		} else {
			partitions = max(1, int(math.Ceil((distance/GradientRatioDivisor)*float64(gradientSlices))))
		}

		fromPos := prevPosition
		for i := 1; i <= partitions; i++ {
			ratio := float64(i) / float64(partitions)
			color := math2.LerpColor(prevColor, currentColor, ratio)

			if interpolationMode == swfsubtypes.GradientInterpolationLinearRGB {
				color = color.ToSRGB()
			}

			toPos := math2.Lerp(prevPosition, currentPosition, ratio)

			result = append(result, GradientSlice{
				Start: fromPos / GradientRatioDivisor,
				End:   toPos / GradientRatioDivisor,
				Color: color,
			})
			fromPos = toPos
		}
		prevItem = item
	}
	return result
}

func GradientItemFromSWF(ratio uint8, color swftypes.Color) GradientItem {
	return GradientItem{
		Ratio: ratio,
		Color: math2.Color{
			R:     color.R(),
			G:     color.G(),
			B:     color.B(),
			Alpha: color.A(),
		},
	}
}
