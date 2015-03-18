package surface

import "testing"
import "math"
import "errors"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/vector"

var err_bs float64 = .0001

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
  expected := [][]float64{[]float64{-3./5., -4./5.}, []float64{0, 0}}

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

