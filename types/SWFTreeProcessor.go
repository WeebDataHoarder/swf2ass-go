package types

import (
	swftag "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
)

type SWFTreeProcessor struct {
	Layout *ViewLayout

	Objects ObjectCollection

	Tags []swftag.Tag

	Index int

	Frame int64

	LastFrame *ViewFrame
	Playing   bool
	Loops     int
}

func NewSWFTreeProcessor(objectId uint16, tags []swftag.Tag, objects ObjectCollection) *SWFTreeProcessor {
	return &SWFTreeProcessor{
		Objects: objects,
		Frame:   0,
		Tags:    tags,
		Layout:  NewViewLayout(objectId, nil, nil),
	}
}

func (p *SWFTreeProcessor) Next() {
	p.Index++
}

func (p *SWFTreeProcessor) Current() swftag.Tag {
	if len(p.Tags) > p.Index {
		return p.Tags[p.Index]
	}
	return nil
}

func (p *SWFTreeProcessor) Process(actions ActionList) (tag swftag.Tag, newActions ActionList) {
	tag = p.Current()
	if tag == nil {
		return nil, nil
	}

	switch node := tag.(type) {
	case *swftag.DefineMorphShape:
		if p.Loops > 0 {
			break
		}
		//TODO: morphs
	case *swftag.DefineMorphShape2:
	//TODO
	case *swftag.ShowFrame:
	case *swftag.End:
	case *swftag.DoAction:
		for _, action := range node.Actions {
			switch action.ActionCode {
			case subtypes.ActionStop:
				actions = append(actions, &StopAction{})
			case subtypes.ActionPlay:
				actions = append(actions, &PlayAction{})
				//TODO ActionGotoFrame
			case subtypes.ActionNextFrame:
				actions = append(actions, &NextFrameAction{})
				//TODO ActionPreviousFrame

			}
		}
		//TODO DoInitAction

	}

	return tag, actions
}

func (p *SWFTreeProcessor) NextFrame() *ViewFrame {
	var actions ActionList
	if !p.Playing {
		return p.LastFrame
	}

	var node swftag.Tag
	for {
		node, actions = p.Process(actions)
		if node == nil {
			break
		}
		p.Next()

		if _, ok := node.(*swftag.ShowFrame); ok {
			break
		} else if _, ok := node.(*swftag.End); ok && p.Frame == 0 {
			break
		}
	}

	if node == nil { //Loop again
		p.Loops++
		p.Frame = 0
		p.Index = 0
		p.Layout = NewViewLayout(p.Layout.GetObjectId(), nil, nil)
		if p.LastFrame != nil {
			return p.NextFrame()
		}
		return nil
	}

	p.Frame++

	frame := p.Layout.NextFrame(actions)

	//TODO: this might need to be elsewhere?
	for _, action := range actions {
		if _, ok := action.(*StopAction); ok {
			p.Playing = false
		} else if _, ok = action.(*PlayAction); ok {
			p.Playing = true
		} else if _, ok = action.(*NextFrameAction); ok {
			return p.NextFrame()
		}
	}

	return frame
}
