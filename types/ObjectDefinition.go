package types

type ObjectDefinition interface {
	GetObjectId() uint16
	GetSafeObject() ObjectDefinition
	GetShapeList(ratio float64) DrawPathList
}

type MultiFrameObjectDefinition interface {
	ObjectDefinition
	NextFrame() *ViewFrame
}
