package ass

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"

type Settings struct {
	SmoothTransitions    bool
	DrawingScale         int64
	VideoScaleMultiplier int64
	VideoRateMultiplier  float64
	BakeTransforms       bool
	GradientSlices       int
}

var GlobalSettings = Settings{
	SmoothTransitions:    false,
	DrawingScale:         DefaultDrawingScale,
	VideoScaleMultiplier: 1,
	VideoRateMultiplier:  1,
	BakeTransforms:       false,
	GradientSlices:       types.GradientAutoSlices,
}
