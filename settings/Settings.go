package settings

type Settings struct {
	// ASSDrawingScale Scale that ASS drawing override tags will use. Coordinates will be multiplied by 2^(ASSDrawingScale-1) to enhance precision
	ASSDrawingScale int

	// ASSDrawingPrecision Number of decimals that ASS drawing override tags will produce
	// Note that at high ASSDrawingScale >= 5 this will be brought down to 0 regardless
	ASSDrawingPrecision int

	// ASSPreciseTiming Enables precise timing via \fade transitions for each event.
	// libass is only precise to 1/100th of a second, causing issues with higher framerate content
	// Additionally extra \err annotations will be placed to note the current error adjusted
	ASSPreciseTiming bool

	// ASSSmoothTransitions Attempt to merge multiple fixed transitions into a single long one if they happen at constant rate
	// \move and \t tags will be emitted
	// Note currently only linear changes are tracked, TODO: track changes in non-linear supported by ASS override tags
	ASSSmoothTransitions bool

	// ASSBakeClips Clip any shape that has a clip directly, instead of emitting \clip tags
	ASSBakeClips bool

	// ASSBakeMatrixTransforms Transforms the shapes directly instead of writing ASS override tags
	// Reduces compressibility, but not affect positioning. \pos tags will still be emitted
	// Enabling this is very expensive on players, also increases output size.
	ASSBakeMatrixTransforms bool

	// VideoRateMultiplier Adjusts the viewport scale. All operations and transforms will be adjusted accordingly
	// For example, VideoScaleMultiplier = 2 will make a 640x480 viewport become 1280x960
	VideoScaleMultiplier float64

	// VideoRateMultiplier Adjusts framerate multiplier. Leave at 1 unless you know what you are doing.
	VideoRateMultiplier float64

	// KeyFrameInterval Flush a shape at least every N frames if it's doing transitions. Set to 0 to disable.
	KeyFrameInterval int

	// GradientSlices Number of slices each gradient will get for each step when rendering them.
	// It is recommended to leave at shapes.GradientAutoSlices as that will automatically pick slices based on color differences across steps.
	GradientSlices int

	// GradientOverlap Overlap between slices, in pixel units, before transformation
	GradientOverlap float64

	// GradientBlur Amount of blur to apply to gradients
	GradientBlur float64

	// BitmapPaletteSize Number of colors bitmap shapes will be paletted into.
	// Transparency is not included and will be handled separately. Set to 0 to disable.
	BitmapPaletteSize int
	// BitmapMaxDimension Highest dimension, width or height, that a bitmap can be sized
	// If higher, the bitmap will be resized accordingly. Set to 0 to disable.
	BitmapMaxDimension int
}

const GradientAutoSlices = -1
const DefaultASSDrawingScale = 6
const DefaultASSDrawingPrecision = 2

var GlobalSettings = Settings{
	ASSDrawingScale:         DefaultASSDrawingScale,
	ASSDrawingPrecision:     DefaultASSDrawingPrecision,
	ASSPreciseTiming:        true,
	ASSSmoothTransitions:    false,
	ASSBakeClips:            false,
	ASSBakeMatrixTransforms: false,

	VideoScaleMultiplier: 1,
	VideoRateMultiplier:  1,

	KeyFrameInterval: 0,

	GradientSlices:  GradientAutoSlices,
	GradientOverlap: 2,
	GradientBlur:    0.1,

	BitmapPaletteSize:  32,
	BitmapMaxDimension: 256,
}
