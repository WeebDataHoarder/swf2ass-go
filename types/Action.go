package types

import "git.gammaspectra.live/WeebDataHoarder/swf-go/subtypes"

type ActionList []Action

type Action interface {
	Code() subtypes.ActionCode
}

type StopAction struct {
}

func (a *StopAction) Code() subtypes.ActionCode {
	return subtypes.ActionStop
}

type PlayAction struct {
}

func (a *PlayAction) Code() subtypes.ActionCode {
	return subtypes.ActionPlay
}

type NextFrameAction struct {
}

func (a *NextFrameAction) Code() subtypes.ActionCode {
	return subtypes.ActionNextFrame
}
