package ass

import (
	"fmt"
	"time"
)

const DecimalPrecision = 2

type Time struct {
	Total         uint64
	TotalAdjusted uint64

	Hours, Minutes, Seconds, Milliseconds uint64

	AdjustedMilliseconds         uint64
	AdjustedMillisecondPrecision uint64
	AdjustedMillisecondError     uint64
}

func NewTime(ms uint64, roundDown bool) (t Time) {
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

func (t Time) String() string {
	return fmt.Sprintf("%01d:%02d:%02d.%0*d", t.Hours, t.Minutes, t.Seconds, t.AdjustedMillisecondPrecision, t.AdjustedMilliseconds)
}

type EventTime struct {
	FrameDuration time.Duration

	Start, End           Time
	StartFrame, EndFrame uint64

	Duration uint64
}

func NewEventTime(startFrame, duration uint64, frameDuration time.Duration) (t EventTime) {
	t.FrameDuration = frameDuration
	t.StartFrame = startFrame
	t.Duration = duration
	t.Start = NewTime(uint64((time.Duration(t.StartFrame) * t.FrameDuration).Milliseconds()), true)
	t.EndFrame = t.StartFrame + t.Duration
	t.End = NewTime(uint64((time.Duration(t.EndFrame) * t.FrameDuration).Milliseconds()), false)

	return t
}

func (t EventTime) GetDurationFromStartOffset(frameOffset int64) time.Duration {
	if (int64(t.StartFrame) + frameOffset) > int64(t.EndFrame) {
		panic("out of bounds")
	}
	return t.FrameDuration*time.Duration(frameOffset) + time.Duration(t.Start.AdjustedMillisecondError)*time.Millisecond
}

func (t EventTime) GetDurationFromEndOffset(frameOffset int64) time.Duration {
	if frameOffset > int64(t.Duration) {
		panic("out of bounds")
	}
	return t.GetDurationFromStartOffset(int64(t.Duration) - frameOffset)
}

/*
func (t EventTime) Slice(frameOffset uint64, frameDuration uint64) EventTime {
	if (t.StartFrame + frameOffset + frameDuration) > t.EndFrame {
		panic("out of bounds")
	}
	return NewEventTime(t.StartFrame+frameOffset, t.FrameDuration)
}
*/

// StringToTimecode Emulates libass parsing
func StringToTimecode(p string) int64 {
	var h, m, s, ms int
	var tm int64
	res, _ := fmt.Sscanf(p, "%d:%d:%d.%d", &h, &m, &s, &ms)
	if res < 4 {
		return 0
	}
	tm = ((int64(h)*60+int64(m))*60+int64(s))*1000 + int64(ms)*10
	return tm
}
