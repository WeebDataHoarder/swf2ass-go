package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
)

func RadialGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) Gradient {
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
				//Create concentric circles to cut out a shape
				var shape Shape
				radiusStart := (item.Start*size)/2 - overlap/4
				radiusEnd := (item.End*size)/2 + overlap/4
				start := NewCircle(math.NewVector2[float64](0, 0), radiusStart).Draw()
				if radiusStart <= 0 {
					start = nil
				}
				end := NewCircle(math.NewVector2[float64](0, 0), radiusEnd).Draw()
				shape.Edges = append(shape.Edges, end...)
				shape.Edges = append(shape.Edges, NewShape(start).Reverse().Edges...)
				paths = append(paths, DrawPathFill(
					&FillStyleRecord{
						Fill: item.Color,
						Blur: blur,
					},
					&shape,
					nil, //TODO: clip here instead of outside
				))
			}
			return paths.ApplyMatrixTransform(self.Transform, true)
		},
	}
}
