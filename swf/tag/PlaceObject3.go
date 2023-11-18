package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/tag/subtypes"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
)

type PlaceObject3 struct {
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
		Reserved          bool
		OpaqueBackground  bool
		HasVisible        bool
		HasImage          bool
		HasClassName      bool
		HasCacheAsBitmap  bool
		HasBlendMode      bool
		HasFilterList     bool
	}
	Depth             uint16
	ClassName         string                        `swfCondition:"HasClassName()"`
	CharacterId       uint16                        `swfCondition:"Flag.HasCharacter"`
	Matrix            types.Matrix                  `swfCondition:"Flag.HasMatrix"`
	ColorTransform    types.ColorTransformWithAlpha `swfCondition:"Flag.HasColorTransform"`
	Ratio             uint16                        `swfCondition:"Flag.HasRatio"`
	Name              string                        `swfCondition:"Flag.HasName"`
	ClipDepth         uint16                        `swfCondition:"Flag.HasClipDepth"`
	SurfaceFilterList subtypes.FILTERLIST           `swfCondition:"Flag.HasFilterList"`
	BlendMode         uint8                         `swfCondition:"Flag.HasBlendMode"`
	BitmapCache       uint8                         `swfCondition:"Flag.HasCacheAsBitmap"`
	Visible           uint8                         `swfCondition:"Flag.HasVisible"`
	BackgroundColor   types.RGBA                    `swfCondition:"Flag.HasBackgroundColor"`
	ClipActions       subtypes.CLIPACTIONS          `swfCondition:"Flag.HasClipActions"`
}

func (t *PlaceObject3) HasClassName(swfVersion uint8) bool {
	return t.Flag.HasClassName || (t.Flag.HasName && t.Flag.HasImage)
}

func (t *PlaceObject3) Code() Code {
	return RecordPlaceObject3
}
