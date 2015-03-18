package polynomialsurfaces

import "testing"
import "math"
//import "fmt"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/vector"
import "github.com/DanielKrawisz/CurvedSpace/surface"

//The strategy for testing here is to create inefficient functions to 
//perform the tensor algebra that can be used to check all the 
//complicated formulas in main file. 

//TODO finish test intersections. 

//TODO Need to adjust error allowance to take propagation of error into account
//or use a more accurate numerical grad tester with 128-bit floats. 
//The error parameter. 
var err_poly float64 = .0001
var h_d float64 = .000001

//Four dimensions is enough to handle all possibilites
func TestLinear(t *testing.T) {
  dim := 4
  for trial := 0; trial < 3; trial ++ {
    a := test.RandFloat(-10, -1)
    b := make([]float64, dim)

    for i := 0; i < dim; i ++ {
      b[i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-4, 4))  
    }

    surface := &linearSurface{dim, b, a}

    //Test 3 random points. 
    for p := 0; p < 3; p ++ {
      point := test.RandFloatVector(-100, 100, dim)

      expect := vector.Dot(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("linear surface error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := test.GradientTester(surface, point, err_poly)

      if !test.VectorCloseEnough(grad, grad_exp, 00001) {
        t.Error("linear surface defined by b = ", b, "grad error. Expected ", grad_exp, ", got ", grad)
      }
    }
  }
}

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestLinearIntersection(t *testing.T) {
  for i := 0; i < 4; i ++ {
    plane := &linearSurface{3, []float64{test.RandFloat(1, 2), test.RandFloat(1, 2), test.RandFloat(1, 2)},
      test.RandFloat(-1, 1)}

    for j := 0; j < 4; j ++ {

      p2 := []float64{test.RandFloat(4, 6), test.RandFloat(4, 6), test.RandFloat(4, 6)}
      p1 := []float64{test.RandFloat(-4, -6), test.RandFloat(-4, -6), test.RandFloat(-4, -6)}

      test.IntersectionTester(plane, p1, p2, t)
    }
  }
}

func TestNewPlane(t *testing.T) {
  if nil != NewPlaneByPointAndNormal(nil, []float64{1, 0}, true) {
    t.Error("New linear surface error 1")
  }
  if nil != NewPlaneByPointAndNormal([]float64{1, 0}, nil, true) {
    t.Error("New linear surface error 2")
  }
  if nil != NewPlaneByPointAndNormal([]float64{1, 0}, []float64{0, 0}, true) {
    t.Error("New linear surface error 3")
  }

  if nil == NewPlaneByPointAndNormal([]float64{1, 0}, []float64{1, 0}, true) {
    t.Error("New linear surface error 4")
  }

  if nil != NewPlaneByPointsAndSignature(nil, 1) {
    t.Error("New linear surface error 5")
  }
  if nil != NewPlaneByPointsAndSignature([][]float64{nil, []float64{0, 1, 0}, []float64{0, 0, 1}}, 1) {
    t.Error("New linear surface error 6")
  }
  if nil != NewPlaneByPointsAndSignature([][]float64{
      []float64{1, 0, 0}, []float64{0, 1}, []float64{0, 0, 1}}, 1) {
    t.Error("New linear surface error 7")
  }
  if nil != NewPlaneByPointsAndSignature([][]float64{
      []float64{1, 0, 0}, []float64{0, 1, 0, 0}, []float64{0, 0, 1}}, 1) {
    t.Error("New linear surface error 8")
  }
  if nil != NewPlaneByPointsAndSignature([][]float64{
      []float64{1, 0, 0}, []float64{1, 0, 0}, []float64{0, 0, 1}}, 1) {
    t.Error("New linear surface error 9")
  }

  if nil == NewPlaneByPointsAndSignature([][]float64{
      []float64{1, 0, 0}, []float64{0, 1, 0}, []float64{0, 0, 1}}, 1) {
    t.Error("New linear surface error 10")
  }

  for i := 0; i < 4; i ++ {
    dim := 4

    point := test.RandFloatVector(-10, 10, dim)
    normal := test.RandFloatVector(-4, 4, dim)
    right_side_out := test.RandBool()

    plane := NewPlaneByPointAndNormal(point, normal, right_side_out)

    for j := 0; j < dim; j ++ {
      test_point := test.RandFloatVector(-10, 10, dim)

      test_f := plane.F(test_point)

      exp := -vector.Dot(vector.Minus(test_point, point), normal)
      if !right_side_out {exp *= -1}

      if !test.CloseEnough(test_f, exp, err_poly) {
        t.Error("plane error: point ", point[j], "; expected ", exp, " got ", test_f)
      }
    }
  }

  //testing that the plane comes out correctly when given n - 1 points. 
  for i := 0; i < 4; i ++ {
    dim := 3 + i % 2

    point := make([][]float64, dim)
    for j := 0; j < dim; j ++ {
      point[j] = test.RandFloatVector(-10, 10, dim)
    }

    plane := NewPlaneByPointsAndSignature(point, 1)

    for j := 0; j < dim; j ++ {
      test_f := plane.F(point[j])
      if !test.CloseEnough(test_f, 0, err_poly) {
        t.Error("plane error: point ", point[j], " should be on the plane; got ", test_f)
      }
    }
  }

  for i := 0; i < 4; i ++ {
    dim := 3 + i % 2

    point := test.RandFloatVector(-10, 10, dim)
    normal := test.RandFloatVector(-5, 5, dim)

    plane := NewPlaneByPointAndNormal(point, normal, true)

    for j := 0; j < dim; j ++ {
      test_p := test.RandFloatVector(-5, 5, dim)
      test_f := plane.F(test_p)
      exp_p := -vector.Dot(vector.Minus(test_p, point), normal)
      if !test.CloseEnough(test_f, exp_p, err_poly) {
        t.Error("plane error: point ", point[j], " should be on the plane; got ", test_f)
      }
    }
  }
}

func TestQuadratic(t *testing.T) {
  dim := 4

  //Test trials with a and b equal to zero. 
  for trial := 0; trial < 5; trial ++ {
    b := make([]float64, dim)
    c := make([][]float64, dim)
    for i := 0; i < dim; i ++ {

      b[i] = 0

      c[i] = make([]float64, i + 1)

      for j := 0; j <= i; j ++ { 
        if i == j {
          c[i][i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-4, 4))
        } else {
          c[i][j] = test.RandFloat(-10, 10)
        }
      } 
    }

    surface := &quadraticSurface{dim, c, b, 0}

    //Test 5 random points. 
    for p := 0; p < 5; p ++ {
      point := make([]float64, dim)
      for i := 0; i < dim; i ++ { 
        point[i] = test.RandFloat(-100, 100)
      }

      expect := vector.Dot(vector.ContractSymmetricTensor(c, point), point)
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quadratic surface defined by c = ", c, ", b = 0, a = 0 has error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := test.GradientTester(surface, point, h_d)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = 0, a = 0; grad error at point ", point, ". Expected ", grad_exp, ", got ", grad)
      }
    }
  }

  //Test arbitrary quadratic surface. 
  for trial := 0; trial < 5; trial ++ {
    a := test.RandFloat(-25, -1)
    b := make([]float64, dim)
    c := make([][]float64, dim)
    for i := 0; i < dim; i ++ {

      b[i] = test.RandFloat(-5, 5)

      c[i] = make([]float64, i + 1)

      for j := 0; j <= i; j ++ { 
        if i == j {
          c[i][i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-3, 1))
        } else {
          c[i][j] = test.RandFloat(-1, 1)
        }
      } 
    }

    surface := &quadraticSurface{dim, c, b, a}

    //Test 5 random points. 
    for p := 0; p < 5; p ++ {
      point := make([]float64, dim)
      for i := 0; i < dim; i ++ { 
        point[i] = test.RandFloat(-25, 25)
      }

      expect := vector.Dot(
        vector.ContractSymmetricTensor(c, point), point) + vector.Dot(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " has error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := test.GradientTester(surface, point, err_poly)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " grad error at point ", point, ". Expected ", grad_exp, ", got ", grad)
      }
    }
  }
}

