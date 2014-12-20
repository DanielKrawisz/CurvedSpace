package surface

import "testing"
import "math"
import "sort"
//import "fmt"
import "../test"

//The strategy for testing here is to create inefficient functions to 
//perform the tensor algebra that can be used to check all the 
//complicated formulas in main file. 

//TODO test intersections. 

//TODO Need to adjust error allowance to take propagation of error into account
//or use a more accurate numerical grad tester with 128-bit floats. 
//The error parameter. 
var err_poly float64 = .001
var h_d float64 = .000001

//Inefficient functions for contracting tensors. Testing use only! 
func contractSymmetric4Tensor(t [][][][]float64, x []float64 ) [][][]float64 {
  z := make([][][]float64, len(x))
  index := make([]int, 4)
  symind := make([]int, 4)
  conind := make([]int, 3)

  for i := 0; i < len(x); i ++ {
    z[i] = make([][]float64, i + 1)
    for j := 0; j <= i; j ++ {
      z[i][j] = make([]float64, j + 1)
    }
  }

  for index[3] = 0; index[3] < len(x); index[3] ++ {
    for index[2] = 0; index[2] < len(x); index[2] ++ {
      for index[1] = 0; index[1] < len(x); index[1] ++ {
        for index[0] = 0; index[0] < len(x); index[0] ++ {
          for d := 0; d < 4; d ++ {
            symind[d] = index[d]
          }
          for d := 1; d < 4; d++ {
            conind[d - 1] = index[d]
          }
          if sort.IntsAreSorted(conind) {
            sort.Ints(symind)

            z[conind[2]][conind[1]][conind[0]] += t[symind[3]][symind[2]][symind[1]][symind[0]] * x[index[0]]
          }
        }
      }
    }
  }

  return z
}

func contractSymmetric3Tensor(t [][][]float64, x []float64 ) [][]float64 {
  z := make([][]float64, len(x))
  index := make([]int, 3)
  symind := make([]int, 3)
  conind := make([]int, 2)

  for i := 0; i < len(x); i ++ {
    z[i] = make([]float64, i + 1)
  }

  for index[2] = 0; index[2] < len(x); index[2] ++ {
    for index[1] = 0; index[1] < len(x); index[1] ++ {
      for index[0] = 0; index[0] < len(x); index[0] ++ {
        for d := 0; d < 3; d ++ {
          symind[d] = index[d]
        }
        for d := 1; d < 3; d++ {
          conind[d - 1] = index[d]
        }
        if sort.IntsAreSorted(conind) {
          sort.Ints(symind)

          z[conind[1]][conind[0]] += t[symind[2]][symind[1]][symind[0]] * x[index[0]]
        }
      }
    }
  }

  return z
}

// t is a symmetric 2-tensor. 
func contractSymmetricTensor(t [][]float64, x []float64 ) []float64 {
  z := make([]float64, len(x))
  index := make([]int, 2)
  symind := make([]int, 2)

  for index[1] = 0; index[1] < len(x); index[1] ++ {
    for index[0] = 0; index[0] < len(x); index[0] ++ {
      for d := 0; d < 2; d ++ {
        symind[d] = index[d]
      }
      sort.Ints(symind)

      z[index[0]] += t[symind[1]][symind[0]] * x[index[1]]
    }
  }

  return z
}

