package records

import (
	"git.gammaspectra.live/WeebDataHoarder/swf2ass-go/types/math"
	stdmath "math"
)

const BezierRecursionLimit = 32
const BezierCurveCollinearityEpsilon = 1e-30
const BezierCurveAngleToleranceEpsilon = 0.01
const BezierCurveAngleTolerance = 0.4

// CubicRecursiveBezier Converts a Cubic bezier curve into line segments
func CubicRecursiveBezier(points []math.Vector2[float64], cuspLimit, angleTolerance, distanceToleranceSquare float64, v1, v2, v3, v4 math.Vector2[float64], level uint) []math.Vector2[float64] {
	if level > BezierRecursionLimit {
		return points
	}

	// Calculate all the mid-points of the line segments
	//----------------------
	x12 := (v1.X + v2.X) / 2
	y12 := (v1.Y + v2.Y) / 2
	x23 := (v2.X + v3.X) / 2
	y23 := (v2.Y + v3.Y) / 2
	x34 := (v3.X + v4.X) / 2
	y34 := (v3.Y + v4.Y) / 2
	x123 := (x12 + x23) / 2
	y123 := (y12 + y23) / 2
	x234 := (x23 + x34) / 2
	y234 := (y23 + y34) / 2
	x1234 := (x123 + x234) / 2
	y1234 := (y123 + y234) / 2

	// Try to approximate the full cubic curve by a single straight line
	//------------------
	dx := v4.X - v1.X
	dy := v4.Y - v1.Y

	d2 := stdmath.Abs((v2.X-v4.X)*dy - (v2.Y-v4.Y)*dx)
	d3 := stdmath.Abs((v3.X-v4.X)*dy - (v3.Y-v4.Y)*dx)

	var da1, da2, k float64

	var val int
	if d2 > BezierCurveCollinearityEpsilon {
		val = 1 << 1
	}
	if d3 > BezierCurveCollinearityEpsilon {
		val = 1
	}

	switch val {
	case 0:
		// All collinear OR p1==p4
		//----------------------
		k = dx*dx + dy*dy
		if k == 0 {
			d2 = v1.SquaredDistance(v2)
			d3 = v4.SquaredDistance(v3)
		} else {
			k = 1 / k
			da1 = v2.X - v1.X
			da2 = v2.Y - v1.Y
			d2 = k * (da1*dx + da2*dy)
			da1 = v3.X - v1.X
			da2 = v3.Y - v1.Y
			d3 = k * (da1*dx + da2*dy)
			if d2 > 0 && d2 < 1 && d3 > 0 && d3 < 1 {
				// Simple collinear case, 1---2---3---4
				// We can leave just two endpoints
				return points
			}
			if d2 <= 0 {
				d2 = v2.SquaredDistance(v1)
			} else if d2 >= 1 {
				d2 = v2.SquaredDistance(v4)
			} else {
				d2 = v2.SquaredDistance(v1.AddVector(math.NewVector2(d2*dx, d2*dy)))
			}

			if d3 <= 0 {
				d3 = v3.SquaredDistance(v1)
			} else if d3 >= 1 {
				d3 = v3.SquaredDistance(v4)
			} else {
				d3 = v3.SquaredDistance(v1.AddVector(math.NewVector2(d2*dx, d2*dy)))
			}
		}
		if d2 > d3 {
			if d2 < distanceToleranceSquare {
				return append(points, v2)
			}
		} else {
			if d3 < distanceToleranceSquare {
				return append(points, v3)
			}
		}
		break

	case 1:
		// p1,p2,p4 are collinear, p3 is significant
		//----------------------
		if d3*d3 <= distanceToleranceSquare*(dx*dx+dy*dy) {
			if angleTolerance < BezierCurveAngleToleranceEpsilon {
				return append(points, math.NewVector2(x23, y23))
			}

			// Angle Condition
			//----------------------
			da1 = stdmath.Abs(stdmath.Atan2(v4.Y-v3.Y, v4.X-v3.X) - stdmath.Atan2(v3.Y-v2.Y, v3.X-v2.X))
			if da1 >= stdmath.Pi {
				da1 = 2*stdmath.Pi - da1
			}

			if da1 < angleTolerance {
				return append(points, v2, v3)
			}

			if cuspLimit != 0.0 {
				if da1 > cuspLimit {
					return append(points, v3)
				}
			}
		}
		break

	case 2:
		// p1,p3,p4 are collinear, p2 is significant
		//----------------------
		if d2*d2 <= distanceToleranceSquare*(dx*dx+dy*dy) {
			if angleTolerance < BezierCurveAngleToleranceEpsilon {
				return append(points, math.NewVector2(x23, y23))
			}

			// Angle Condition
			//----------------------
			da1 = stdmath.Abs(stdmath.Atan2(v3.Y-v2.Y, v3.X-v2.X) - stdmath.Atan2(v2.Y-v1.Y, v2.X-v1.X))
			if da1 >= stdmath.Pi {
				da1 = 2*stdmath.Pi - da1
			}

			if da1 < angleTolerance {
				return append(points, v2, v3)
			}

			if cuspLimit != 0.0 {
				if da1 > cuspLimit {
					return append(points, v2)
				}
			}
		}
		break

	case 3:
		// Regular case
		//-----------------
		if (d2+d3)*(d2+d3) <= distanceToleranceSquare*(dx*dx+dy*dy) {
			// If the curvature doesn't exceed the distance_tolerance value
			// we tend to finish subdivisions.
			//----------------------
			if angleTolerance < BezierCurveAngleToleranceEpsilon {
				return append(points, math.NewVector2(x23, y23))
			}

			// Angle & Cusp Condition
			//----------------------
			k = stdmath.Atan2(v3.Y-v2.Y, v3.X-v2.X)
			da1 = stdmath.Abs(k - stdmath.Atan2(v2.Y-v1.Y, v2.X-v1.X))
			da2 = stdmath.Abs(stdmath.Atan2(v4.Y-v3.Y, v4.X-v3.X) - k)
			if da1 >= stdmath.Pi {
				da1 = 2*stdmath.Pi - da1
			}
			if da2 >= stdmath.Pi {
				da2 = 2*stdmath.Pi - da2
			}

			if da1+da2 < angleTolerance {
				// Finally we can stop the recursion
				//----------------------
				return append(points, math.NewVector2(x23, y23))
			}

			if cuspLimit != 0.0 {
				if da1 > cuspLimit {
					return append(points, v2)
				}

				if da2 > cuspLimit {
					return append(points, v3)
				}
			}
		}
		break
	}

	// Continue subdivision
	//----------------------
	points = CubicRecursiveBezier(points, cuspLimit, angleTolerance, distanceToleranceSquare, v1, math.NewVector2(x12, y12), math.NewVector2(x123, y123), math.NewVector2(x1234, y1234), level+1)
	return CubicRecursiveBezier(points, cuspLimit, angleTolerance, distanceToleranceSquare, math.NewVector2(x1234, y1234), math.NewVector2(x234, y234), math.NewVector2(x34, y34), v4, level+1)
}

