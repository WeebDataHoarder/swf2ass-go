package shapes

import (
	types2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/records"
)

type Complex interface {
	Draw() []types2.Record
}
