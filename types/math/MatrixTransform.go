package math

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"gonum.org/v1/gonum/mat"
	"math"
)

type MatrixTransform struct {
	matrix *mat.Dense
}

var DefaultScale = Vector2[float64]{
	X: 1,
	Y: 1,
}

var DefaultRotateSkew = Vector2[float64]{
	X: 0,
	Y: 0,
}

var DefaultTranslation = Vector2[types.Twip]{
	X: 0,
	Y: 0,
}

func NewMatrixTransform(scale, rotateSkew Vector2[float64], translation Vector2[types.Twip]) MatrixTransform {
	return MatrixTransform{
		//TODO: check order
		matrix: mat.NewDense(3, 3, []float64{
			/* a */ /* c */
			scale.X, rotateSkew.Y, 0,
			/* b */ /* d */
			rotateSkew.X, scale.Y, 0,
			translation.X.Float64(), translation.Y.Float64(), 1,
		}),
	}
}

func ScaleTransform(scale Vector2[float64]) MatrixTransform {
	return NewMatrixTransform(scale, DefaultRotateSkew, DefaultTranslation)
}

func RotateTransform(angle float64) MatrixTransform {
	//TODO: check sin sign location
	sin, cos := math.Sincos(angle)
	return NewMatrixTransform(NewVector2(cos, cos), NewVector2(-sin, sin), DefaultTranslation)
}

func TranslateTransform(translate Vector2[types.Twip]) MatrixTransform {
	return NewMatrixTransform(DefaultScale, DefaultRotateSkew, translate)
}

func IdentityTransform() MatrixTransform {
	return NewMatrixTransform(DefaultScale, DefaultRotateSkew, DefaultTranslation)
}

func SkewXTransform(angle float64) MatrixTransform {
	return NewMatrixTransform(DefaultScale, NewVector2(math.Tan(angle), 0), DefaultTranslation)
}

func SkewYTransform(angle float64) MatrixTransform {
	return NewMatrixTransform(DefaultScale, NewVector2(0, math.Tan(angle)), DefaultTranslation)
}

func (m MatrixTransform) Multiply(o MatrixTransform) MatrixTransform {
	var r mat.Dense
	r.Mul(m.matrix, o.matrix)
	return MatrixTransform{
		matrix: &r,
	}
}

var identityTransform = IdentityTransform()

func (m MatrixTransform) IsIdentity() bool {
	return m.EqualsExact(identityTransform)
}

func (m MatrixTransform) GetA() float64 {
	return m.matrix.At(0, 0)
}

func (m MatrixTransform) GetB() float64 {
	return m.matrix.At(1, 0)
}

func (m MatrixTransform) GetC() float64 {
	return m.matrix.At(0, 1)
}

func (m MatrixTransform) GetD() float64 {
	return m.matrix.At(1, 1)
}

func (m MatrixTransform) GetTX() float64 {
	return m.matrix.At(0, 2)
}

func (m MatrixTransform) GetTY() float64 {
	return m.matrix.At(1, 2)
}

func (m MatrixTransform) GetMatrix() mat.Matrix {
	return m.matrix
}

func (m MatrixTransform) GetTranslation() Vector2[float64] {
	return m.ApplyToVector(NewVector2[float64](0, 0), true)
}

func (m MatrixTransform) ApplyToVector(v Vector2[float64], applyTranslation bool) Vector2[float64] {
	var r mat.VecDense
	if applyTranslation {
		//TODO: check order
		r.MulVec(m.matrix, mat.NewVecDense(3, []float64{v.X, v.Y, 1}))
	} else {
		//TODO: check order
		r.MulVec(m.matrix.Slice(0, 0, 1, 1), mat.NewVecDense(2, []float64{v.X, v.Y}))
	}
	return NewVector2[float64](r.At(0, 0), r.At(1, 0))
}

func (m MatrixTransform) EqualsExact(o MatrixTransform) bool {
	return mat.Equal(m.matrix, o.matrix)
}

const TransformCompareEpsilon = 1e-12

func (m MatrixTransform) Equals(o MatrixTransform, epsilon float64) bool {
	return mat.EqualApprox(m.matrix, o.matrix, epsilon)
}
func (m MatrixTransform) EqualsWithoutTranslation(o MatrixTransform, epsilon float64) bool {
	return mat.EqualApprox(m.matrix.Slice(0, 0, 1, 1), o.matrix.Slice(0, 0, 1, 1), epsilon)
}

func MatrixTransformFromSWF(m types.MATRIX) MatrixTransform {
	return NewMatrixTransform(
		NewVector2[float64](m.ScaleX.Float64(), m.ScaleY.Float64()),
		NewVector2[float64](m.RotateSkew0.Float64(), m.RotateSkew1.Float64()),
		NewVector2[types.Twip](m.TranslateX, m.TranslateY),
	)
}