package records

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"

func LerpRecord(start, end Record, ratio float64) Record {
	if start.SameType(end) {
		switch s := start.(type) {
		case LineRecord:
			return LineRecord{
				To:    math.LerpVector2(s.To, end.(LineRecord).To, ratio),
				Start: math.LerpVector2(s.Start, end.(LineRecord).Start, ratio),
			}
		case MoveRecord:
			return MoveRecord{
				To:    math.LerpVector2(s.To, end.(MoveRecord).To, ratio),
				Start: math.LerpVector2(s.Start, end.(MoveRecord).Start, ratio),
			}
		case QuadraticCurveRecord:
			return QuadraticCurveRecord{
				Control: math.LerpVector2(s.Control, end.(QuadraticCurveRecord).Control, ratio),
				Anchor:  math.LerpVector2(s.Anchor, end.(QuadraticCurveRecord).Anchor, ratio),
				Start:   math.LerpVector2(s.Start, end.(QuadraticCurveRecord).Start, ratio),
			}
		case CubicCurveRecord:
			return CubicCurveRecord{
				Control1: math.LerpVector2(s.Control1, end.(CubicCurveRecord).Control1, ratio),
				Control2: math.LerpVector2(s.Control2, end.(CubicCurveRecord).Control2, ratio),
				Anchor:   math.LerpVector2(s.Anchor, end.(CubicCurveRecord).Anchor, ratio),
				Start:    math.LerpVector2(s.Start, end.(CubicCurveRecord).Start, ratio),
			}
		default:
			panic("not supported")
		}
	} else {
		startLine, startLineOk := start.(LineRecord)
		startQuadratic, startQuadraticOk := start.(QuadraticCurveRecord)
		endLine, endLineOk := end.(LineRecord)
		endQuadratic, endQuadraticOk := end.(QuadraticCurveRecord)

		if startLineOk && endQuadraticOk {
			startQuadratic = QuadraticCurveFromLineRecord(startLine)
			return QuadraticCurveRecord{
				Control: math.LerpVector2(startQuadratic.Control, endQuadratic.Control, ratio),
				Anchor:  math.LerpVector2(startQuadratic.Anchor, endQuadratic.Anchor, ratio),
				Start:   math.LerpVector2(startQuadratic.Start, endQuadratic.Start, ratio),
			}
		} else if startQuadraticOk && endLineOk {
			endQuadratic = QuadraticCurveFromLineRecord(endLine)
			return QuadraticCurveRecord{
				Control: math.LerpVector2(startQuadratic.Control, endQuadratic.Control, ratio),
				Anchor:  math.LerpVector2(startQuadratic.Anchor, endQuadratic.Anchor, ratio),
				Start:   math.LerpVector2(startQuadratic.Start, endQuadratic.Start, ratio),
			}
		} else {
			panic("not supported")
		}
	}
}
