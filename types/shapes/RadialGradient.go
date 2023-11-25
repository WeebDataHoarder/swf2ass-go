package shapes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"slices"
)

type RadialGradient struct {
	Colors []GradientItem

	Transform         math.MatrixTransform
	SpreadMode        swfsubtypes.GradientSpreadMode
	InterpolationMode swfsubtypes.GradientInterpolationMode
}

func (g RadialGradient) GetSpreadMode() swfsubtypes.GradientSpreadMode {
	return g.SpreadMode
}

func (g RadialGradient) GetInterpolationMode() swfsubtypes.GradientInterpolationMode {
	return g.InterpolationMode
}

func (g RadialGradient) GetItems() []GradientItem {
	return g.Colors
}

func (g RadialGradient) GetInterpolatedDrawPaths(overlap, blur float64, gradientSlices int) DrawPathList {
	//items is max size 8 to 15 depending on SWF version
	size := GradientBounds.Width()

	//TODO spreadMode

	var paths DrawPathList
	for _, item := range LerpGradient(g, gradientSlices) {
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
	return paths.ApplyMatrixTransform(g.Transform, true)
}

func (g RadialGradient) GetMatrixTransform() math.MatrixTransform {
	return g.Transform
}

func (g RadialGradient) ApplyColorTransform(transform math.ColorTransform) Fillable {
	g2 := g
	g2.Colors = slices.Clone(g2.Colors)
	for i, g := range g2.Colors {
		g2.Colors[i] = GradientItem{
			Ratio: g.Ratio,
			Color: transform.ApplyToColor(g.Color),
		}
	}
	return &g2
}

func (g RadialGradient) Fill(shape *Shape) DrawPathList {
	return g.GetInterpolatedDrawPaths(settings.GlobalSettings.GradientOverlap, settings.GlobalSettings.GradientBlur, settings.GlobalSettings.GradientSlices).Fill(shape)
}

func RadialGradientFromSWF(records []swfsubtypes.GRADRECORD, transform types.MATRIX, spreadMode swfsubtypes.GradientSpreadMode, interpolationMode swfsubtypes.GradientInterpolationMode) DrawPathList {
	items := make([]GradientItem, 0, len(records))
	for _, r := range records {
		items = append(items, GradientItemFromSWF(r.Ratio, r.Color))
	}

	//TODO: interpolationMode, spreadMode

	return RadialGradient{
		Colors: items,
		//TODO: do we need to scale this to pixel world from twips?
		Transform:         math.MatrixTransformFromSWF(transform),
		SpreadMode:        spreadMode,
		InterpolationMode: interpolationMode,
	}.GetInterpolatedDrawPaths(settings.GlobalSettings.GradientOverlap, settings.GlobalSettings.GradientBlur, settings.GlobalSettings.GradientSlices)
}
