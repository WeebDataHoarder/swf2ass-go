package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
	math2 "math"
	"strconv"
)

type DrawingTag interface {
	Tag
	ApplyMatrixTransform(transform math.MatrixTransform, applyTranslation bool) DrawingTag
	AsShape() shapes.Shape
	GetCommands(scale, precision int) string
}

type BaseDrawingTag shapes.Shape

func entryToPrecisionAndScaleTag(buf []byte, tag string, scale, precision int, vectors ...math.Vector2[float64]) []byte {
	if len(buf) > 0 {
		buf = append(buf, ' ')
	}
	if len(tag) > 0 {
		buf = append(buf, tag...)
		buf = append(buf, ' ')
	}
	for i, v := range vectors {
		if i > 0 {
			buf = append(buf, ' ')
		}
		buf = vectorToPrecisionAndScale(buf, scale, precision, v)
	}
	return buf
}

func vectorToPrecisionAndScale(buf []byte, scale, precision int, v math.Vector2[float64]) []byte {
	coords := v.Multiply(float64(scale))
	if precision == 0 {
		//fast path
		buf = strconv.AppendInt(buf, int64(math2.Round(coords.X)), 10)
		buf = append(buf, ' ')
		buf = strconv.AppendInt(buf, int64(math2.Round(coords.Y)), 10)
		return buf
	}
	buf = strconv.AppendFloat(buf, coords.X, 'f', precision, 64)
	buf = append(buf, ' ')
	buf = strconv.AppendFloat(buf, coords.Y, 'f', precision, 64)
	return buf
}

func (b *BaseDrawingTag) AsShape() shapes.Shape {
	return *(*shapes.Shape)(b)
}

func (b *BaseDrawingTag) GetCommands(scale, precision int) string {
	var lastEdge records.Record

	commands := make([]byte, 0, len(*b)*2*10)
	for _, edge := range *b {
		moveRecord, isMoveRecord := edge.(records.MoveRecord)
		if !isMoveRecord {
			if lastEdge == nil {
				commands = entryToPrecisionAndScaleTag(commands, "m", scale, precision, edge.GetStart())
			} else if !lastEdge.GetEnd().Equals(edge.GetStart()) {
				commands = entryToPrecisionAndScaleTag(commands, "m", scale, precision, edge.GetStart())
				lastEdge = nil
			}
		}

		if isMoveRecord {
			commands = entryToPrecisionAndScaleTag(commands, "m", scale, precision, moveRecord.To)
		} else if lineRecord, ok := edge.(records.LineRecord); ok {
			if _, ok = lastEdge.(records.LineRecord); ok {
				commands = entryToPrecisionAndScaleTag(commands, "", scale, precision, lineRecord.To)
			} else {
				commands = entryToPrecisionAndScaleTag(commands, "l", scale, precision, lineRecord.To)
			}
		} else if quadraticRecord, ok := edge.(records.QuadraticCurveRecord); ok {
			edge = records.CubicCurveFromQuadraticRecord(quadraticRecord)
		}

		if cubicRecord, ok := edge.(records.CubicCurveRecord); ok {
			if _, ok = lastEdge.(records.CubicCurveRecord); ok {
				commands = entryToPrecisionAndScaleTag(commands, "", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor)
			} else {
				commands = entryToPrecisionAndScaleTag(commands, "b", scale, precision, cubicRecord.Control1, cubicRecord.Control2, cubicRecord.Anchor)
			}
		} else if cubicSplineRecord, ok := edge.(records.CubicSplineCurveRecord); ok {
			_ = cubicSplineRecord
			panic("not implemented")
		}

		lastEdge = edge
	}

	/*if(!$this->shape->is_closed()){
	    $coords = $this->shape->start()->multiply($scale / Constants::TWIP_SIZE);
	    $commands[] = "n " . round($coords->x, $precision) . " " . round($coords->y, $precision);
	}*/

	return string(commands)
}
