package surface

import "testing"
import "math"
import "sort"
//import "fmt"
import "../test"

//The strategy for testing here is to create inefficient functions to 
//perform the tensor algebra that can be used to check all the 
//complicated formulas in main file. 

//TODO get grad tester to work.

//TODO test intersections. 

//The error parameter. 
var err_poly float64 = .000001

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

    surface := &LinearCurve{dim, b, a}

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

      /*grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, .000000000001)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("linear surface defined by b = ", b, "grad error. Expected ", grad_exp, ", got ", grad)
      }*/
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

    surface := &QuadraticCurve{dim, c, b, 0}

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

    surface := &QuadraticCurve{dim, c, b, a}

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

      /*grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, .000001)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " grad error. Expected ", grad_exp, ", got ", grad)
      }*/
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

    surface := &CubicCurve{dim, d, c, b, a}

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

      /*grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, .000000001)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " grad error. Expected ", grad_exp, ", got ", grad)
      }*/
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

    surface := &QuarticCurve{dim, e, d, c, b, 0.0}

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

    surface := &QuarticCurve{dim, e, d, c, b, a}

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

      /*grad := surface.Gradient(point)
      grad_exp := testGradient(surface, point, .000000001)

      var grad_match bool = true

      for i := 0; i < dim; i++ {
        grad_match = grad_match && test.CloseEnough(grad[i], grad_exp[i], err_poly)
      }

      if !grad_match {
        t.Error("quadratic surface defined by c = ", c, ", b = ", b, ", a = ", a, " grad error. Expected ", grad_exp, ", got ", grad)
      }*/
    }
  }
}
