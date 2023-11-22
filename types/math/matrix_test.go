package math

import (
	"golang.org/x/exp/rand"
	"math"
	"testing"
)

type testVector [4][2]float64

const X = 0
const Y = 1

const (
	In = iota
	Scale
	RotSkew
	Trans
)

func TestMatrixTransform(t *testing.T) {
	testVectors := []testVector{
		{{0, 0}, {1, 1}, {0, 0}, {1000, -1000}},
		{{1, 1}, {1, 1}, {0, 0}, {-1000, 1000}},
		{{1, 1}, {1, 1}, {1, 0}, {0, 0}},
		{{1, 1}, {1, 1}, {0, 1}, {0, 0}},
		{{1, 1}, {1, 1}, {-1, 0.5}, {0, 0}},
		{{1, 1}, {1, 1}, {0.5, -1}, {0, 0}},
		{{1, 1}, {1, 1}, {1, -1}, {0, 0}},
		{{1, 1}, {1, 1}, {-1, 1}, {0, 0}},
		{{1, 1}, {1, 1}, {-1, 1}, {0, 0}},
		{{1, 1}, {1, 1}, {1, -1}, {0, 0}},
		{{1, 1}, {1, 1}, {0.5, -1}, {0, 0}},
		{{1, 1}, {1, 1}, {-1, 0.5}, {0, 0}},
		{{1, 1}, {1, 1}, {-1, 0.5}, {0, 0}},
		{{1, 1}, {1, 1}, {0.5, -1}, {0, 0}},
	}

	evilRand := func(r float64) float64 {
		if rand.Intn(9) == 0 {
			return 0
		}
		return r - rand.Float64()*2*r
	}

	for i := 0; i < 100; i++ {
		testVectors = append(testVectors, testVector{
			{1000 - rand.Float64()*2*1000, 1000 - rand.Float64()*2*1000}, {evilRand(6), evilRand(6)}, {evilRand(math.Pi), evilRand(math.Pi)}, {evilRand(1000), evilRand(1000)},
		})
	}

	for i, s := range testVectors {
		in := NewVector2(s[In][X], s[In][1])
		/*
			a = ScaleX
			b = RotSkewX
			c = RotSkewY
			d = ScaleY
			[a c tx] * [x] = [a*x + c*y + tx]
			[b d ty]   [y]   [b*x + d*y + ty]
			[0 0 1 ]   [1]   [1             ]
		*/
		outputAlt := NewVector2(
			s[Scale][X]*in.X+
				s[RotSkew][Y]*in.Y+
				s[Trans][X],
			s[RotSkew][X]*in.X+
				s[Scale][Y]*in.Y+
				s[Trans][Y],
		)
		m := NewMatrixTransform(NewVector2(s[Scale][X], s[Scale][Y]), NewVector2(s[RotSkew][X], s[RotSkew][Y]), NewVector2(s[Trans][X], s[Trans][Y]))
		output := m.ApplyToVector(in, true)
		t.Logf("[#%d] Input\t%s", i, in.String())
		t.Logf("[#%d] Output\t%s", i, output.String())
		t.Logf("[#%d] Expected\t%s", i, outputAlt.String())
		t.Logf("[#%d] Matrix\n%s", i, m.String())
		if !outputAlt.Equals(output) {
			t.Errorf("[#%d] Failed: output mismatch!\n\n", i)
		}
		if m.GetA() != s[Scale][X] {
			t.Fatal()
		}
		if m.GetB() != s[RotSkew][X] {
			t.Fatal()
		}
		if m.GetC() != s[RotSkew][Y] {
			t.Fatal()
		}
		if m.GetD() != s[Scale][Y] {
			t.Fatal()
		}
		if m.GetTX() != s[Trans][X] {
			t.Fatal()
		}
		if m.GetTY() != s[Trans][Y] {
			t.Fatal()
		}
	}
}
