package shapes

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"slices"
)

type LinearGradient struct {
	Colors []GradientItem

	Transform         math.MatrixTransform
	SpreadMode        swfsubtypes.GradientSpreadMode
	InterpolationMode swfsubtypes.GradientInterpolationMode
}

func (g *LinearGradient) GetSpreadMode() swfsubtypes.GradientSpreadMode {
	return g.SpreadMode
}

func (g *LinearGradient) GetInterpolationMode() swfsubtypes.GradientInterpolationMode {
	return g.InterpolationMode
}

func (g *LinearGradient) GetItems() []GradientItem {
	return g.Colors
}

func (g *LinearGradient) GetInterpolatedDrawPaths(overlap, blur float64, gradientSlices int) DrawPathList {
	//items is max size 8 to 15 depending on SWF version
	size := GradientBounds.Width()

	//TODO spreadMode

	var paths DrawPathList
	for _, item := range LerpGradient(g, gradientSlices) {
		paths = append(paths, DrawPathFill(
			&FillStyleRecord{
				Fill: item.Color,
				Blur: blur,
			},
			NewShape(Rectangle[float64]{
				TopLeft:     math.NewVector2(GradientBounds.TopLeft.X+item.Start*size-overlap/2, GradientBounds.TopLeft.Y),
				BottomRight: math.NewVector2(GradientBounds.TopLeft.X+item.End*size+overlap/2, GradientBounds.BottomRight.Y),
			}.Draw()),
		))
	}
	return paths
}

func (g *LinearGradient) GetMatrixTransform() math.MatrixTransform {
	return g.Transform
}

func (g *LinearGradient) ApplyColorTransform(transform math.ColorTransform) Gradient {
	g2 := *g
	g2.Colors = slices.Clone(g2.Colors)
	for i, g := range g2.Colors {
		g2.Colors[i] = GradientItem{
			Ratio: g.Ratio,
			Color: transform.ApplyToColor(g.Color),
		}
	}
	return &g2
}

func LinearGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) *LinearGradient {
	items := make([]GradientItem, 0, len(records))
	for _, r := range records {
		items = append(items, GradientItemFromSWF(r.Ratio, r.Color))
	}

	//TODO: interpolationMode, spreadMode

	return &LinearGradient{
		Colors:            items,
		Transform:         math.MatrixTransformFromSWF(transform),
		SpreadMode:        spreadMode,
		InterpolationMode: interpolationMode,
	}
}
