package types

import (
	swfsubtypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
)

const GradientAutoSlices = -1

type Gradient interface {
	GetSpreadMode() swfsubtypes.GradientSpreadMode
	GetInterpolationMode() swfsubtypes.GradientInterpolationMode
	GetItems() []GradientItem
	GetInterpolatedDrawPaths(overlap int, slices int) DrawPathList
	GetMatrixTransform() MatrixTransform
	ApplyColorTransform(transform ColorTransform) Gradient
}

type GradientItem struct {
	Ratio uint8
	Color Color
}
