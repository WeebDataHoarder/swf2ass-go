package types

import (
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/shapes"
)

const BackgroundObjectId = 0

const BackgroundObjectDepth = 0

type SWFProcessor struct {
	SWFTreeProcessor

	Background         *shapes.FillStyleRecord
	ViewPort           shapes.Rectangle[types.Twip]
	FrameRate          float64
	ExpectedFrameCount int64
}

func NewSWFProcessor(tags []swftag.Tag, viewPort shapes.Rectangle[types.Twip], frameRate float64, frameCount int64) *SWFProcessor {
	p := &SWFProcessor{
		SWFTreeProcessor: *NewSWFTreeProcessor(0, tags, make(ObjectCollection)),
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
		//TODO: handle sound
	}
	return p.process(actions)
}

func (p *SWFProcessor) NextFrameOutput() *FrameInformation {
	frame := p.NextFrame()
	if frame == nil {
		return nil
	}
	/*
		if(!$this->isPlaying() and ($this->audio === null or $this->audio->getStartFrame() === null) or $this->getFrame() === 1){ //Force play till finding audio, or first frame is 0
			$this->playing = true;
		}
	*/
	if !p.Playing && (p.Frame == 1) {
		p.Playing = true
	}

	//TODO: actions?

	frame.AddChild(BackgroundObjectDepth, NewViewFrame(BackgroundObjectId, &shapes.DrawPathList{shapes.DrawPathFill(p.Background, shapes.NewShape(p.ViewPort.Draw()))}))
	return &FrameInformation{
		FrameNumber: p.Frame - 1,
		FrameRate:   p.FrameRate,
		Frame:       frame,
	}
}