func TestNewQuadraticSurface(t *testing.T) {
  if nil != NewQuadraticSurface(nil, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 1")
  }
  if nil != NewQuadraticSurface([]float64{0,0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 2")
  }
  if nil != NewQuadraticSurface([]float64{0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 3")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, nil, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 4")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{}, nil, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 5")
  }
  if nil == NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{[]float64{0, 1}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 6")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{nil, []float64{0, 1}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 7")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{nil}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 8")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 9")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 10")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 11")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 12")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, nil, 1) {
    t.Error("New quadratic surface error 13")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0, 0}, 1) {
    t.Error("New quadratic surface error 14")
  }
  if nil != NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0}, 1) {
    t.Error("New quadratic surface error 15")
  }

  if nil == NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 16")
  }
  if nil == NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{[]float64{0, 1}}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 17")
  }
  if nil == NewQuadraticSurface([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 18")
  }

  //Test whether the surface actually comes out as expected. 
  dim := 4; 
  for i := 0; i < 20; i ++ {
    point := make([]float64, dim)
    basis := make([][]float64, dim)
    y := make([]float64, dim)
    a := test.RandFloat(-5, 5)

    for j := 0; j < dim; j ++ {
      point[j] = test.RandFloat(-5, 5)
      basis[j] = make([]float64, dim)
      y[j] = test.RandFloat(-2, 2)

      for k := 0; k < dim; k ++ {
        basis[j][k] = test.RandFloat(-2, 2)
      }
    }

    ib := test.RandInt(0, dim)

    vp := basis[0:ib]
    vn := basis[ib:dim]

    quadratic := NewQuadraticSurface(point, vp, vn, y, a)

    for j := 0; j < 5; j ++ {
      test_point := make([]float64, dim)
      tp := make([]float64, dim)
      for k := 0; k < dim; k ++ {
        test_point[k] = test.RandFloat(-2, 2)
        tp[k] = test_point[k] - point[k]
      }

      f_test := quadratic.F(test_point)

      f_exp := a

      for k := 0; k < dim; k ++ {
        f_exp += -y[k] * tp[k]

        for l := 0; l < dim; l ++ {
          for m := 0; m < len(vp); m ++ {
            f_exp -= tp[k] * vp[m][k] * vp[m][l] * tp[l]
          }
          for m := 0; m < len(vn); m ++ {
            f_exp += tp[k] * vn[m][k] * vn[m][l] * tp[l]
          }
        }
      }

      if !test.CloseEnough(f_test, f_exp, .000001) {
        t.Error("Quadratic surface error for ", quadratic, "; ib = ", ib, ", p = ",
          point, ", vp = , ", vp, ", vn = ", vn, ", y = ", y, 
          "; x = ", test_point, "; expected ", f_exp, " got ", f_test)
      }
    }
  }
}

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. Thus, the kind of quadratic surface being tested
//is somewhat restricted, but that shouldn't matter because the quadratic
//formula underlying it has been elsewhere tested in all its generality
//and this test only needs to ensure that the intersection problem is
//being translated into a quadratic formula properly. 
func TestQuadraticIntersection(t *testing.T) {

  for i := 0; i < 4; i ++ {
    basis := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    point := []float64{test.RandFloat(-2, 2), test.RandFloat(-2, 2), test.RandFloat(-2, 2)}
    b := []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)}
    a := test.RandFloat(8, 80)

    quadratic := NewQuadraticSurface(point, basis, [][]float64{}, b, a)

    for j := 0; j < 4; j ++ {

      var p1, p2[]float64

      n := 0
      for {
        p1 = []float64{test.RandFloat(-10, 10), test.RandFloat(-10, 10), test.RandFloat(-10, 10)}
        n++
        if !surface.SurfaceInterior(quadratic, p1) { break }
        if n > 100 { return }
      }

      n = 0
      for {
        n++
        p2 = []float64{test.RandFloat(point[0] - 1, point[0] + 1),
          test.RandFloat(point[1] - 1, point[1] + 1), test.RandFloat(point[2] - 1, point[2] + 1)}
        if surface.SurfaceInterior(quadratic, p2) { break }
      }

      test.IntersectionTester(quadratic, p1, p2, t)
    }
  }
}

