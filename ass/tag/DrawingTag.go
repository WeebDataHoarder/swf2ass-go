package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	"strings"
)

type DrawingTag interface {
	Tag
	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) DrawingTag
	AsShape() *shapes.Shape
	GetCommands(scale, precision int64) []string
}

type BaseDrawingTag shapes.Shape

func entryToPrecisionAndScaleTag(tag string, scale, precision int64, vectors ...math.Vector2[float64]) string {
	result := make([]string, 0, len(vectors)+1)
	if len(tag) > 0 {
		result = append(result, tag)
	}
	for _, v := range vectors {
		result = append(result, vectorToPrecisionAndScale(scale, precision, v))
	}
	return strings.Join(result, " ")
}

func vectorToPrecisionAndScale(scale, precision int64, v math.Vector2[float64]) string {
	coords := v.Multiply(float64(scale))
	return fmt.Sprintf("%.*f %.*f", precision, coords.X, precision, coords.Y)
}

func (b *BaseDrawingTag) AsShape() *shapes.Shape {
	return (*shapes.Shape)(b)
}

func (b *BaseDrawingTag) GetCommands(scale, precision int64) []string {
	commands := make([]string, 0, len(b.Edges)*2)
	var lastEdge records.Record

	for _, edge := range b.Edges {
		moveRecord, isMoveRecord := edge.(*records.MoveRecord)
		if !isMoveRecord {
			if lastEdge == nil {
				commands = append(commands, entryToPrecisionAndScaleTag("m", scale, precision, edge.GetStart()))
			} else if !lastEdge.GetEnd().Equals(edge.GetStart()) {
				commands = append(commands, entryToPrecisionAndScaleTag("m", scale, precision, edge.GetStart()))
				lastEdge = nil
			}
		}

		if isMoveRecord {
			commands = append(commands, entryToPrecisionAndScaleTag("m", scale, precision, moveRecord.To))
		} else if lineRecord, ok := edge.(*records.LineRecord); ok {
			if _, ok = lastEdge.(*records.LineRecord); ok {
				commands = append(commands, entryToPrecisionAndScaleTag("", scale, precision, lineRecord.To))
			} else {
				commands = append(commands, entryToPrecisionAndScaleTag("l", scale, precision, lineRecord.To))
			}
		} else if quadraticRecord, ok := edge.(*records.QuadraticCurveRecord); ok {
			edge = records.CubicCurveFromQuadraticRecord(quadraticRecord)
		}

		if cubicRecord, ok := edge.(*records.CubicCurveRecord); ok {
			if _, ok = lastEdge.(*records.CubicCurveRecord); ok {
				commands = append(commands, entryToPrecisionAndScaleTag("", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor))
			} else {
				commands = append(commands, entryToPrecisionAndScaleTag("b", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor))
			}
		} else if cubicSplineRecord, ok := edge.(*records.CubicSplineCurveRecord); ok {
			_ = cubicSplineRecord
			panic("not implemented")
		}

		lastEdge = edge
	}

	/*if(!$this->shape->is_closed()){
	    $coords = $this->shape->start()->multiply($scale / Constants::TWIP_SIZE);
	    $commands[] = "n " . round($coords->x, $precision) . " " . round($coords->y, $precision);
	}*/

	return commands
}
