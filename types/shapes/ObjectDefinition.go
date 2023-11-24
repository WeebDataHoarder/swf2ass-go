package shapes

type ObjectDefinition interface {
	GetObjectId() uint16
	GetSafeObject() ObjectDefinition
	GetShapeList(ratio float64) DrawPathList
}
