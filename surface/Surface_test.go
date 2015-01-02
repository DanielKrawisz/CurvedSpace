package surface

import "testing"
import "../test"

//This mock object is for testing the Interior and Normal functions.
type mockTestSurface struct {
  dim int
  f float64
  grad []float64
  intersect []float64
}

func (m *mockTestSurface) Dimension() int {
  return m.dim
}

func (m *mockTestSurface) F(x []float64) float64 {
  return m.f
}

func (m *mockTestSurface) Intersection(x, v []float64) []float64 {
  return m.intersect
}

func (m *mockTestSurface) Gradient(x []float64) []float64 {
  return m.grad
}

func (m *mockTestSurface) String() string {
  return "mock test surface"
}


func TestInterior(t *testing.T) {
  var mock Surface
  x := []float64{0,0}

  mock = &mockTestSurface{2, 1, []float64{1,0}, []float64{}}
  if !SurfaceInterior(mock, x) {
    t.Error("Interior error 0")
  }

  mock = &mockTestSurface{2, -1, []float64{1,0}, []float64{}}
  if SurfaceInterior(mock, x) {
    t.Error("Interior error 1")
  }
}

func TestNormal(t *testing.T) {
  var mock Surface
  var dim int = 2
  x := []float64{0, 0}

  tests := [][]float64{[]float64{3, 4}, []float64{0, 0}}
  expected := [][]float64{[]float64{3./5., 4./5.}, []float64{0, 0}}

  for i, ttt := range tests {
    mock = &mockTestSurface{dim, 1, ttt, []float64{}}
    norm := SurfaceNormal(mock, x)
    var close_enough = true

    for j := 0; j < dim; j ++ {
      if !test.CloseEnough(norm[j], expected[i][j], .0000001) {
        close_enough = false;
      }
    }

    if !close_enough {
      t.Error("Normal error ", i)
    }
  }
}

//Some functions useful for testing surfaces. 

//Given an arbitrary surface, and some points defining the endpoints
//of a line segment, this function tests its intersection function
//against the generic one using Newton's method.
func intersectionTester(s Surface, p1, p2 []float64, t *testing.T) {
  if len(p1) != s.Dimension() || len(p2) != s.Dimension() {
    t.Error("Testing error: intersection lines have incorrect dimension.")
    return
  }

  v := make([]float64, s.Dimension())
  for i := 0; i < s.Dimension(); i ++ {
    v[i] = p2[i] - p1[i]
  }

  u := s.Intersection(p1, v)
  u_test, err := testIntersection(s, p1, v, 100)
  if err != nil {
    //TODO Need to check on how rapidly the intersections are converging.
    //We reach the maximum number of steps a lot and sometimes the
    //estimated intersections aren't accurate enough so the test fails.
    t.Error("Max steps reached in intersection test!")
    //return 
  }

  //The given p1 and p2 should attempt to intersect the surface. 
  if len(u) == 0 {
    t.Error("No intersection point found for ", s.String(), " p1 = ", p1, "; v = ", v)
    return 
  } 

  if len(u_test) == 0 {
    t.Error("No test intersection point found for ", s.String(), " p1 = ", p1, "; v = ", v)
    return
  }

  intersection_point := make([]float64, s.Dimension())

  close_enough_F := true
  f := make([]float64, len(u))
  for q, uu := range u {
    for i := 0; i < s.Dimension(); i++ {
      intersection_point[i] = p1[i] + uu * v[i]
    }

    f[q] = s.F(intersection_point)

    //Every intersection point p that the surface returns should satisfy F(p) == 0
    if !test.CloseEnough(f[q], 0.0, err_bs) {
      close_enough_F = false
    }
  }

  close_enough_test, _ := test.MemberCloseEnough(u_test[0], u, err_bs)

  if !close_enough_test {
    t.Error("test point and intersection point do not agree: u = ", u, "; u_test = ", u_test)
  }

  if !close_enough_F {
    t.Error("intersection error for ", s.String(), ", p1 = ",
      p1, "; v = ", v, "; u = ", u, "; F = ", f)
  }
}