func TestCubic(t *testing.T) {
  dim := 4
  for trial := 0; trial < 5; trial ++ {
    a := test.RandFloat(-25, -1)
    b := make([]float64, dim)
    c := make([][]float64, dim)
    d := make([][][]float64, dim)
    for i := 0; i < dim; i ++ {

      b[i] = test.RandFloat(-5, 5)

      c[i] = make([]float64, i + 1)
      d[i] = make([][]float64, i + 1)

      for j := 0; j <= i; j ++ { 
        c[i][j] = test.RandFloat(-1, 1)

        d[i][j] = make([]float64, j + 1)

        for k := 0; k <= j; k ++ {
          if i == j && j == k {
            d[i][i][i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-4, 0))
          } else {
            d[i][j][k] = test.RandFloat(-.2, .2)
          }
        }
      } 
    }

    surface := &cubicSurface{dim, d, c, b, a}

    //Test 5 random points. 
    for p := 0; p < 5; p ++ {
      point := make([]float64, dim)
      for i := 0; i < dim; i ++ { 
        point[i] = test.RandFloat(-25, 25)
      }

      expect := vector.Dot(
        vector.ContractSymmetricTensor(vector.ContractSymmetric3Tensor(d, point), point), point) + 
        vector.Dot(
          vector.ContractSymmetricTensor(c, point), point) + vector.Dot(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("cubic surface defined by d = ", d, ", c = ", c, ", b = ", b, ", a = ",
          a, " error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := test.GradientTester(surface, point, err_poly)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ",
          c, ", b = ", b, ", a = ", a, " grad error. Expected ", grad_exp, ", got ", grad)
      }
    }
  }
}

