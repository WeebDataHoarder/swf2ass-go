package ass

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types"
	"math"
	"strings"
)

type MatrixTransformTag struct {
	Scale    ScaleTag
	Rotation RotationTag
	Shearing ShearingTag

	Transform types.MatrixTransform
}

func NewMatrixTransformTag(transform types.MatrixTransform, scale types.Vector2[float64], rotationX, rotationY, rotationZ float64, shear types.Vector2[float64]) *MatrixTransformTag {
	return &MatrixTransformTag{
		Scale: ScaleTag{
			Scale: scale,
		},
		Rotation: RotationTag{
			RotationX: rotationX,
			RotationY: rotationY,
			RotationZ: rotationZ,
		},
		Shearing: ShearingTag{
			Shear: shear,
		},
		Transform: transform,
	}
}

func (t *MatrixTransformTag) TransitionMatrixTransform(line *Line, transform types.MatrixTransform) PositioningTag {
	t2 := &MatrixTransformTag{}
	return t2.FromMatrixTransform(transform)
}

func (t *MatrixTransformTag) FromMatrixTransform(transform types.MatrixTransform) PositioningTag {
	*t = *MatrixTransformTagFromTransformStable(transform)
	return t
}

func (t *MatrixTransformTag) Encode(event EventTime) string {
	return strings.Join([]string{
		t.Scale.Encode(event),
		t.Rotation.Encode(event),
		t.Shearing.Encode(event),
	}, "")
}

func (t *MatrixTransformTag) Equals(tag Tag) bool {
	if o, ok := tag.(*MatrixTransformTag); ok {
		return t.Transform.Equals(o.Transform, types.TransformCompareEpsilon) && t.Scale.Equals(&o.Scale) && t.Rotation.Equals(&o.Rotation) && t.Shearing.Equals(&o.Shearing)
	}
	return false
}

func MatrixTransformTagFromTransformStable(transform types.MatrixTransform) *MatrixTransformTag {
	//Numerically stable implementation by MrSmile

	a := transform.GetA()
	b := transform.GetB()
	c := transform.GetC()
	d := transform.GetD()

	ac2 := (a * a) + (c * c)
	bd2 := (b * b) + (d * d)

	det := (a * d) - (b * c)
	dot := (a * b) + (c * d)

	var scaleX, scaleY, frx, fry, frz, fax, fay float64

	if ac2 > bd2 {
		if ac2 > 0 {
			frz = math.Atan2(c, a) * (180 / math.Pi)
			scaleX = math.Sqrt(ac2)
			scaleY = math.Abs(det) / math.Sqrt(ac2)
			fax = dot / ac2

			if det < 0 {
				frz = -frz
				frx = 180
			}
		}
	} else {
		if bd2 > 0 {
			frz = math.Atan2(-b, d) * (180 / math.Pi)
			scaleX = math.Abs(det) / math.Sqrt(bd2)
			scaleY = math.Sqrt(bd2)
			fay = dot / bd2

			if det < 0 {
				frz = -frz
				fry = 180
			}
		}
	}

	//TODO: ???
	frz = -frz
	fscx := math.Abs(scaleX) * 100
	fscy := math.Abs(scaleY) * 100

	return NewMatrixTransformTag(transform, types.NewVector2(fscx, fscy), frx, fry, frz, types.NewVector2(fax, fay))
}
