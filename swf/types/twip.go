package types

const TwipFactor = 20

type Twip int64

func (t Twip) IntInt() (pixel int64, subPixel int64) {
	return int64(t / TwipFactor), int64(t % TwipFactor)
}

func (t Twip) Float64() float64 {
	return float64(t) / TwipFactor
}

func (t Twip) Float32() float32 {
	return float32(t) / TwipFactor
}