func TestNewCubic(t *testing.T) {
  dv3 := [][]float64{[]float64{1, 0, 0}, []float64{0, 1, 0}}
  dv2 := [][]float64{[]float64{1, 0}, []float64{0, 1}}

  if nil != NewCubicSurface([]float64{0,0}, nil, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 1")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv3, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 1")
  }
  if nil != NewCubicSurface(nil, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 1")
  }
  if nil != NewCubicSurface([]float64{0,0,0}, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 2")
  }
  if nil != NewCubicSurface([]float64{0}, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 3")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, nil, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 4")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{}, nil, []float64{0,0}, 1) {
    t.Error("New cubic surface error 5")
  }
  if nil == NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{[]float64{0, 1}}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 6")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{nil, []float64{0, 1}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 7")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{nil}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 8")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 9")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0, 0}}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 10")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 11")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1}}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 12")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, nil, 1) {
    t.Error("New cubic surface error 13")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0, 0}, 1) {
    t.Error("New cubic surface error 14")
  }
  if nil != NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0}, 1) {
    t.Error("New cubic surface error 15")
  }

  if nil == NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 16")
  }
  if nil == NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{[]float64{0, 1}}, []float64{0,0}, 1) {
    t.Error("New cubic surface error 17")
  }
  if nil == NewCubicSurface([]float64{0,0}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New cubic surface error 18")
  }

  //Test whether the surface actually comes out as expected. 
  dim := 1; 
  var a float64
  for i := 0; i < 20; i ++ { 
    point := make([]float64, dim)
    basis := make([][]float64, dim)
    dbasis := make([][]float64, dim)
    y := make([]float64, dim)
    a = test.RandFloat(-5, 5)

    for j := 0; j < dim; j ++ {
      point[j] = test.RandFloat(-5, 5)
      basis[j] = make([]float64, dim)
      dbasis[j] = make([]float64, dim)
      y[j] = test.RandFloat(-2, 2)

      for k := 0; k < dim; k ++ {
        basis[j][k] = test.RandFloat(-2, 2)
        dbasis[j][k] = test.RandFloat(-2, 2)
      }
    }

    ib := test.RandInt(0, dim)

    vp := basis[0:ib]
    vn := basis[ib:dim]

    cubic := NewCubicSurface(point, dbasis, vp, vn, y, a)

    for j := 0; j < 5; j ++ { 
      test_point := make([]float64, dim)
      tp := make([]float64, dim)
      for k := 0; k < dim; k ++ {
        test_point[k] = test.RandFloat(-2, 2)
        tp[k] = test_point[k] - point[k]
      }

      f_test := cubic.F(test_point)

      f_exp := a

      for k := 0; k < dim; k ++ {
        f_exp -= y[k] * tp[k]

        for l := 0; l < dim; l ++ {
          for m := 0; m < len(vp); m ++ {
            f_exp -= tp[k] * vp[m][k] * vp[m][l] * tp[l]
          }
          for m := 0; m < len(vn); m ++ {
            f_exp += tp[k] * vn[m][k] * vn[m][l] * tp[l]
          }

          for m := 0; m < dim; m ++ {
            for n := 0; n < len(dbasis); n ++ {
              f_exp -= dbasis[n][k] * dbasis[n][l] * dbasis[n][m] * tp[k] * tp[l] * tp[m]
            }
          }
        }
      }

      if !test.CloseEnough(f_test, f_exp, .000001) {
        t.Error("Cubic surface error for ", cubic, "; ib = ", ib, ", p = ",
          point, ", vp = , ", vp, ", vn = ", vn, ", y = ", y, 
          "; x = ", test_point, "; expected ", f_exp, " got ", f_test)
      }
    }
  }
}

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestCubicIntersection(t *testing.T) {
  dim := 3 
  for i := 0; i < 4; i ++ { 
    basisC := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    basisD := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    point := []float64{test.RandFloat(-2, 2), test.RandFloat(-2, 2), test.RandFloat(-2, 2)}
    a := test.RandFloat(8, 80)

    cubic := NewCubicSurface(point, basisD, basisC, [][]float64{}, make([]float64, dim), a)

    for j := 0; j < 4; j ++ { 

      var p1, p2[]float64

      n := 0
      for {
        p1 = []float64{test.RandFloat(-10, 10), test.RandFloat(-10, 10), test.RandFloat(-10, 10)}
        n++
        if !surface.SurfaceInterior(cubic, p1) { break }
        if n > 100 { return }
      }

      n = 0
      for {
        n++
        p2 = []float64{test.RandFloat(point[0] - 1, point[0] + 1),
          test.RandFloat(point[1] - 1, point[1] + 1), test.RandFloat(point[2] - 1, point[2] + 1)}
        if surface.SurfaceInterior(cubic, p2) { break }
      }

      test.IntersectionTester(cubic, p1, p2, t)
    }
  }
}

