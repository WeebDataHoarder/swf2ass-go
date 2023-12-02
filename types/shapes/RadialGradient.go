package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	math2 "math"
)

func interpolateRadialGradient(self Gradient, overlap, blur float64, gradientSlices int, bb Rectangle[float64]) DrawPathList {
	//items is max size 8 to 15 depending on SWF version
	size := GradientBounds.Width()

	//TODO spreadMode

	items := InterpolateGradient(self, gradientSlices)
	for _, i := range items {
		if i.Color.Alpha != 255 {
			//transparency! remove overlaps
			blur = 0
			overlap = 0
			break
		}
	}

	radius := size / 2

	var maxRadius float64

	switch self.SpreadMode {
	case swfsubtypes.GradientSpreadPad:
		r := max(math2.Abs(bb.TopLeft.X), math2.Abs(bb.TopLeft.Y), math2.Abs(bb.BottomRight.X), math2.Abs(bb.BottomRight.Y))
		maxRadius = math.NewVector2(r, r).Length()
	}

	center := math.NewVector2[float64](0, 0)

	var paths DrawPathList
	for _, item := range items {
		//Create concentric circles to cut out a shape
		var shape Shape
		radiusStart := math.Lerp(0, radius, item.Start) - overlap/4
		radiusEnd := math.Lerp(0, radius, item.End) + overlap/4
		start := NewCircle(center, radiusStart).Draw()
		if radiusStart <= 0 {
			start = nil
		}
		end := NewCircle(center, radiusEnd).Draw()
		shape = append(shape, end...)
		shape = append(shape, start.Reverse()...)
		paths = append(paths, DrawPathFill(
			&FillStyleRecord{
				Fill: item.Color,
				Blur: blur,
			},
			shape,
		))
		if item.End == 1 {
			if maxRadius > radius {
				var shape Shape
				start := NewCircle(center, radiusEnd).Draw()
				if radiusEnd <= 0 {
					start = nil
				}
				end := NewCircle(center, maxRadius).Draw()
				shape = append(shape, end...)
				shape = append(shape, start.Reverse()...)
				paths = append(paths, DrawPathFill(
					&FillStyleRecord{
						Fill: item.Color,
						Blur: blur,
					},
					shape,
				))
			}
		}
	}
	return paths
}

func RadialGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) Gradient {
	items := make([]GradientItem, 0, len(records))
	for _, r := range records {
		items = append(items, GradientItemFromSWF(r.Ratio, r.Color))
	}

	//TODO: interpolationMode, spreadMode

	return Gradient{
		Records: items,
		//TODO: do we need to scale this to pixel world from twips?
		Transform:         math.MatrixTransformFromSWF(transform, 1),
		SpreadMode:        spreadMode,
		InterpolationMode: interpolationMode,
		Interpolation:     interpolateRadialGradient,
	}
}
