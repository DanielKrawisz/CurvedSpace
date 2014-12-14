package surface

import "testing"

//Test construction of booleans as well as the Dimension,
// F, Gradient, Interior, and Intersection functions.
func TestBooleans(t *testing.T) {
  p1 := []float64{-2,0}
  p2 := []float64{2,0}
  p3 := []float64{0,0,0}
  var r float64 = 3
  a := NewSphere(p1, r)
  b := NewSphere(p2, r)
  c := NewSphere(p3, r)

  if NewAddition(nil, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewAddition(a, nil) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewSubtraction(nil, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewSubtraction(a, nil) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewIntersection(nil, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewIntersection(a, nil) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewAddition(c, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewAddition(a, c) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewSubtraction(c, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewSubtraction(a, c) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewIntersection(c, b) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }
  if NewIntersection(a, c) != nil {
    t.Error("Non-nil value returned for invalid inputs")
  }

  add := NewAddition(a, b)
  sub := NewSubtraction(a, b)
  sec := NewIntersection(a, b)

  if add == nil {
    t.Error("nil value returned for valid inputs")
  }
  if sub == nil {
    t.Error("nil value returned for valid inputs")
  }
  if sec == nil {
    t.Error("nil value returned for valid inputs")
  }

  if add.Dimension() != 2 {
    t.Error("invalid dimension")
  }
  if sub.Dimension() != 2 {
    t.Error("invalid dimension")
  }
  if sec.Dimension() != 2 {
    t.Error("invalid dimension")
  }

  //Tests for F
  test_points := [][]float64{[]float64{-2, 0}, []float64{0, 0}, []float64{2, 0}}
  interior_tests := [][]bool{
    []bool{true,  true,  true}, 
    []bool{true,  false, false},
    []bool{false, true,  false}}

  for i, boolean := range []Surface{add, sub, sec} {
    for j, point := range test_points {
      if boolean.Interior(point) != interior_tests[i][j] {
        t.Error("For boolean ", i, ", point ", point, " is on the wrong side.")
      }

      f := boolean.F(point)
      if (f > 0) != interior_tests[i][j] {
        t.Error("For boolean ", i, ", point ", point, " F disagrees with Interior.")
      }
    }
  }

  //TODO Tests for Gradient

  //TODO complete tests for intersection functions.
}