func TestQuartic(t *testing.T) {
  dim := 4

  //Some trials with all parameters other than e set to zero. 
  for trial := 0; trial < 5; trial ++ {
    b := make([]float64, dim)
    c := make([][]float64, dim)
    d := make([][][]float64, dim)
    e := make([][][][]float64, dim)
    for i := 0; i < dim; i ++ {

      b[i] = 0

      c[i] = make([]float64, i + 1)
      d[i] = make([][]float64, i + 1)
      e[i] = make([][][]float64, i + 1)

      for j := 0; j <= i; j ++ { 
        c[i][j] = 0

        d[i][j] = make([]float64, j + 1)
        e[i][j] = make([][]float64, j + 1)

        for k := 0; k <= j; k ++ {
          d[i][j][k] = 0

          e[i][j][k] = make([]float64, k + 1)

          for l := 0; l <= k; l ++ {

            if i == j && j == k && k == l {
              e[i][i][i][i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-5, -1))
            } else {
              e[i][j][k][l] = test.RandFloat(-.04, .04)
            }
          }
        }
      } 
    }

    surface := &quarticSurface{dim, e, d, c, b, 0.0}

    //Test 5 random points. 
    for p := 0; p < 5; p ++ {
      point := make([]float64, dim)
      for i := 0; i < dim; i ++ { 
        point[i] = test.RandFloat(-25, 25)
      }

      expect := vector.Dot(
                  vector.ContractSymmetricTensor(
                    vector.ContractSymmetric3Tensor(
                      vector.ContractSymmetric4Tensor(e,
                        point), point), point), point) 
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quartic surface e = ", e, ", d = 0, c = 0, b = 0, a = 0 error point ", point, "; expected ", expect, ", got ", val)
      }
    }
  }

  //Trials with random values for all parameters.
  for trial := 0; trial < 5; trial ++ {
    a := test.RandFloat(-25, -1)
    b := make([]float64, dim)
    c := make([][]float64, dim)
    d := make([][][]float64, dim)
    e := make([][][][]float64, dim)
    for i := 0; i < dim; i ++ {

      b[i] = test.RandFloat(-5, 5)

      c[i] = make([]float64, i + 1)
      d[i] = make([][]float64, i + 1)
      e[i] = make([][][]float64, i + 1)

      for j := 0; j <= i; j ++ { 
        c[i][j] = test.RandFloat(-1, 1)

        d[i][j] = make([]float64, j + 1)
        e[i][j] = make([][]float64, j + 1)

        for k := 0; k <= j; k ++ {
          d[i][j][k] = test.RandFloat(-.2, .2)

          e[i][j][k] = make([]float64, k + 1)

          for l := 0; l <= k; l ++ {

            if i == j && j == k && k == l {
              e[i][i][i][i] = float64(test.RandSign()) * math.Exp(test.RandFloat(-5, -1))
            } else {
              e[i][j][k][l] = test.RandFloat(-.04, .04)
            }
          }
        }
      } 
    }

    surface := &quarticSurface{dim, e, d, c, b, a}

    //Test 5 random points. 
    for p := 0; p < 5; p ++ {
      point := make([]float64, dim)
      for i := 0; i < dim; i ++ { 
        point[i] = test.RandFloat(-25, 25)
      }

      expect := vector.Dot(
                  vector.ContractSymmetricTensor(
                    vector.ContractSymmetric3Tensor(
                      vector.ContractSymmetric4Tensor(e,
                        point), point), point), point) +
        vector.Dot(
          vector.ContractSymmetricTensor(
            vector.ContractSymmetric3Tensor(d, point), point), point) + 
        vector.Dot(
          vector.ContractSymmetricTensor(c, point), point) +
        vector.Dot(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quartic surface e = ", e, ", d = ", d, ", c = ", c, ", b = ", b, ", a = ", a," error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := test.GradientTester(surface, point, err_poly)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " grad error. Expected ", grad_exp, ", got ", grad)
      }
    }
  }
}

