package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

func LinearGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) Gradient {
	items := make([]GradientItem, 0, len(records))
	for _, r := range records {
		items = append(items, GradientItemFromSWF(r.Ratio, r.Color))
	}

	//TODO: interpolationMode, spreadMode

	return Gradient{
		Records: items,
		//TODO: do we need to scale this to pixel world from twips?
		Transform:         math.MatrixTransformFromSWF(transform),
		SpreadMode:        spreadMode,
		InterpolationMode: interpolationMode,
		Interpolation: func(self Gradient, overlap, blur float64, gradientSlices int) DrawPathList {
			//items is max size 8 to 15 depending on SWF version
			size := GradientBounds.Width()

			//TODO spreadMode

			var paths DrawPathList
			for _, item := range InterpolateGradient(self, gradientSlices) {
				paths = append(paths, DrawPathFill(
					&FillStyleRecord{
						Fill: item.Color,
						Blur: blur,
					},
					NewShape(Rectangle[float64]{
						TopLeft:     math.NewVector2(GradientBounds.TopLeft.X+item.Start*size-overlap/2, GradientBounds.TopLeft.Y),
						BottomRight: math.NewVector2(GradientBounds.TopLeft.X+item.End*size+overlap/2, GradientBounds.BottomRight.Y),
					}.Draw()),
					nil, //TODO: clip here instead of outside
				).ApplyMatrixTransform(self.Transform, true))
			}
			return paths
		},
	}
}
