package surface

import "testing"

//Test construction of booleans as well as the Dimension,
// F, Gradient, and Intersection functions.
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
      if SurfaceInterior(boolean, point) != interior_tests[i][j] {
        t.Error("For boolean ", i, ", point ", point, " is on the wrong side.")
      }

      f := boolean.F(point)
      if (f > 0) != interior_tests[i][j] {
        t.Error("For boolean ", i, ", point ", point, " F disagrees with Interior.")
      }
    }
  }

  //Tests for Gradient
  grad_test_points := [][]float64{[]float64{-2, 0}, []float64{2, 0}}
  var grad_test bool

  //addition grad tests.
  grad := add.Gradient(grad_test_points[0])
  grad_expect := a.Gradient(grad_test_points[0])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for addition, test point 1, ", grad_test_points[0])
  }

  grad = add.Gradient(grad_test_points[1])
  grad_expect = b.Gradient(grad_test_points[1])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for addition, test point 2, ", grad_test_points[1])
  }

  //subtraction grad tests.
  grad = sub.Gradient(grad_test_points[0])
  grad_expect = b.Gradient(grad_test_points[0])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == -grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for subtraction, test point 1, ", grad_test_points[0])
  }

  grad = sub.Gradient(grad_test_points[1])
  grad_expect = b.Gradient(grad_test_points[1])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == -grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for subtraction, test point 2, ", grad_test_points[1])
  }

  //intersection grad tests.
  grad = sec.Gradient(grad_test_points[0])
  grad_expect = b.Gradient(grad_test_points[0])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for intersection and test point 1, ", grad_test_points[0])
  }

  grad = sec.Gradient(grad_test_points[1])
  grad_expect = a.Gradient(grad_test_points[1])
  grad_test = true
  for i := 0; i < len(grad); i ++ {
    grad_test = grad_test && (grad[i] == grad_expect[i])
  }
  if !grad_test {
    t.Error("Grad test error for intersection and test point 2, ", grad_test_points[1])
  }

  //TODO intersection tests
}