func TestNewQuartic(t *testing.T) {
  dv3 := [][]float64{[]float64{1, 0, 0}, []float64{0, 1, 0}}
  dv2 := [][]float64{[]float64{1, 0}, []float64{0, 1}}

  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, nil,
      [][]float64{[]float64{1,0}, []float64{0, 1}}, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 1")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv3, [][]float64{}, dv3,
      [][]float64{[]float64{1,0}, []float64{0, 1}}, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 1")
  }
  if nil != NewQuarticSurface(nil, dv2, [][]float64{}, dv2, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 1")
  }
  if nil != NewQuarticSurface([]float64{0,0,0}, dv2, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}}, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 2")
  }
  if nil != NewQuarticSurface([]float64{0}, dv2, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}}, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 3")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, nil,
    [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 4")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2,
    [][]float64{}, nil, []float64{0,0}, 1) {
    t.Error("New quartic surface error 5")
  }
  if nil == NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{[]float64{0, 1}}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 6")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{nil, []float64{0, 1}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 7")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{nil}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 8")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 9")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0, 0}}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 10")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 11")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1}}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 12")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, nil, 1) {
    t.Error("New quartic surface error 13")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0, 0}, 1) {
    t.Error("New quartic surface error 14")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0}, 1) {
    t.Error("New quartic surface error 15")
  }
  if nil != NewQuarticSurface([]float64{0,0}, nil, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 16")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv2, nil, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 17")
  }
  if nil != NewQuarticSurface([]float64{0,0}, dv3, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 18")
  }
  if nil != NewQuarticSurface([]float64{0,0}, [][]float64{}, dv3, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 19")
  }

  if nil == NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2,
      [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 16")
  }
  if nil == NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{[]float64{0, 1}}, []float64{0,0}, 1) {
    t.Error("New quartic surface error 17")
  }
  if nil == NewQuarticSurface([]float64{0,0}, dv2, [][]float64{}, dv2, [][]float64{[]float64{1,0}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quartic surface error 18")
  }

  //Test whether the surface actually comes out as expected. 
  dim := 4; 
  for i := 0; i < 20; i ++ {
    point := make([]float64, dim)
    basis := make([][]float64, dim)
    dbasis := make([][]float64, dim)
    ebasis := make([][]float64, dim)
    y := make([]float64, dim)
    a := test.RandFloat(-5, 5)

    for j := 0; j < dim; j ++ {
      point[j] = test.RandFloat(-5, 5)
      basis[j] = make([]float64, dim)
      dbasis[j] = make([]float64, dim)
      ebasis[j] = make([]float64, dim)
      y[j] = test.RandFloat(-2, 2)

      for k := 0; k < dim; k ++ {
        basis[j][k] = test.RandFloat(-2, 2)
        dbasis[j][k] = test.RandFloat(-2, 2)
        ebasis[j][k] = test.RandFloat(-2, 2)
      }
    }

    ib := test.RandInt(0, dim)
    ie := test.RandInt(0, dim)

    vp := basis[0:ib]
    vn := basis[ib:dim]

    vqp := basis[0:ie]
    vqn := basis[ie:dim]

    quartic := NewQuarticSurface(point, vqp, vqn, dbasis, vp, vn, y, a)

    for j := 0; j < 5; j ++ {
      test_point := make([]float64, dim)
      tp := make([]float64, dim)
      for k := 0; k < dim; k ++ {
        test_point[k] = test.RandFloat(-2, 2)
        tp[k] = test_point[k] - point[k]
      }

      f_test := quartic.F(test_point)

      f_exp := a

      for k := 0; k < dim; k ++ {
        f_exp += -y[k] * tp[k]

        for l := 0; l < dim; l ++ {
          for m := 0; m < len(vp); m ++ {
            f_exp -= tp[k] * vp[m][k] * vp[m][l] * tp[l]
          }
          for m := 0; m < len(vn); m ++ {
            f_exp += tp[k] * vn[m][k] * vn[m][l] * tp[l]
          }

          for m := 0; m < dim; m ++ {
            for n := 0; n < len(dbasis); n ++ {
              f_exp -= dbasis[n][k] * dbasis[n][l] * dbasis[n][m] * tp[k] * tp[l] * tp[m]
            }

            for n := 0; n < dim; n ++ {
              for o := 0; o < len(vqp); o ++ {
                f_exp -= vqp[o][k] * vqp[o][l] * vqp[o][m] * vqp[o][n] * tp[k] * tp[l] * tp[m] * tp[n]
              }
              for o := 0; o < len(vqn); o ++ {
                f_exp += vqn[o][k] * vqn[o][l] * vqn[o][m] * vqn[o][n] * tp[k] * tp[l] * tp[m] * tp[n]
              }
            }
          }
        }
      }

      if !test.CloseEnough(f_test, f_exp, .000001) {
        t.Error("Quartic surface error for ", quartic, "; ib = ", ib, ", p = ",
          point, ", vp = , ", vp, ", vn = ", vn, ", y = ", y, 
          "; x = ", test_point, "; expected ", f_exp, " got ", f_test)
      }
    }
  }
}

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestQuarticIntersection(t *testing.T) {

  dim := 3
  for i := 0; i < 4; i ++ { 
    basisC := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    basisD := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    basisE := [][]float64{
      []float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    point := []float64{test.RandFloat(-2, 2), test.RandFloat(-2, 2), test.RandFloat(-2, 2)}
    a := test.RandFloat(8, 80)

    quartic := NewQuarticSurface(point, basisE, [][]float64{}, basisD, basisC, [][]float64{}, make([]float64, dim), a)

    for j := 0; j < 4; j ++ { 

      var p1, p2[]float64

      n := 0
      for {
        p1 = []float64{test.RandFloat(-10, 10), test.RandFloat(-10, 10), test.RandFloat(-10, 10)}
        n++
        if !surface.SurfaceInterior(quartic, p1) { break }
        if n > 100 { return }
      }

      n = 0
      for {
        n++
        p2 = []float64{test.RandFloat(point[0] - 1, point[0] + 1),
          test.RandFloat(point[1] - 1, point[1] + 1), test.RandFloat(point[2] - 1, point[2] + 1)}
        if surface.SurfaceInterior(quartic, p2) { break }
      }

      test.IntersectionTester(quartic, p1, p2, t)
    }
  }
}

