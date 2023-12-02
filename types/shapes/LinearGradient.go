package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf-go/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

func interpolateLinearGradient(self Gradient, overlap, blur float64, gradientSlices int, bb Rectangle[float64]) DrawPathList {
	//items is max size 8 to 15 depending on SWF version

	height0 := GradientBoundsMin.Float64()
	height1 := GradientBoundsMax.Float64()

	switch self.SpreadMode {
	case swfsubtypes.GradientSpreadPad:
		height0 = min(height0, bb.TopLeft.Y)
		height1 = max(height1, bb.BottomRight.Y)
	}

	topLeft0 := math.NewVector2(GradientBoundsMin.Float64(), height0)
	topLeft1 := math.NewVector2(GradientBoundsMax.Float64(), height0)
	bottomRight0 := math.NewVector2(GradientBoundsMin.Float64(), height1)
	bottomRight1 := math.NewVector2(GradientBoundsMax.Float64(), height1)

	//TODO: more spreadMode, generalize

	vOverlap := math.NewVector2(overlap, 0).Divide(2)

	items := InterpolateGradient(self, gradientSlices)
	for _, i := range items {
		if i.Color.Alpha != 255 {
			//transparency! remove overlaps
			blur = 0
			overlap = 0
			break
		}
	}

	var paths DrawPathList
	for _, item := range items {
		if item.Start == 0 {
			switch self.SpreadMode {
			case swfsubtypes.GradientSpreadPad:
				if bb.TopLeft.X < topLeft0.X {
					paths = append(paths, DrawPathFill(
						&FillStyleRecord{
							Fill: item.Color,
							Blur: blur,
						},
						Rectangle[float64]{
							TopLeft:     math.NewVector2(bb.TopLeft.X, height0),
							BottomRight: math.NewVector2(topLeft0.X, height1).AddVector(vOverlap),
						}.Draw(),
					))
				}
			}
		}
		paths = append(paths, DrawPathFill(
			&FillStyleRecord{
				Fill: item.Color,
				Blur: blur,
			},
			Rectangle[float64]{
				TopLeft:     math.LerpVector2(topLeft0, topLeft1, item.Start).SubVector(vOverlap),
				BottomRight: math.LerpVector2(bottomRight0, bottomRight1, item.End).AddVector(vOverlap),
			}.Draw(),
		))
		if item.End == 1 {
			switch self.SpreadMode {
			case swfsubtypes.GradientSpreadPad:
				if bb.BottomRight.X > bottomRight1.X {
					paths = append(paths, DrawPathFill(
						&FillStyleRecord{
							Fill: item.Color,
							Blur: blur,
						},
						Rectangle[float64]{
							TopLeft:     math.NewVector2(topLeft1.X, height0).SubVector(vOverlap),
							BottomRight: math.NewVector2(bb.BottomRight.X, height1),
						}.Draw(),
					))
				}
			}
		}
	}

	return paths
}

func LinearGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) Gradient {
	items := make([]GradientItem, 0, len(records))
	for _, r := range records {
		items = append(items, GradientItemFromSWF(r.Ratio, r.Color))
	}

	//TODO: interpolationMode, spreadMode

	return Gradient{
		Records:           items,
		Transform:         math.MatrixTransformFromSWF(transform, 1),
		SpreadMode:        spreadMode,
		InterpolationMode: interpolationMode,
		Interpolation:     interpolateLinearGradient,
	}
}
