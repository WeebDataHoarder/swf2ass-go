package math

import (
	"fmt"
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/swf/types"
	"gonum.org/v1/gonum/mat"
	"math"
)

// MatrixTransform The transformation matrix used by Flash display objects.
// The matrix is a 2x3 affine transformation matrix. A Vector2(x, y) is transformed by the matrix in the following way:
//
//	[a c tx] * [x] = [a*x + c*y + tx]
//	[b d ty]   [y]   [b*x + d*y + ty]
//	[0 0 1 ]   [1]   [1             ]
//
// Objects in Flash can only move in units of types.Twip, or 1/20 pixels.
//
// [SWF19 pp.22-24](https://web.archive.org/web/20220205011833if_/https://www.adobe.com/content/dam/acom/en/devnet/pdf/swf-file-format-spec.pdf#page=22)
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

var DefaultTranslation = Vector2[float64]{
	X: 0,
	Y: 0,
}

func NewMatrixTransform(scale, rotateSkew, translation Vector2[float64]) MatrixTransform {
	return MatrixTransform{
		matrix: mat.NewDense(3, 3, []float64{
			/* a */ /* c */ /* tx */
			scale.X, rotateSkew.Y, translation.X,
			/* b */ /* d */ /* ty */
			rotateSkew.X, scale.Y, translation.Y,
			0, 0, 1,
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

func TranslateTransform[T ~int64 | ~float64](translate Vector2[T]) MatrixTransform {
	return NewMatrixTransform(DefaultScale, DefaultRotateSkew, translate.Float64())
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

// GetA Gets the ScaleX factor
func (m MatrixTransform) GetA() float64 {
	return m.matrix.At(0, 0)
}

// GetB Gets the RotateSkewX factor
func (m MatrixTransform) GetB() float64 {
	return m.matrix.At(1, 0)
}

// GetC Gets the RotateSkewY factor
func (m MatrixTransform) GetC() float64 {
	return m.matrix.At(0, 1)
}

// GetD Gets the ScaleY factor
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

// MinimumStrokeWidth
// Given a matrix, calculates the scale for stroke widths.
// TODO: Verify the actual behavior; I think it's more like the average between scaleX and scaleY.
// Does not yet support vertical/horizontal stroke scaling flags.
func (m MatrixTransform) MinimumStrokeWidth() float64 {
	sx := math.Sqrt(m.GetA()*m.GetA() + m.GetB()*m.GetB())
	sy := math.Sqrt(m.GetC()*m.GetC() + m.GetD()*m.GetD())
	return max(sx, sy)
}

func (m MatrixTransform) GetMatrixWithoutTranslation() mat.Matrix {
	return m.matrix.Slice(0, 2, 0, 2)
}

func (m MatrixTransform) GetTranslation() Vector2[float64] {
	return m.ApplyToVector(NewVector2[float64](0, 0), true)
}

func (m MatrixTransform) Inverse() *MatrixTransform {
	var r mat.Dense
	err := r.Inverse(m.matrix)
	if err != nil {
		return nil
	}
	return &MatrixTransform{
		matrix: &r,
	}
}

func MatrixTransformApplyToVector[T ~int64 | ~float64](m MatrixTransform, v Vector2[T], applyTranslation bool) Vector2[T] {
	return Vector2ToType[float64, T](m.ApplyToVector(v.Float64(), applyTranslation))
}

func (m MatrixTransform) ApplyToVector(v Vector2[float64], applyTranslation bool) Vector2[float64] {
	var r mat.VecDense
	if applyTranslation {
		/*
			[a c tx] * [x] = [a*x + c*y + tx]
			[b d ty]   [y]   [b*x + d*y + ty]
			[0 0 1 ]   [1]   [1             ]
		*/
		r.MulVec(m.matrix, mat.NewVecDense(3, []float64{v.X, v.Y, 1}))
	} else {
		/*
			[a c] * [x] = [a*x + c*y]
			[b d]   [y]   [b*x + d*y]
		*/
		r.MulVec(m.GetMatrixWithoutTranslation(), mat.NewVecDense(2, []float64{v.X, v.Y}))
	}
	return NewVector2[float64](r.AtVec(0), r.AtVec(1))
}

func (m MatrixTransform) EqualsExact(o MatrixTransform) bool {
	return mat.Equal(m.matrix, o.matrix)
}

const TransformCompareEpsilon = 1e-12

func (m MatrixTransform) Equals(o MatrixTransform, epsilon float64) bool {
	return mat.EqualApprox(m.matrix, o.matrix, epsilon)
}

func (m MatrixTransform) EqualsWithoutTranslation(o MatrixTransform, epsilon float64) bool {
	return mat.EqualApprox(m.GetMatrixWithoutTranslation(), o.GetMatrixWithoutTranslation(), epsilon)
}

func (m MatrixTransform) String() string {
	return fmt.Sprintf("%#v", mat.Formatted(m.matrix, mat.FormatPython()))
}

func MatrixTransformFromSWF(m types.MATRIX) MatrixTransform {
	return NewMatrixTransform(
		NewVector2(m.ScaleX.Float64(), m.ScaleY.Float64()),
		NewVector2(m.RotateSkew0.Float64(), m.RotateSkew1.Float64()),
		NewVector2(m.TranslateX.Float64(), m.TranslateY.Float64()),
	)
}