func contractVector(t []float64, x []float64 ) float64 {
  var z float64 = 0.0

  for i := 0; i < len(x); i ++ {
    z += t[i] * x[i]
  }

  return z
}

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
      point := make([]float64, dim)
      point[0] = test.RandFloat(-100, 100)
      point[1] = test.RandFloat(-100, 100)
      point[2] = test.RandFloat(-100, 100)
      point[3] = test.RandFloat(-100, 100)

      expect := contractVector(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("linear surface error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, err_poly)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        //TODO using less error than usual here. 
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], .0001)
      }

      if !grad_match {
        t.Error("linear surface defined by b = ", b, "grad error. Expected ", grad_exp, ", got ", grad)
      }
    }
  }
}

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestLinearIntersection(t *testing.T) {
  dim := 3
  for i := 0; i < 4; i ++ {
    plane := &linearSurface{3, []float64{test.RandFloat(1, 2), test.RandFloat(1, 2), test.RandFloat(1, 2)},
      test.RandFloat(-1, 1)}

    for j := 0; j < 4; j ++ {

      p2 := []float64{test.RandFloat(4, 6), test.RandFloat(4, 6), test.RandFloat(4, 6)}
      p1 := []float64{test.RandFloat(-4, -6), test.RandFloat(-4, -6), test.RandFloat(-4, -6)}
      v := make([]float64, len(p1))
      for i := 0; i < len(p1); i ++ {
        v[i] = p2[i] - p1[i]
      }

      u := plane.Intersection(p1, v)
      u_test := testIntersection(plane, p1, v, 100)

      if len(u) == 0 {
        t.Error("No intersection point found for ", plane.String(), ", p1 = ", p1, "; v = ", v)
        return 
      } 

      if len(u_test) == 0 {
        t.Error("No test intersection point found for ", plane.String(), ", p1 = ", p1, "; v = ", v)
        return
      }

      intersection_point := make([]float64, dim)
      for i := 0; i < dim; i++ {
        intersection_point[i] = p1[i] + u[0] * v[i]
      }

      if !test.CloseEnough(plane.F(intersection_point), 0.0, err_bs) {
        t.Error("intersection error for linear surface {",
          plane.b, ", ", plane.a, "} p1 = ", p1, "; v = ", v, "; u = ", u, "; F = ", plane.F(intersection_point))
      }

      close_enough := false
      for _, uu := range u {
        if test.CloseEnough(u_test[0], uu, err_poly) {
          close_enough = true
          break
        }
      }

      if !close_enough {
        t.Error("test point and intersection point do not agree: u = ", u, "; u_test = ", u_test)
      }
    }
  }
}

