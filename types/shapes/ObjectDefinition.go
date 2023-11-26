package shapes

type ObjectDefinition interface {
	GetObjectId() uint16
	GetSafeObject() ObjectDefinition
	GetShapeList(p ObjectProperties) DrawPathList
}

type ObjectProperties struct {
	Ratio float64
	// Data can be any value internal to the object itself
	Data any
}
