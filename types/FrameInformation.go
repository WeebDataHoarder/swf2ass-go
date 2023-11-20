package types

import "time"

type FrameInformation struct {
	FrameNumber int64
	FrameOffset int64

	FrameRate float64

	Frame *ViewFrame
}

func (i FrameInformation) GetFrameNumber() int64 {
	return i.FrameNumber - i.FrameOffset
}

func (i FrameInformation) GetStartTime() time.Duration {
	//TODO: check this
	return (time.Duration(i.GetFrameNumber()) * time.Second) / time.Duration(float64(time.Second)*i.FrameRate)
}

func (i FrameInformation) GetEndTime() time.Duration {
	return i.GetStartTime() + i.GetFrameDuration()
}

func (i FrameInformation) GetFrameDuration() time.Duration {
	return time.Duration(float64(time.Second) * (1 / i.FrameRate))
}

func (i FrameInformation) Difference(o FrameInformation) FrameInformation {
	return FrameInformation{
		FrameNumber: i.GetFrameNumber() - o.GetFrameNumber(),
		FrameRate:   i.FrameRate,
		FrameOffset: 0,
		Frame:       nil,
	}
}
