package ass

type Settings struct {
	SmoothTransitions bool
	DrawingScale      int64
}

var GlobalSettings = Settings{
	SmoothTransitions: false,
	DrawingScale:      DefaultDrawingScale,
}