// QuadraticRecursiveBezier Converts a Quadratic bezier curve into line segments
func QuadraticRecursiveBezier(points []math.Vector2[float64], angleTolerance, distanceToleranceSquare float64, v1, v2, v3 math.Vector2[float64], level uint) []math.Vector2[float64] {
	if level > BezierRecursionLimit {
		return points
	}

	// Calculate all the mid-points of the line segments
	//----------------------
	x12 := (v1.X + v2.X) / 2
	y12 := (v1.Y + v2.Y) / 2
	x23 := (v2.X + v3.X) / 2
	y23 := (v2.Y + v3.Y) / 2
	x123 := (x12 + x23) / 2
	y123 := (y12 + y23) / 2

	dx := v3.X - v1.X
	dy := v3.Y - v1.Y
	d := stdmath.Abs((v2.X-v3.X)*dy - (v2.Y-v3.Y)*dx)

	if d > BezierCurveCollinearityEpsilon {
		// Regular case
		//-----------------
		if d*d <= distanceToleranceSquare*(dx*dx+dy*dy) {
			// If the curvature doesn't exceed the distance_tolerance value
			// we tend to finish subdivisions.
			//----------------------
			if angleTolerance < BezierCurveAngleToleranceEpsilon {
				return append(points, math.NewVector2(x123, y123))
			}

			// Angle & Cusp Condition
			//----------------------
			da := stdmath.Abs(stdmath.Atan2(v3.Y-v2.Y, v3.X-v2.X) - stdmath.Atan2(v2.Y-v1.Y, v2.X-v1.X))
			if da >= stdmath.Pi {
				da = 2*stdmath.Pi - da
			}

			if da < angleTolerance {
				// Finally we can stop the recursion
				//----------------------
				return append(points, math.NewVector2(x123, y123))
			}
		}
	} else {

		// Collinear case
		//------------------
		da := dx*dx + dy*dy
		if da == 0 {
			d = v1.SquaredDistance(v2)
		} else {
			d = ((v2.X-v1.X)*dx + (v2.Y-v1.Y)*dy) / da
			if d > 0 && d < 1 {
				// Simple collinear case, 1---2---3
				// We can leave just two endpoints
				return points
			}
			if d <= 0 {
				d = v2.SquaredDistance(v1)
			} else if d >= 1 {
				d = v2.SquaredDistance(v3)
			} else {
				d = v2.SquaredDistance(v1.AddVector(math.NewVector2(d*dx, d*dy)))
			}
		}
		if d < distanceToleranceSquare {
			return append(points, v2)
		}
	}

	// Continue subdivision
	//----------------------
	points = QuadraticRecursiveBezier(points, angleTolerance, distanceToleranceSquare, v1, math.NewVector2(x12, y12), math.NewVector2(x123, y123), level+1)
	return QuadraticRecursiveBezier(points, angleTolerance, distanceToleranceSquare, math.NewVector2(x123, y123), math.NewVector2(x23, y23), v3, level+1)
}
