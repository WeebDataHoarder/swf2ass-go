package tag

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/settings"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
	"strconv"
)

type PositionTag struct {
	From, To   math2.Vector2[float64]
	Start, End int64
}

func (t *PositionTag) TransitionMatrixTransform(event Event, transform math2.MatrixTransform) PositioningTag {
	translation := transform.GetTranslation()

	frame := event.GetEnd() - event.GetStart()

	isInitialState := t.Start == frame && t.End == frame
	isMovingState := t.Start < frame && t.End == frame
	isMovedState := t.Start < frame && t.End < frame

	if t.To.Equals(translation) {
		if isInitialState {
			return &PositionTag{
				From:  t.From,
				To:    t.To,
				Start: t.Start + 1, //TODO: should this be +1?
				End:   t.End + 1,
			}
		} else if isMovingState || isMovedState {
			return &PositionTag{
				From:  t.From,
				To:    t.To,
				Start: t.Start,
				End:   t.End,
			}
		} else {
			panic("logic error")
		}
	}

	if isInitialState { //Always allow initial move
		return &PositionTag{
			From:  t.From,
			To:    translation,
			Start: t.Start,
			End:   t.End + 1,
		}
	} else if isMovingState {
		duration := t.End - t.Start + 1

		direction := t.To.SubVector(t.From).Normalize()
		//TODO: maybe use larger epsilon?
		if math.Abs(direction.Dot(translation.Normalize())-1) <= math.SmallestNonzeroFloat64 { //Same direction, extend
			length := t.To.SubVector(t.From).Divide(float64(duration)).SquaredLength()
			length2 := translation.SubVector(t.To).SquaredLength()

			if math.Abs(length-length2) <= math.SmallestNonzeroFloat64 { //same length
				return &PositionTag{
					From:  t.From,
					To:    translation,
					Start: t.Start,
					End:   t.End + 1,
				}
			}
		}
		return nil
	} else if isMovedState {
		return nil
	} else {
		panic("logic error")
	}
}

func (t *PositionTag) Encode(event time.EventTime) string {
	hasMoved := t.Start != t.End

	shift := t.End - t.Start

	if hasMoved {
		var start, end int64
		if shift > 1 || settings.GlobalSettings.ASSSmoothTransitions {
			start = event.GetDurationFromStartOffset(t.Start - 1).Milliseconds()
			end = event.GetDurationFromStartOffset(t.End).Milliseconds()
		} else {
			start = event.GetDurationFromStartOffset(t.Start).Milliseconds() - 1
			end = event.GetDurationFromStartOffset(t.Start).Milliseconds()
		}
		return fmt.Sprintf("\\move(%s,%s,%s,%s,%d,%d)", strconv.FormatFloat(t.From.X, 'f', -1, 64), strconv.FormatFloat(t.From.Y, 'f', -1, 64), strconv.FormatFloat(t.To.X, 'f', -1, 64), strconv.FormatFloat(t.To.Y, 'f', -1, 64), start, end)
	}

	return fmt.Sprintf("\\pos(%s,%s)", strconv.FormatFloat(t.From.X, 'f', -1, 64), strconv.FormatFloat(t.From.Y, 'f', -1, 64))
}

func (t *PositionTag) Equals(tag Tag) bool {
	if o, ok := tag.(*PositionTag); ok {
		return *t == *o
	}
	return false
}

func (t *PositionTag) FromMatrixTransform(transform math2.MatrixTransform) PositioningTag {
	translation := transform.GetTranslation()
	t.From = translation
	t.To = translation
	t.Start = 1
	t.End = 1
	return t
}
