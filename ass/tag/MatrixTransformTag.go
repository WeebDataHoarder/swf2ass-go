package tag

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/ass/time"
	math2 "git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	"math"
	"strings"
)

type MatrixTransformTag struct {
	Scale    ScaleTag
	Rotation RotationTag
	Shearing ShearingTag

	Transform math2.MatrixTransform
}

func NewMatrixTransformTag(transform math2.MatrixTransform, scale math2.Vector2[float64], rotationX, rotationY, rotationZ float64, shear math2.Vector2[float64]) *MatrixTransformTag {
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

func (t *MatrixTransformTag) TransitionMatrixTransform(event Event, transform math2.MatrixTransform) PositioningTag {
	t2 := &MatrixTransformTag{}
	return t2.FromMatrixTransform(transform)
}

func (t *MatrixTransformTag) FromMatrixTransform(transform math2.MatrixTransform) PositioningTag {
	*t = *MatrixTransformTagFromTransformStable(transform)
	return t
}

func (t *MatrixTransformTag) Encode(event time.EventTime) string {
	return strings.Join([]string{
		t.Scale.Encode(event),
		t.Rotation.Encode(event),
		t.Shearing.Encode(event),
	}, "")
}

func (t *MatrixTransformTag) Equals(tag Tag) bool {
	if o, ok := tag.(*MatrixTransformTag); ok {
		return t.Transform.Equals(o.Transform, math2.TransformCompareEpsilon) && t.Scale.Equals(&o.Scale) && t.Rotation.Equals(&o.Rotation) && t.Shearing.Equals(&o.Shearing)
	}
	return false
}

// MatrixTransformTagFromTransformUnstable Finds matching ASS \fscx, \fscy, \frx, \fry, \frz, \fax, \fay for a given math.MatrixTransform
// Numerically unstable implementation by Oneric
func MatrixTransformTagFromTransformUnstable(transform math2.MatrixTransform) *MatrixTransformTag {

	a := transform.GetA()
	b := transform.GetB()
	c := transform.GetC()
	d := transform.GetD()

	var scaleX, scaleY, frx, fry, frz, fax, fay float64

	isZero := func(v float64) bool {
		return math.Abs(v) <= math2.TransformCompareEpsilon
	}

	if !((isZero(a) && !isZero(b)) || (isZero(d) && !isZero(c))) {
		//Trivial case
		scaleX = a
		scaleY = d

		if !isZero(a) {
			fax = b / a
		}
		if !isZero(d) {
			fay = c / d
		}

		if scaleX < 0 {
			fry = 180
		}

		if scaleY < 0 {
			frx = 180
		}
	} else if !((isZero(b) && !isZero(a)) || (isZero(c) && !isZero(d))) {
		//Rowswap
		frz = 90
		scaleX = c
		scaleY = -b

		if !isZero(c) {
			fax = d / c
		}
		if !isZero(b) {
			fay = a / b
		}

		if scaleX < 0 {
			fry = 180
		}

		if scaleY < 0 {
			frx = 180
		}
	} else if isZero(a) && isZero(c) && !isZero(b) && !isZero(d) {
		//Zero col left
		scaleX = math.Sqrt(b*b + d*d)
		frz = math.Atan(-b/d) * (180 / math.Pi)
		if a < 0 { // atan always yields positive cos
			frz += 180
		}
	} else if !isZero(a) && !isZero(c) && isZero(b) && isZero(d) {
		//Zero col right
		scaleX = math.Sqrt(a*a + c*c)
		frz = math.Atan(c/a) * (180 / math.Pi)
		if a < 0 { // atan always yields positive cos
			frz += 180
		}
	} else {
		panic("invalid transform state")
	}

	frz = -frz

	fscx := math.Abs(scaleX) * 100
	fscy := math.Abs(scaleY) * 100

	return NewMatrixTransformTag(transform, math2.NewVector2(fscx, fscy), frx, fry, frz, math2.NewVector2(fax, fay))
}

// MatrixTransformTagFromTransformStable Finds matching ASS \fscx, \fscy, \frx, \fry, \frz, \fax, \fay for a given math.MatrixTransform
// Numerically stable implementation by MrSmile
func MatrixTransformTagFromTransformStable(transform math2.MatrixTransform) *MatrixTransformTag {

	a := transform.GetA()
	b := transform.GetB()
	c := transform.GetC()
	d := transform.GetD()

	ab2 := (a * a) + (b * b)
	cd2 := (c * c) + (d * d)

	det := (a * d) - (c * b)
	dot := (a * c) + (b * d)

	var scaleX, scaleY, frx, fry, frz, fax, fay float64

	if ab2 > cd2 {
		if ab2 > 0 {
			frz = math.Atan2(b, a) * (180 / math.Pi)
			scaleX = math.Sqrt(ab2)
			scaleY = math.Abs(det) / math.Sqrt(ab2)
			fax = dot / ab2

			if det < 0 {
				frz = -frz
				frx = 180
			}
		}
	} else {
		if cd2 > 0 {
			frz = math.Atan2(-c, d) * (180 / math.Pi)
			scaleX = math.Abs(det) / math.Sqrt(cd2)
			scaleY = math.Sqrt(cd2)
			fay = dot / cd2

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

	return NewMatrixTransformTag(transform, math2.NewVector2(fscx, fscy), frx, fry, frz, math2.NewVector2(fax, fay))
}
