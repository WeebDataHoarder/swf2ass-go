package ass

import (
	"fmt"
	swftypes "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"strings"
)

type DrawingTag interface {
	Tag
	ApplyMatrixTransform(transform types.MatrixTransform, applyTranslation bool) DrawingTag
	AsShape() *types.Shape
	GetCommands(scale, precision int64) []string
}

const DefaultDrawingScale = 6
const DefaultDrawingPrecision = 2

type BaseDrawingTag types.Shape

func twipEntryToPrecisionAndScaleTag(tag string, scale, precision int64, vectors ...types.Vector2[swftypes.Twip]) string {
	result := make([]string, 0, len(vectors)+1)
	if len(tag) > 0 {
		result = append(result, tag)
	}
	for _, v := range vectors {
		result = append(result, twipVectorToPrecisionAndScale(scale, precision, v))
	}
	return strings.Join(result, " ")
}

func twipVectorToPrecisionAndScale(scale, precision int64, v types.Vector2[swftypes.Twip]) string {
	coords := v.Multiply(swftypes.Twip(scale))
	return fmt.Sprintf("%.*f %.*f", precision, coords.X.Float64(), precision, coords.Y.Float64())
}

func (b *BaseDrawingTag) AsShape() *types.Shape {
	return (*types.Shape)(b)
}

func (b *BaseDrawingTag) GetCommands(scale, precision int64) []string {
	commands := make([]string, 0, len(b.Edges)*2)
	var lastEdge types.Record

	for _, edge := range b.Edges {
		moveRecord, isMoveRecord := edge.(*types.MoveRecord)
		if !isMoveRecord {
			if lastEdge == nil {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("m ", scale, precision, edge.GetStart()))
			} else if !lastEdge.GetEnd().Equals(edge.GetStart()) {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("m ", scale, precision, edge.GetStart()))
				lastEdge = nil
			}
		}

		if isMoveRecord {
			commands = append(commands, twipEntryToPrecisionAndScaleTag("m ", scale, precision, moveRecord.To))
		} else if lineRecord, ok := edge.(*types.LineRecord); ok {
			if _, ok = lastEdge.(*types.LineRecord); ok {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("", scale, precision, lineRecord.To))
			} else {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("l ", scale, precision, lineRecord.To))
			}
		} else if quadraticRecord, ok := edge.(*types.QuadraticCurveRecord); ok {
			edge = types.CubicCurveFromQuadraticRecord(quadraticRecord)
		}

		if cubicRecord, ok := edge.(*types.CubicCurveRecord); ok {
			if _, ok = lastEdge.(*types.CubicCurveRecord); ok {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor))
			} else {
				commands = append(commands, twipEntryToPrecisionAndScaleTag("b ", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor))
			}
		} else if cubicSplineRecord, ok := edge.(*types.CubicSplineCurveRecord); ok {
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
