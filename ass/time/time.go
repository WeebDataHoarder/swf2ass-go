package time

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

const DecimalPrecision = 2

type Time struct {
	Total         int64
	TotalAdjusted int64

	Hours, Minutes, Seconds, Milliseconds int64

	AdjustedMilliseconds         int64
	AdjustedMillisecondPrecision int64
	AdjustedMillisecondError     int64
}

func NewTime(ms int64, roundDown bool) (t Time) {
	t.Total = ms
	t.AdjustedMillisecondPrecision = DecimalPrecision

	const HoursMilliseconds = 1000 * 3600
	const MinutesMilliseconds = 1000 * 60

	t.Hours = ms / HoursMilliseconds
	ms -= t.Hours * HoursMilliseconds
	t.Minutes = ms / MinutesMilliseconds
	ms -= t.Minutes * MinutesMilliseconds

	t.Seconds = ms / 1000
	t.Milliseconds = ms - t.Seconds*1000

	msAdjustement := 10 * (3 - t.AdjustedMillisecondPrecision)
	t.AdjustedMilliseconds = t.Milliseconds / msAdjustement

	t.AdjustedMillisecondError = t.Milliseconds - t.AdjustedMilliseconds*msAdjustement

	if !roundDown && t.AdjustedMillisecondError > 0 {
		t.AdjustedMilliseconds++
		t.AdjustedMillisecondError -= msAdjustement
	}
	t.TotalAdjusted = t.Total + t.AdjustedMillisecondError

	return t
}

func (t Time) Duration() time.Duration {
	//TODO
	return 0
}

func (t Time) String() string {
	return fmt.Sprintf("%01d:%02d:%02d.%0*d", t.Hours, t.Minutes, t.Seconds, t.AdjustedMillisecondPrecision, t.AdjustedMilliseconds)
}

func FromString(p string) (t Time, err error) {
	t.AdjustedMillisecondPrecision = DecimalPrecision

	res, _ := fmt.Sscanf(p, "%d:%d:%d.%d", &t.Hours, &t.Minutes, &t.Seconds, &t.Milliseconds)
	if res < 4 {
		return t, errors.New("not parsed properly")
	}

	//TODO
	return t, nil
}

func (t Time) AdjustWithError(ms, adjMsErr int64) Time {
	t.AdjustedMillisecondPrecision = DecimalPrecision
	t.AdjustedMillisecondError = adjMsErr
	if t.AdjustedMillisecondError != 0 {
		t.Milliseconds = ms
	}
	//TODO
	return t
}

type EventTime struct {
	FrameDuration time.Duration

	Start, End           Time
	StartFrame, EndFrame int64

	Duration int64
}

func NewEventTime(startFrame, duration int64, frameDuration time.Duration) (t EventTime) {
	t.FrameDuration = frameDuration
	t.StartFrame = startFrame
	t.Duration = duration
	t.Start = NewTime((time.Duration(t.StartFrame) * t.FrameDuration).Milliseconds(), true)
	t.EndFrame = t.StartFrame + t.Duration
	t.End = NewTime((time.Duration(t.EndFrame) * t.FrameDuration).Milliseconds(), false)

	return t
}

func (t EventTime) GetDurationFromStartOffset(frameOffset int64) time.Duration {
	if (t.StartFrame + frameOffset) > t.EndFrame {
		panic("out of bounds")
	}
	return t.FrameDuration*time.Duration(frameOffset) + time.Duration(t.Start.AdjustedMillisecondError)*time.Millisecond
}

func (t EventTime) GetDurationFromEndOffset(frameOffset int64) time.Duration {
	if frameOffset > t.Duration {
		panic("out of bounds")
	}
	return t.GetDurationFromStartOffset(int64(t.Duration) - frameOffset)
}

func (t EventTime) Encode() string {
	if t.Start.AdjustedMillisecondError != 0 || t.End.AdjustedMillisecondError != 0 {
		//Adjust frame precision exactly to frame boundaries. This is necessary due to low ASS timing precision
		//TODO: Maybe use fade?
		frameStartTime := t.GetDurationFromStartOffset(0).Milliseconds()
		frameEndTime := t.GetDurationFromEndOffset(0).Milliseconds()
		//TODO: maybe needs to be -1?
		return fmt.Sprintf(
			"{\\fade(255,0,255,%d,%d,%d,%d)\\err(%d~%d,%d~%d)}", frameStartTime, frameStartTime, frameEndTime, frameEndTime, t.Start.Milliseconds, t.Start.AdjustedMillisecondError, t.End.Milliseconds, t.End.AdjustedMillisecondError)
	}
	return ""
}

var eventTimeRegexp = regexp.MustCompile(`^\{\\fade\(255,0,255,(?P<FrameStartTimeMs>\d+),(?P<FrameStartTimeMs2>\d+),(?P<FrameEndTimeMs>\d+),(?P<FrameEndTimeMs2>\d+)\)\\err\((?P<StartMs>\d+)~(?P<StartAdjustedErrMs>\d+),(?P<EndMs>\d+)~(?P<EndAdjustedErrMs>\d+)\)}`)

func EventLineFromText(start, end Time, text string) (t EventTime) {
	var frameStartTimeMs, frameStartTimeMs2, frameEndTimeMs, frameEndTimeMs2 int64
	var startMs, startAdjustedErrMs, endMs, endAdjustedMs int64

	matches := eventTimeRegexp.FindStringSubmatch(text)
	if matches != nil {
		//not exact

		for i, name := range eventTimeRegexp.SubexpNames() {
			if name == "" {
				continue
			}
			val, err := strconv.ParseInt(matches[i], 10, 0)
			if err != nil {
				panic(err)
			}
			switch name {
			case "FrameStartTimeMs":
				frameStartTimeMs = val
			case "FrameStartTimeMs2":
				frameStartTimeMs2 = val
			case "FrameEndTimeMs":
				frameEndTimeMs = val
			case "FrameEndTimeMs2":
				frameEndTimeMs2 = val
			case "StartMs":
				startMs = val
			case "StartAdjustedErrMs":
				startAdjustedErrMs = val
			case "EndMs":
				endMs = val
			case "EndAdjustedErrMs":
				endAdjustedMs = val

			default:
				panic("not implemented")

			}
		}
	}

	if frameStartTimeMs != frameStartTimeMs2 {
		panic("frame start match")
	}

	if frameEndTimeMs != frameEndTimeMs2 {
		panic("frame start match")
	}

	t.Start = start.AdjustWithError(startMs, startAdjustedErrMs)
	t.End = end.AdjustWithError(endMs, endAdjustedMs)
	//TODO
	return t
}
