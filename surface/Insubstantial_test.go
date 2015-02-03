package surface

import "testing"
import "math/rand"
import "github.com/DanielKrawisz/CurvedSpace/test"

var mockInsubstantialValue float64

func mockInsubstantialRand() float64 {
  return mockInsubstantialValue
}

func TestInsubstantial(t *testing.T) {
  insubstantialRand = mockInsubstantialRand

  insub := NewInsubstantialSurface(1, 1.05)

  x, v := []float64{0}, []float64{1}

  for i, val := range [][]float64{{1, 0}, {.5, 1.05}, {.25, 2.1}, {.125, 3.15}} {
    mockInsubstantialValue = val[0]

    got := insub.Intersection(x, v)

    if ! test.CloseEnough(got[0], val[1], .000001) {
      t.Error("insubstantial test case ", i, " expected ", val[1], " got ", got)
    }
  }

  insubstantialRand = rand.Float64
}
