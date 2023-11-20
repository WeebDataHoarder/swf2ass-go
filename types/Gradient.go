package types

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
	"slices"
)

const GradientAutoSlices = -1

type Gradient interface {
	GetSpreadMode() swfsubtypes.GradientSpreadMode
	GetInterpolationMode() swfsubtypes.GradientInterpolationMode
	GetItems() []GradientItem
	GetInterpolatedDrawPaths(overlap int, slices int) DrawPathList
	GetMatrixTransform() math2.MatrixTransform
	ApplyColorTransform(transform math2.ColorTransform) Gradient
}

type GradientItem struct {
	Ratio uint8
	Color math2.Color
}

type GradientSlice struct {
	Start, End float64
	Color      math2.Color
}

func LerpGradient(gradient Gradient, gradientSlices int) (result []GradientSlice) {
	items := gradient.GetItems()
	//TODO: spread modes

	first := items[0]
	last := items[len(items)-1]

	interpolationMode := gradient.GetInterpolationMode()

	if first.Ratio != 0 {
		first = GradientItem{
			Ratio: 0,
			Color: first.Color,
		}
		items = slices.Insert(items, 0, first)
	}

	if last.Ratio != 255 {
		last = GradientItem{
			Ratio: 255,
			Color: last.Color,
		}
		items = append(items, last)
	}

	prevItem := items[0]
	for _, item := range items[1:] {
		prevItem = item
		prevColor := prevItem.Color
		currentColor := item.Color
		if interpolationMode == swfsubtypes.GradientInterpolationLinearRGB {
			prevColor = prevColor.ToLinearRGB()
			currentColor = currentColor.ToLinearRGB()
		}

		maxColorDistance := max(math.Abs(float64(prevColor.R)-float64(currentColor.R)), math.Abs(float64(prevColor.G)-float64(currentColor.G)), math.Abs(float64(prevColor.B)-float64(currentColor.B)), math.Abs(float64(prevColor.Alpha)-float64(currentColor.Alpha)))

		prevPosition := float64(prevItem.Ratio)
		currentPosition := float64(item.Ratio)
		distance := math.Abs(prevPosition - currentPosition)

		var partitions int
		if maxColorDistance < math.SmallestNonzeroFloat64 {
			partitions = 1
		} else if gradientSlices == GradientAutoSlices {
			partitions = max(1, int(math.Ceil(min(255/float64(len(items)+1), max(1, math.Ceil(maxColorDistance))))))
		} else {
			partitions = max(1, int(math.Ceil((distance/255)*float64(gradientSlices))))
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
				Start: fromPos / 255,
				End:   toPos / 255,
				Color: color,
			})
			fromPos = toPos
		}
	}
	return result
}
