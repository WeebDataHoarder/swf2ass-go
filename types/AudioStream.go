package types

import (
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
)

type AudioStream struct {
	Format     swftag.SoundFormat
	SampleRate int
	SampleSize int
	Channels   int

	SoundId   uint16
	SoundData []byte

	Start *int64

	Data []byte
}

func (s *AudioStream) AddStreamBlock(node *swftag.SoundStreamBlock) {
	switch s.Format {
	case swftag.SoundFormatMP3:
		s.addMP3Data(node.Data)
	}
}

func (s *AudioStream) addMP3Data(data []byte) {
	s.Data = append(s.Data, data[4:]...)
}

func AudioStreamFromSWF(soundRate swftag.SoundRate, soundSize uint8, isStereo bool, soundFormat swftag.SoundFormat) *AudioStream {
	rate := 0
	switch soundRate {
	case swftag.SoundRate5512Hz:
		rate = 5512
	case swftag.SoundRate11025Hz:
		rate = 11025
	case swftag.SoundRate22050Hz:
		rate = 22050
	case swftag.SoundRate44100Hz:
		rate = 44100
	}
	sampleSize := 8
	if soundSize == 1 {
		sampleSize = 16
	}

	channels := 1
	if isStereo {
		channels = 2
	}
	return &AudioStream{
		Format:     soundFormat,
		SampleRate: rate,
		SampleSize: sampleSize,
		Channels:   channels,
	}
}
