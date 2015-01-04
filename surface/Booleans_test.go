package surface

import "testing"
import "sort"
import "../test"

var b_err float64 = .00001

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
}

func TestBooleanRayIntersections(t *testing.T) {
  p1 := []float64{-2,0}
  p2 := []float64{2,0}
  var r float64 = 4
  s1 := NewSphere(p1, r)
  s2 := NewSphere(p2, r)

  rays := [][][]float64{
    [][]float64{[]float64{-3, -5}, []float64{0, 1}},
    [][]float64{[]float64{-1, -5}, []float64{0, 1}}, 
    [][]float64{[]float64{1, -5}, []float64{0, 1}},
    [][]float64{[]float64{3, -5}, []float64{0, 1}}}

  int1 := make([][]float64, 4)
  int2 := make([][]float64, 4)
  intAdd := make([][]float64, 4)
  intSub := make([][]float64, 4)
  intSec := make([][]float64, 4)

  for i := 0; i < 4; i ++ {
    int1[i] = s1.Intersection(rays[i][0], rays[i][1])
    int2[i] = s2.Intersection(rays[i][0], rays[i][1])
    sort.Float64s(int1[i])
    sort.Float64s(int2[i])
  }

  add := NewAddition(s1, s2)
  sub := NewSubtraction(s1, s2)
  sec := NewIntersection(s1, s2)

  for i := 0; i < 4; i ++ {
    intAdd[i] = add.Intersection(rays[i][0], rays[i][1])
    intSub[i] = sub.Intersection(rays[i][0], rays[i][1])
    intSec[i] = sec.Intersection(rays[i][0], rays[i][1])
    sort.Float64s(intAdd[i])
    sort.Float64s(intSub[i])
    sort.Float64s(intSec[i])
  }

  expAdd := [][]float64{int1[0], int1[1], int2[2], int2[3]}
  expSub := [][]float64{int1[0],
    []float64{int1[1][0], int1[1][1], int2[1][0], int2[1][1]},
    []float64{}, []float64{}}
  expSec := [][]float64{[]float64{}, int2[1], int1[2], []float64{}}
  sort.Float64s(expSub[1])

  for i := 0; i < 4; i ++ {
    if !test.VectorCloseEnough(intAdd[i], expAdd[i], b_err) {
      t.Error("boolean add intersection error; case ", i, "; expected = ", expAdd[i], "; got = ", intAdd[i])
    }
    if !test.VectorCloseEnough(intSub[i], expSub[i], b_err) {
      t.Error("boolean sub intersection error; case ", i, "; expected = ", expSub[i], "; got = ", intSub[i])
    }
    if !test.VectorCloseEnough(intSec[i], expSec[i], b_err) {
      t.Error("boolean sec intersection error; case ", i, "; expected = ", expSec[i], "; got = ", intSec[i])
    }
  }
}