func TestNewPlane(t *testing.T) {
  if nil != NewPlaneByPointAndNormal(nil, []float64{1, 0}) {
    t.Error("New linear surface error 1")
  }
  if nil != NewPlaneByPointAndNormal([]float64{1, 0}, nil) {
    t.Error("New linear surface error 2")
  }
  if nil != NewPlaneByPointAndNormal([]float64{1, 0}, []float64{0, 0}) {
    t.Error("New linear surface error 3")
  }

  if nil == NewPlaneByPointAndNormal([]float64{1, 0}, []float64{1, 0}) {
    t.Error("New linear surface error 4")
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

      expect := contractVector(contractSymmetricTensor(c, point), point)
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quadratic surface defined by c = ", c, ", b = 0, a = 0 has error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, h_d)

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

      expect := contractVector(contractSymmetricTensor(c, point), point) + contractVector(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " has error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, err_poly)

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
  if nil != NewQuadraticSurfaceByCenterVectorList(nil, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 1")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 2")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 3")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, nil, [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 4")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{}, nil, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 5")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{[]float64{0, 1}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 6")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{nil, []float64{0, 1}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 7")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{nil}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 8")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 9")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 10")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 11")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1}}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 12")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, nil, 1) {
    t.Error("New quadratic surface error 13")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0, 0, 0}, 1) {
    t.Error("New quadratic surface error 14")
  }
  if nil != NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1, 0}},
      [][]float64{[]float64{1, 0}}, []float64{0}, 1) {
    t.Error("New quadratic surface error 15")
  }

  if nil == NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1,0}, []float64{0, 1}},
      [][]float64{}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 16")
  }
  if nil == NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{[]float64{0, 1}}, []float64{0,0}, 1) {
    t.Error("New quadratic surface error 17")
  }
  if nil == NewQuadraticSurfaceByCenterVectorList([]float64{0,0}, [][]float64{[]float64{1,0}},
      [][]float64{}, []float64{0, 0}, 1) {
    t.Error("New quadratic surface error 18")
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
  test.SetSeed(7820)

  dim := 3
  for i := 0; i < 4; i ++ {
    basis := [][]float64{[]float64{1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25), test.RandFloat(-.25, .25)},
      []float64{test.RandFloat(-.25, .25), test.RandFloat(-.25, .25), 1 + test.RandFloat(-.25, .25)}}
    point := []float64{test.RandFloat(-2, 2), test.RandFloat(-2, 2), test.RandFloat(-2, 2)}
    a := test.RandFloat(8, 80)

    quadratic := NewQuadraticSurfaceByCenterVectorList(point, basis, [][]float64{}, make([]float64, dim), a)

    for j := 0; j < 4; j ++ {

      var p1, p2, v []float64

      n := 0
      for {
        p1 = []float64{test.RandFloat(-10, 10), test.RandFloat(-10, 10), test.RandFloat(-10, 10)}
        //fmt.Println("X trial: ", n, "; ", quadratic, "; point = ", p1, " F = ", quadratic.F(p1)) 
        n++
        if !quadratic.Interior(p1) { break }
        if n > 100 { return }
      }

      n = 0
      for {
        //fmt.Println("Y trial ", n) 
        n++
        p2 = []float64{test.RandFloat(point[0] - 1, point[0] + 1),
          test.RandFloat(point[1] - 1, point[1] + 1), test.RandFloat(point[2] - 1, point[2] + 1)}
        if quadratic.Interior(p2) { break }
      }

      v = make([]float64, len(p1))
      for i := 0; i < len(p1); i ++ {
        v[i] = p2[i] - p1[i]
      }

      u := quadratic.Intersection(p1, v)
      u_test := testIntersection(quadratic, p1, v, 100)

      if len(u) == 0 {
        t.Error("No intersection point found for ", quadratic.String(), ", p1 = ", p1, "; v = ", v)
        return 
      } 

      if len(u_test) == 0 {
        t.Error("No test intersection point found for ", quadratic.String(), ", p1 = ", p1, "; v = ", v)
        return
      }

      intersection_point := make([]float64, dim)

      close_enough_test := false
      close_enough_F := true
      f := make([]float64, len(u))
      for q, uu := range u {
        for i := 0; i < dim; i++ {
          intersection_point[i] = p1[i] + uu * v[i]
        }

        f[q] = quadratic.F(intersection_point)

        if !test.CloseEnough(f[q], 0.0, err_bs) {
          close_enough_F = false
        }

        if test.CloseEnough(u_test[0], uu, err_bs) {
          close_enough_test = true
        }
      }

      if !close_enough_F {
        t.Error("intersection error for ", quadratic.String(), ", p1 = ",
          p1, "; v = ", v, "; u = ", u, "; F = ", f)
      }

      if !close_enough_test {
        t.Error("test point and intersection point do not agree: u = ", u, "; u_test = ", u_test)
      }
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

      expect := contractVector(contractSymmetricTensor(contractSymmetric3Tensor(d, point), point), point) + 
        contractVector(contractSymmetricTensor(c, point), point) + contractVector(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("cubic surface defined by d = ", d, ", c = ", c, ", b = ", b, ", a = ", a, " error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, err_poly)

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

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestCubicIntersection(t *testing.T) {

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

      expect := contractVector(contractSymmetricTensor(contractSymmetric3Tensor(contractSymmetric4Tensor(e,
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

      expect := contractVector(contractSymmetricTensor(contractSymmetric3Tensor(contractSymmetric4Tensor(e,
          point), point), point), point) +
        contractVector(contractSymmetricTensor(contractSymmetric3Tensor(d, point), point), point) + 
        contractVector(contractSymmetricTensor(c, point), point) + contractVector(b, point) + a
      val := surface.F(point)

      if ! test.CloseEnough(val, expect, err_poly) {
        t.Error("quartic surface e = ", e, ", d = ", d, ", c = ", c, ", b = ", b, ", a = ", a," error point ", point, "; expected ", expect, ", got ", val)
      }

      grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, err_poly)

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

//The strategy of this test is to ensure that one point of the segment
//is inside and the other is outside the surface. There should always
//be one intersection. 
func TestQuarticIntersection(t *testing.T) {
  
}

