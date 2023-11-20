package types

import "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"

type StyleRecord interface {
}

type LineStyleRecord struct {
	Width types.Twip
	Color Color
}

type FillStyleRecord struct {
	// Fill can be a Color or Gradient
	Fill   any
	Border *LineStyleRecord
}
