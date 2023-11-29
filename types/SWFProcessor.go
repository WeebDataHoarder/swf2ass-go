package types

import (
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

const BackgroundObjectId = 0

const BackgroundObjectDepth = 0

type SWFProcessor struct {
	SWFTreeProcessor

	Background         *shapes.FillStyleRecord
	ViewPort           shapes.Rectangle[float64]
	FrameRate          float64
	ExpectedFrameCount int64

	Audio *AudioStream
}

func NewSWFProcessor(tags []swftag.Tag, viewPort shapes.Rectangle[float64], frameRate float64, frameCount int64, version uint8) *SWFProcessor {
	p := &SWFProcessor{
		SWFTreeProcessor: *NewSWFTreeProcessor(0, tags, make(shapes.ObjectCollection), version),
		Background: &shapes.FillStyleRecord{
			Fill: math.Color{
				R:     255,
				G:     255,
				B:     255,
				Alpha: 255,
			},
			Border: nil,
		},
		ViewPort:           viewPort,
		FrameRate:          frameRate,
		ExpectedFrameCount: frameCount,
	}
	p.processFunc = p.subProcess
	return p
}

func (p *SWFProcessor) subProcess(actions ActionList) (tag swftag.Tag, newActions ActionList) {
	tag = p.Current()
	if tag == nil {
		return nil, nil
	}
	switch node := tag.(type) {
	case *swftag.SetBackgroundColor:
		p.Background = &shapes.FillStyleRecord{
			Fill: math.Color{
				R:     node.BackgroundColor.R(),
				G:     node.BackgroundColor.G(),
				B:     node.BackgroundColor.B(),
				Alpha: node.BackgroundColor.A(),
			},
			Border: nil,
		}
	case *swftag.SoundStreamHead:
		if p.Loops > 0 {
			break
		}
		//fixes swf without actual audio but a head
		if !func() bool {
			for _, t := range p.Tags {
				if _, ok := t.(*swftag.SoundStreamBlock); ok {
					return true
				}
			}
			return false
		}() {
			break
		}
		p.Audio = AudioStreamFromSWF(node.StreamSoundRate, node.StreamSoundSize, node.StreamIsStereo, swftag.SoundFormat(node.StreamSoundCompression))
	case *swftag.SoundStreamHead2:
		if p.Loops > 0 {
			break
		}
		//fixes swf without actual audio but a head
		if !func() bool {
			for _, t := range p.Tags {
				if _, ok := t.(*swftag.SoundStreamBlock); ok {
					return true
				}
			}
			return false
		}() {
			break
		}
		p.Audio = AudioStreamFromSWF(node.StreamSoundRate, node.StreamSoundSize, node.StreamIsStereo, node.StreamSoundFormat)
	case *swftag.SoundStreamBlock:
		if p.Loops > 0 {
			break
		}
		if p.Audio != nil {
			if p.Audio.Start == nil {
				f := p.Frame
				p.Audio.Start = &f
			}
			p.Audio.AddStreamBlock(node)
		}
	case *swftag.DefineSound:
		if p.Loops > 0 {
			break
		}
		if p.Audio != nil {
			break
		}
		p.Audio = AudioStreamFromSWF(node.SoundRate, node.SoundSize, node.IsStereo, node.SoundFormat)
		p.Audio.SoundId = node.SoundId
		p.Audio.SoundData = node.SoundData
	case *swftag.StartSound:
		if p.Loops > 0 {
			break
		}
		if p.Audio != nil && p.Audio.SoundId == node.SoundId {
			if p.Audio.Start == nil {
				f := p.Frame
				p.Audio.Start = &f
			}
			p.Audio.Data = p.Audio.SoundData
		}
	case *swftag.StartSound2:
		if p.Loops > 0 {
			break
		}
		if p.Audio != nil && p.Audio.SoundId == node.SoundId {
			if p.Audio.Start == nil {
				f := p.Frame
				p.Audio.Start = &f
			}
			p.Audio.Data = p.Audio.SoundData
		}
	}
	return p.process(actions)
}

func (p *SWFProcessor) NextFrameOutput() *FrameInformation {
	frame := p.NextFrame(false)
	if frame == nil {
		return nil
	}
	// Stop looping main video
	if p.Loops > 0 {
		return nil
	}

	if !p.Playing && (p.Audio == nil || p.Audio.Start == nil) || p.Frame == 1 { //Force play till finding audio, or first frame is 0
		p.Playing = true
	}

	//TODO: actions?

	frame.AddChild(BackgroundObjectDepth, NewViewFrame(BackgroundObjectId, &shapes.DrawPathList{shapes.DrawPathFill(p.Background, p.ViewPort.Draw())}))
	return &FrameInformation{
		FrameNumber: p.Frame - 1,
		FrameRate:   p.FrameRate,
		Frame:       frame,
	}
}
