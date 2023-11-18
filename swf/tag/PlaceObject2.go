package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type PlaceObject2 struct {
	_    struct{} `swfFlags:"root,align"`
	Flag struct {
		HasClipActions    bool
		HasClipDepth      bool
		HasName           bool
		HasRatio          bool
		HasColorTransform bool
		HasMatrix         bool
		HasCharacter      bool
		Move              bool
	}
	Depth          uint16
	CharacterId    uint16                        `swfCondition:"Flag.HasCharacter"`
	Matrix         types.Matrix                  `swfCondition:"Flag.HasMatrix"`
	ColorTransform types.ColorTransformWithAlpha `swfCondition:"Flag.HasColorTransform"`
	Ratio          uint16                        `swfCondition:"Flag.HasRatio"`
	Name           string                        `swfCondition:"Flag.HasName"`
	ClipDepth      uint16                        `swfCondition:"Flag.HasClipDepth"`
	ClipActions    subtypes.CLIPACTIONS          `swfCondition:"Flag.HasClipActions"`
}

func (t *PlaceObject2) Code() Code {
	return RecordPlaceObject2
}
