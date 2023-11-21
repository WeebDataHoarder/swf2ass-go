package shapes

type DrawPath struct {
	Style    StyleRecord
	Commands *Shape
}

func DrawPathFill(record *FillStyleRecord, shape *Shape) DrawPath {
	return DrawPath{
		Style:    record,
		Commands: shape,
	}
}

func DrawPathStroke(record *LineStyleRecord, shape *Shape) DrawPath {
	return DrawPath{
		Style:    record,
		Commands: shape,
	}
}
