package settings

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

type Settings struct {
	// ASSDrawingScale Scale that ASS drawing override tags will use. Coordinates will be multiplied by 2^(ASSDrawingScale-1) to enhance precision
	ASSDrawingScale int64

	// ASSDrawingPrecision Number of decimals that ASS drawing override tags will produce
	// Note that at high ASSDrawingScale >= 5 this will be brought down to 0 regardless
	ASSDrawingPrecision int64

	// VideoRateMultiplier Adjusts the viewport scale. All operations and transforms will be adjusted accordingly
	// For example, VideoScaleMultiplier = 2 will make a 640x480 viewport become 1280x960
	VideoScaleMultiplier float64

	// VideoRateMultiplier Adjusts framerate multiplier. Leave at 1 unless you know what you are doing.
	VideoRateMultiplier float64

	// BakeClips Clip any shape that has a clip directly, instead of emitting \clip tags
	BakeClips bool

	// BakeMatrixTransforms Transforms the shapes directly instead of writing ASS override tags
	// Reduces compressibility, but not affect positioning. \pos tags will still be emitted
	// Enabling this is very expensive on players, also increases output size.
	BakeMatrixTransforms bool

	// SmoothTransitions Attempt to merge multiple fixed transitions into a single long one if they happen at constant rate
	// \move and \t tags will be emitted
	// Note currently only linear changes are tracked, TODO: track changes in non-linear supported by ASS override tags
	SmoothTransitions bool

	// GradientSlices Number of slices each gradient will get for each step when rendering them.
	// It is recommended to leave at shapes.GradientAutoSlices as that will automatically pick slices based on color differences across steps.
	GradientSlices int

	// GradientOverlap Overlap between slices, in pixel units, before transformation
	GradientOverlap float64

	// GradientBlur Amount of blur to apply to gradients
	GradientBlur float64
}

const DefaultASSDrawingScale = 6
const DefaultASSDrawingPrecision = 2

var GlobalSettings = Settings{
	ASSDrawingScale:      DefaultASSDrawingScale,
	ASSDrawingPrecision:  DefaultASSDrawingPrecision,
	VideoScaleMultiplier: 1,
	VideoRateMultiplier:  1,

	BakeMatrixTransforms: true,

	BakeClips: false,

	GradientSlices:  shapes.GradientAutoSlices,
	GradientOverlap: 2,
	GradientBlur:    0.1,
}
