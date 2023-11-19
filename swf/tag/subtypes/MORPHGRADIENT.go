package subtypes

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type MORPHGRADIENT struct {
	_            struct{}          `swfFlags:"root"`
	NumGradients uint8             `swfBits:",4"`
	Records      []MORPHGRADRECORD `swfCount:"NumGradients"`
}

type MORPHGRADRECORD struct {
	StartRatio uint8
	StartColor types.RGBA
	EndRatio   uint8
	EndColor   types.RGBA
}
