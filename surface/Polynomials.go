package surface

import "strings"
import "fmt"
import "./polynomials"
import "../vector"

//Polynomial surfaces for degrees 1 to 4.
//Since the solutions to the polynomials are already given, 
//this all ends up being just a bunch of tensor algebra. 

//Each polynomial surface of degree n is given as a set of
//totally symmetric tensors of degrees 0 to n. The symmetric
//tensors are represented as nested arrays whose values are
//only defined for indices given in decreasing order. In
//other words, c[1][2] is not defined and will probably give
//an out of bounds exception, whereas c[2][1] or c[2][2] is
//fine.

type linearSurface struct {
  dimension int
  b []float64
  a float64
}

func (s *linearSurface) Dimension() int {
  return s.dimension
}

func (s *linearSurface) String() string {
  return strings.Join([]string{"linearSurface{", fmt.Sprint(s.b), ", ", fmt.Sprint(s.a), "}"}, "")
}

func (s *linearSurface) F(x []float64) float64 {
  var f float64 = 0
  for i := 0; i < s.dimension; i++ {
    f += s.b[i]*x[i]
  }
  return f + s.a
}

func (s *linearSurface) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  for i := 0; i < s.dimension; i++ {
    z[i] = s.b[i]
  }
  return z
}

//Solving for b (x + u v) + a == 0
func (s *linearSurface) Intersection(x, v []float64) []float64 {
  var f = s.F(x)
  var q float64 = 0
  for i := 0; i < s.dimension; i ++ {
    q += v[i] * s.b[i]
  }
  if q == 0.0 {
    return []float64{}
  } else {
    return []float64{-f / q}
  }
}

type quadraticSurface struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  //Thus, all indices must be in descending order. 
  c [][]float64
  b []float64
  a float64
}

func (s *quadraticSurface) Dimension() int {
  return s.dimension
}

func (s *quadraticSurface) String() string {
  return strings.Join([]string{"quadraticSurface{", fmt.Sprint(s.c), ", ", fmt.Sprint(s.b), ", ", fmt.Sprint(s.a), "}"}, "")
}

func (s *quadraticSurface) F(x []float64) float64 {
  var f float64 

  for i := 0; i < s.dimension; i ++ {
    f += s.b[i]*x[i] 

    for j := 0; j < i; j ++ {
      f += 2 * s.c[i][j]*x[i]*x[j]
    }

    f += s.c[i][i]*x[i]*x[i]
  }

  return f + s.a;
}

func (s *quadraticSurface) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)

  for i := 0; i < s.dimension; i++ { 
    z[i] += s.b[i]

    for j := 0; j <= i; j++ {
      z[i] += 2 * s.c[i][j] * x[j]
    }

    for j := i+1; j < s.dimension; j ++ {
      z[i] += 2 * s.c[j][i] * x[j]
    }
  }

  return z
}

//Solving for c (x + u v) (x + u v) + b (x + u v) + a == 0
func (s *quadraticSurface) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*v[i]

    for j := 0; j < i; j ++ { 
      cxx += 2 * s.c[i][j]*x[i]*x[j]
      cvx += s.c[i][j] * (v[i]*x[j] + v[j]*x[i])
      cvv += 2 * s.c[i][j]*v[i]*v[j]
    } 

    cxx += s.c[i][i] * x[i] * x[i]
    cvx += s.c[i][i] * x[i] * v[i]
    cvv += s.c[i][i] * v[i] * v[i]
  }

  pa := s.a + cxx + bx
  pb := bv + 2 * cvx

  if cvv == 0.0 {
    if pb == 0.0 {
      return []float64{}
    } 

    return []float64{-pa / pb}
  } 

  return polynomials.QuadraticFormula(pa / cvv, pb / cvv)
}

type cubicSurface struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  d [][][]float64
  c [][]float64
  b []float64
  a float64
}

func (s *cubicSurface) Dimension() int {
  return s.dimension
}

func (s *cubicSurface) String() string {
  return strings.Join([]string{"cubicSurface{", fmt.Sprint(s.d), ", ",
    fmt.Sprint(s.c), ", ", fmt.Sprint(s.b), ", ", fmt.Sprint(s.a), "}"}, "")
}

func (s *cubicSurface) F(x []float64) float64 {
  var f float64

  for i := 0; i < s.dimension; i ++ {
    f += s.b[i]*x[i]

    for j := 0; j < i; j ++ {
      f += 2 * s.c[i][j]*x[i]*x[j]

      for k := 0; k < j; k ++ {
        f += 6 * s.d[i][j][k]*x[i]*x[j]*x[k]
      }

      f += 3 * s.d[i][j][j]*x[i]*x[j]*x[j]
      f += 3 * s.d[i][i][j]*x[i]*x[i]*x[j]
    }

    f += s.c[i][i]*x[i]*x[i]
    f += s.d[i][i][i]*x[i]*x[i]*x[i]
  }

  return f + s.a
}

func (s *cubicSurface) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  var index [3]int
  var ordering [][]int = make([][]int, 2)
  var inverse, swap int

  ordering[0] = make([]int, 2)
  ordering[1] = make([]int, 3)

  for index[0] = 0; index[0] < s.dimension; index[0] ++ { 
    ordering[0][0] = 0
    ordering[0][1] = 1

    z[index[0]] += s.b[index[0]]

    for index[1] = 0; index[1] < s.dimension; index[1] ++ {

      if index[ordering[0][1]] > index[ordering[0][0]] { 
        swap = ordering[0][1]
        ordering[0][1] = ordering[0][0]
        ordering[0][0] = swap
      }

      z[index[0]] += 2 * s.c[index[ordering[0][0]]][index[ordering[0][1]]] * x[index[1]]

      ordering[1][0] = ordering[0][0]
      ordering[1][1] = ordering[0][1]
      ordering[1][2] = 2
      inverse = 2

      for index[2] = 0; index[2] < s.dimension; index[2] ++ {

        for inverse > 0 && index[2] > index[ordering[1][inverse - 1]] {

          swap = ordering[1][inverse]
          ordering[1][inverse] = ordering[1][inverse - 1]
          ordering[1][inverse - 1] = swap

          inverse --
        }

        z[index[0]] += 3 * s.d[index[ordering[1][0]]][index[ordering[1][1]]][index[ordering[1][2]]] * x[index[1]] * x[index[2]]
      }
    }
  }

  return z
}

//Solving for d (x + u v)^3 + c (x + u v)^2 + b (x + u v) + a == 0
func (s *cubicSurface) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv, dvvv, dxxx, dvvx, dvxx float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*v[i]

    for j := 0; j < i; j ++ { 
      cxx += 2 * s.c[i][j]*x[i]*x[j]
      cvx += s.c[i][j] * (v[i]*x[j] + v[j]*x[i])
      cvv += 2 * s.c[i][j]*v[i]*v[j]

      for k := 0; k < j; k ++ {
        dvvv += 6 * s.d[i][j][k]*v[i]*v[j]*v[k]
        dxxx += 6 * s.d[i][j][k]*x[i]*x[j]*x[k]
        dvvx += s.d[i][j][k] * (2 * v[i]*v[j]*x[k] + 2 * v[i]*v[k]*x[j] + 2 * v[k]*v[j]*x[i])
        dvxx += s.d[i][j][k] * (2 * v[i]*x[j]*x[k] + 2 * v[j]*x[i]*x[k] + 2 * v[k]*x[j]*x[i])
      }

      dvvv += 3 * s.d[i][j][j]*v[i]*v[j]*v[j]
      dxxx += 3 * s.d[i][j][j]*x[i]*x[j]*x[j]
      dvvx += s.d[i][j][j] * (v[j]*v[j]*x[i] + 2 * v[i]*v[j]*x[j])
      dvxx += s.d[i][j][j] * (v[i]*x[j]*x[j] + 2 * v[j]*x[i]*x[j])

      dvvv += 3 * s.d[i][i][j]*v[i]*v[i]*v[j]
      dxxx += 3 * s.d[i][i][j]*x[i]*x[i]*x[j]
      dvvx += s.d[i][i][j] * (2 * v[j]*v[i]*x[i] + v[i]*v[i]*x[j])
      dvxx += s.d[i][i][j] * (v[i]*x[i]*x[j] + 2 * v[j]*x[i]*x[i])
    } 

    cxx += s.c[i][i] * x[i] * x[i]
    cvx += s.c[i][i] * x[i] * v[i]
    cvv += s.c[i][i] * v[i] * v[i]

    dvvv += s.d[i][i][i] * v[i] * v[i] * v[i]
    dxxx += s.d[i][i][i] * x[i] * x[i] * x[i]
    dvvx += s.d[i][i][i] * v[i] * v[i] * x[i]
    dvxx += s.d[i][i][i] * v[i] * x[i] * x[i]
  }

  pa := s.a + cxx + bx + dxxx
  pb := bv + 2 * cvx + 3 * dvxx
  pc := cvv + 3 * dvvx

  if dvvv == 0.0 {
    if pc == 0.0 {
      if pb == 0.0 {
        return []float64{}
      } 

      return []float64{-pa / pb}
    } 

    return polynomials.QuadraticFormula(pa / pc, pb / pc)
  } 

  return polynomials.CubicFormula(pa / dvvv, pb / dvvv, pc / dvvv)
}

type quarticSurface struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  e [][][][]float64
  d [][][]float64
  c [][]float64
  b []float64
  a float64
}

func (s *quarticSurface) Dimension() int {
  return s.dimension
}

func (s *quarticSurface) String() string {
  return strings.Join([]string{"quarticSurface{", fmt.Sprint(s.e), ", ", fmt.Sprint(s.d), ", ",
    fmt.Sprint(s.c), ", ", fmt.Sprint(s.b), ", ", fmt.Sprint(s.a), "}"}, "")
}

func (s *quarticSurface) F(x []float64) float64 {
  var f float64 

  for i := 0; i < s.dimension; i ++ {
    f += s.b[i]*x[i]

    for j := 0; j < i; j ++ {
      f += 2 * s.c[i][j]*x[i]*x[j]

      //For fewer than three dimensions, these cases do not arise. 
      for k := 0; k < j; k ++ {
        f += 6 * s.d[i][j][k]*x[i]*x[j]*x[k]

        //For fewer than four dimensions, this case does not even arise. 
        for l := 0; l < k; l ++ { 
          f += 24 * s.e[i][j][k][l]*x[i]*x[j]*x[k]*x[l]
        }

        f += 12 * s.e[i][j][k][k]*x[i]*x[j]*x[k]*x[k]
        f += 12 * s.e[i][j][j][k]*x[i]*x[j]*x[j]*x[k]
        f += 12 * s.e[i][i][j][k]*x[i]*x[i]*x[j]*x[k]
      }

      f += 3 * s.d[i][j][j]*x[i]*x[j]*x[j]
      f += 3 * s.d[i][i][j]*x[i]*x[i]*x[j]

      f += 4 * s.e[i][j][j][j]*x[i]*x[j]*x[j]*x[j]
      f += 6 * s.e[i][i][j][j]*x[i]*x[i]*x[j]*x[j]
      f += 4 * s.e[i][i][i][j]*x[i]*x[i]*x[i]*x[j]
    }

    f += s.c[i][i]*x[i]*x[i]
    f += s.d[i][i][i]*x[i]*x[i]*x[i]
    f += s.e[i][i][i][i]*x[i]*x[i]*x[i]*x[i]
  }

  return f + s.a;
}

func (s *quarticSurface) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  var index [4]int
  var ordering [][]int = make([][]int, 3)
  var inverse [2]int
  var swap int

  ordering[0] = make([]int, 2)
  ordering[1] = make([]int, 3)
  ordering[2] = make([]int, 4)

  for index[0] = 0; index[0] < s.dimension; index[0] ++ { 
    ordering[0][0] = 0
    ordering[0][1] = 1

    z[index[0]] += s.b[index[0]]

    for index[1] = 0; index[1] < s.dimension; index[1] ++ {

      if index[ordering[0][1]] > index[ordering[0][0]] { 
        swap = ordering[0][1]
        ordering[0][1] = ordering[0][0]
        ordering[0][0] = swap
      }

      z[index[0]] += 2 * s.c[index[ordering[0][0]]][index[ordering[0][1]]] * x[index[1]]

      ordering[1][0] = ordering[0][0]
      ordering[1][1] = ordering[0][1]
      ordering[1][2] = 2
      inverse[0] = 2

      for index[2] = 0; index[2] < s.dimension; index[2] ++ {

        for inverse[0] > 0 && index[2] > index[ordering[1][inverse[0] - 1]] {

          swap = ordering[1][inverse[0]]
          ordering[1][inverse[0]] = ordering[1][inverse[0] - 1]
          ordering[1][inverse[0] - 1] = swap

          inverse[0] --
        }

        z[index[0]] += 3 * s.d[index[ordering[1][0]]][index[ordering[1][1]]][index[ordering[1][2]]] * x[index[1]] * x[index[2]]

        ordering[2][0] = ordering[1][0]
        ordering[2][1] = ordering[1][1]
        ordering[2][2] = ordering[1][2]
        ordering[2][3] = 3
        inverse[1] = 3

        for index[3] = 0; index[3] < s.dimension; index[3] ++ {

          for inverse[1] > 0 && index[3] > index[ordering[2][inverse[1] - 1]] {
            swap = ordering[2][inverse[1]]
            ordering[2][inverse[1]] = ordering[2][inverse[1] - 1]
            ordering[2][inverse[1] - 1] = swap

            inverse[1] --
          }

          z[index[0]] += s.e[index[ordering[2][0]]][index[ordering[2][1]]][index[ordering[2][2]]][index[ordering[2][3]]] *
            4 * x[index[1]] * x[index[2]] * x[index[3]]
        }
      }
    }
  }

  return z
}

//TODO: must simplify so as to be easier to test. 
//Solving for e (x + u v)^4 + d (x + u v)^3 + c (x + u v)^2 + b (x + u v) + a == 0
func (s *quarticSurface) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv, dvvv, dxxx, dvvx, dvxx, evvvv, evvvx, evvxx, evxxx, exxxx float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*v[i]

    for j := 0; j < i; j ++ { 
      cxx += 2 * s.c[i][j]*x[i]*x[j]
      cvx += s.c[i][j] * (v[i]*x[j] + v[j]*x[i])
      cvv += 2 * s.c[i][j]*v[i]*v[j]

      for k := 0; k < j; k ++ {
        dvvv += 6 * s.d[i][j][k]*v[i]*v[j]*v[k]
        dxxx += 6 * s.d[i][j][k]*x[i]*x[j]*x[k]
        dvvx += s.d[i][j][k] * (2 * v[i]*v[j]*x[k] + 2 * v[i]*v[k]*x[j] + 2 * v[k]*v[j]*x[i])
        dvxx += s.d[i][j][k] * (2 * v[i]*x[j]*x[k] + 2 * v[j]*x[i]*x[k] + 2 * v[k]*x[j]*x[i])

        for l := 0; l < k; l ++ {
          evvvv += 24 * s.e[i][j][k][l] * v[i] * v[j] * v[k] * v[l]
          evvvx += s.e[i][j][k][l] * (6 * v[i] * v[j] * v[k] * x[l] + 6 * v[i] * v[j] * v[l] * x[k] + 
                     6 * v[i] * v[l] * v[k] * x[j] + 6 * v[l] * v[j] * v[k] * x[i])
          evvxx += s.e[i][j][k][l] * (4 * v[i] * v[j] * x[k] * x[l] + 4 * v[i] * v[k] * x[j] * x[l] + 
                     4 * v[i] * v[l] * x[k] * x[j] + 4 * v[j] * v[k] * x[i] * x[l] + 
                     4 * v[j] * v[l] * x[i] * x[k] + 4 * v[k] * v[l] * x[i] * x[j])
          evxxx += s.e[i][j][k][l] * (6 * v[i] * x[j] * x[k] * x[l] + 6 * v[i] * x[j] * x[k] * x[l] + 
                     6 * v[i] * x[j] * x[k] * x[l] + 6 * v[i] * x[j] * x[k] * x[l])
          exxxx += 24 * s.e[i][j][k][l] * x[i] * x[j] * x[k] * x[l]
        }

        //TODO basically there is no way these formulas are all correct as written now. 
        //Testing still not done. 
        evvvv += 12 * s.e[i][j][k][k] * v[i] * v[j] * v[k] * v[k]
        evvvx += s.e[i][j][k][k] * (6 * v[i] * v[j] * v[k] * x[k] + 3 * v[i] * v[k] * v[k] * x[j] + 3 * v[k] * v[j] * v[k] * x[i])
        evvxx += s.e[i][j][k][k] * (2 * v[i] * v[j] * x[k] * x[k] + 2 * v[k] * v[k] * x[i] * x[j] +
                   4 * v[i] * v[k] * x[j] * x[k] + 4 * v[j] * v[k] * x[i] * x[k])
        evxxx += s.e[i][j][k][k] * (3 * v[i] * x[j] * x[k] * x[k] + 3 * v[j] * x[i] * x[k] * x[k] + 6 * v[k] * x[j] * x[i] * x[k])
        exxxx += 12 * s.e[i][j][k][k] * x[i] * x[j] * x[k] * x[k]

        evvvv += 12 * s.e[i][j][j][k] * v[i] * v[j] * v[k] * v[k]
        evvvx += s.e[i][j][j][k] * (6 * v[i] * v[j] * v[k] * x[j] + 3 * v[i] * v[j] * v[j] * x[k] + 3 * v[k] * v[j] * v[j] * x[i])
        evvxx += s.e[i][j][j][k] * (4 * v[i] * v[j] * x[j] * x[k] + 4 * v[j] * v[k] * x[i] * x[j] +
                   2 * v[i] * v[k] * x[j] * x[j] + 2 * v[j] * v[j] * x[i] * x[k])
        evxxx += s.e[i][j][j][k] * (3 * v[i] * x[j] * x[j] * x[k] + 3 * v[j] * x[i] * x[k] * x[k] + 6 * v[j] * x[j] * x[i] * x[k])
        exxxx += 12 * s.e[i][j][j][k] * x[i] * x[j] * x[j] * x[k]

        evvvv += 12 * s.e[i][j][j][k] * v[i] * v[j] * v[k] * v[k]
        evvvx += s.e[i][j][j][k] * (6 * v[i] * v[j] * v[k] * x[i] + 3 * v[i] * v[i] * v[j] * x[k] + 3 * v[i] * v[j] * v[i] * x[k])
        evvxx += s.e[i][j][j][k] * (4 * v[i] * v[j] * x[j] * x[k] + 4 * v[j] * v[k] * x[i] * x[j] +
                   2 * v[i] * v[k] * x[j] * x[j] + 2 * v[j] * v[j] * x[i] * x[k])
        evxxx += s.e[i][j][j][k] * (6 * v[i] * x[i] * x[j] * x[k] + 3 * v[j] * x[i] * x[i] * x[k] + 3 * v[k] * x[j] * x[i] * x[k])
        exxxx += 12 * s.e[i][j][j][k] * x[i] * x[i] * x[j] * x[k]
      }

      dvvv += 3 * s.d[i][j][j]*v[i]*v[j]*v[j]
      dxxx += 3 * s.d[i][j][j]*x[i]*x[j]*x[j]
      dvvx += s.d[i][j][j] * (v[j]*v[j]*x[i] + 2 * v[i]*v[j]*x[j])
      dvxx += s.d[i][j][j] * (v[i]*x[j]*x[j] + 2 * v[j]*x[i]*x[j])

      dvvv += 3 * s.d[i][i][j]*v[i]*v[i]*v[j]
      dxxx += 3 * s.d[i][i][j]*x[i]*x[i]*x[j]
      dvvx += s.d[i][i][j] * (2 * v[j]*v[i]*x[i] + v[i]*v[i]*x[j])
      dvxx += s.d[i][i][j] * (v[i]*x[i]*x[j] + 2 * v[j]*x[i]*x[i])

      evvvv += 4 * s.e[i][j][j][j] * v[i] * v[j] * v[j] * v[j]
      evvvx += s.e[i][j][j][j] * (3 * v[i] * v[j] * v[j] * x[j] + v[j] * v[j] * v[j] * x[i])
      evvxx += s.e[i][j][j][j] * (2 * v[i] * v[j] * x[j] * x[j] + 2 * v[j] * v[j] * x[i] * x[j])
      evxxx += s.e[i][j][j][j] * (v[i] * x[j] * x[j] * x[j] + 3 * v[j] * x[i] * x[j] * x[j])
      exxxx += 4 * s.e[i][j][j][j] * x[i] * x[j] * x[j] * x[j]

      evvvv += 6 * s.e[i][i][j][j] * v[i] * v[i] * v[i] * v[j]
      evvvx += s.e[i][i][j][j] * (3 * v[j] * v[j] * v[i] * x[i] + 3 * v[i] * v[i] * v[j] * x[j])
      evvxx += s.e[i][i][j][j] * (v[i] * v[i] * x[j] * x[j] + 4 * v[i] * v[j] * x[j] * x[i] + v[j] * v[j] * x[i] * x[i])
      evxxx += s.e[i][i][j][j] * (3 * v[j] * x[j] * x[i] * x[i] + 3 * v[i] * x[i] * x[j] * x[j])
      exxxx += 6 * s.e[i][i][j][j] * x[i] * x[i] * x[i] * x[j]

      evvvv += 4 * s.e[i][i][i][j] * v[i] * v[i] * v[i] * v[j]
      evvvx += s.e[i][i][i][j] * (3 * v[j] * v[i] * v[i] * x[i] + v[i] * v[i] * v[i] * x[j])
      evvxx += s.e[i][i][i][j] * (2 * v[j] * v[i] * x[i] * x[i] + 2 * v[i] * v[i] * x[j] * x[i])
      evxxx += s.e[i][i][i][j] * (v[j] * x[i] * x[i] * x[i] + 3 * v[i] * x[j] * x[i] * x[i])
      exxxx += 4 * s.e[i][i][i][j] * x[i] * x[i] * x[i] * x[j]
    } 

    cxx += s.c[i][i] * x[i] * x[i]
    cvx += s.c[i][i] * x[i] * v[i]
    cvv += s.c[i][i] * v[i] * v[i]

    dvvv += s.d[i][i][i] * v[i] * v[i] * v[i]
    dxxx += s.d[i][i][i] * x[i] * x[i] * x[i]
    dvvx += s.d[i][i][i] * v[i] * v[i] * x[i]
    dvxx += s.d[i][i][i] * v[i] * x[i] * x[i]

    evvvv += s.e[i][i][i][i] * v[i] * v[i] * v[i] * v[i]
    evvvx += s.e[i][i][i][i] * v[i] * v[i] * v[i] * x[i]
    evvxx += s.e[i][i][i][i] * v[i] * v[i] * x[i] * x[i]
    evxxx += s.e[i][i][i][i] * v[i] * x[i] * x[i] * x[i] 
    exxxx += s.e[i][i][i][i] * x[i] * x[i] * x[i] * x[i]
  }

  pa := s.a + bx + cxx + dxxx + exxxx
  pb := bv + 2 * cvx + 3 * dvxx + 4 * evxxx
  pc := cvv + 3 * dvvx + 6 * evvxx
  pd := dvvv + 4 * evvvx

  if evvvv == 0.0 {
    if pd == 0.0 {
      if pc == 0.0 {
        if pb == 0.0 {
          return []float64{}
        }

        return []float64{-pa / pb}
      } 

      return polynomials.QuadraticFormula(pa / pc, pb / pc)
    } 

    return polynomials.CubicFormula(pa / pd, pb / pd, pc / pd)
  }

  return polynomials.QuarticFormula(pa / evvvv, pb / evvvv, pc / evvvv, pd / evvvv)
}

//May return nil
func NewPlaneByPointAndNormal(point, norm []float64) Surface {
  if point == nil || norm == nil {return nil}
  if len(point) != len(norm) {return nil}

  if vector.Dot(norm, norm) == 0 {return nil}
  return &linearSurface{len(norm), vector.Negative(norm), vector.Dot(norm, point)}
}

//Takes n points and returns the n-1 dimensional
//flat surface that has those n points on it. 
//May return nil
func NewPlaneByPointsAndSignature(p [][]float64, sig int) Surface {
  if p == nil { return nil }
  dim := len(p)
  if dim <= 0 { return nil }
  for i := 0; i < dim; i ++ {
    if p[i] == nil { return nil }
    if len(p[i]) != dim { return nil }
  }

  v := make([][]float64, dim - 1)
  for i := 1; i < dim; i ++ {
    v[i - 1] = make([]float64, dim)
    for j := 0; j < dim; j ++ {
      v[i - 1][j] = p[i][j] - p[0][j]
    }
  }

  b := vector.Cross(v)
  return NewPlaneByPointAndNormal(p[0], vector.Times(-float64(sig), b))
}

//Functions that translate a polynomial object along a given vector.
func translateLinear(s *linearSurface, v []float64) Surface {
  p := vector.Negative(v)
  s.a += vector.Dot(s.b, p)
  return s
}

//The algebra is 
//
// c (x + p)^2 + b (x + p) + a = c x^2 + (b + 2 c p) x + (a + b p + c p^2)
//
func translateQuadratic(s *quadraticSurface, v []float64) Surface {
  p := vector.Negative(v)
  bp := vector.ContractSymmetricTensor(s.c, p)
  s.a += vector.Dot(s.b, p) + vector.Dot(bp, p)
  s.b = vector.Plus(s.b, vector.Times(2, bp))
  return s
}

//The algebra is 
//
// d (x + p)^3 + c (x + p)^2 + b (x + p) + a
//   = d x^3 + (c + 3 d p) x^2 + (b + 2 c p + 3 d p^2) x + (a + b p + c p^2 + d p^3)
//
func translateCubic(s *cubicSurface, v []float64) Surface {
  p    := vector.Negative(v)
  dp   := vector.ContractSymmetric3Tensor(s.d, p)
  dpp  := vector.ContractSymmetricTensor(dp, p)
  dppp := vector.Dot(dpp, p)
  cp   := vector.ContractSymmetricTensor(s.c, p)
  cpp  := vector.Dot(cp, p)
  bp   := vector.Dot(s.b, p)

  s.a += bp + cpp + dppp
  s.b = vector.Plus(s.b, vector.Plus(vector.Times(2, cp), vector.Times(3, dpp)))
  vector.AddToSymmetricTensor(s.c, vector.SymmetricTensorTimes(3, dp))
  return s
}

//The algebra is 
//
// e (x + p)^4 + d (x + p)^3 + c (x + p)^2 + b (x + p) + a
//   = e x^4 + (d + 4 e p) x^3 + (c + 3 d p + 6 e p^2) x^2 + (b + 2 c p + 3 d p^2 + 4 e p^3) x 
//     + (a + b p + c p^2 + d p^3 + e p^4)
//
func translateQuartic(s *quarticSurface, v[]float64) Surface {
  p     := vector.Negative(v)
  ep    := vector.ContractSymmetric4Tensor(s.e, p)
  epp   := vector.ContractSymmetric3Tensor(ep, p)
  eppp  := vector.ContractSymmetricTensor(epp, p)
  epppp := vector.Dot(eppp, p)
  dp    := vector.ContractSymmetric3Tensor(s.d, p)
  dpp   := vector.ContractSymmetricTensor(dp, p)
  dppp  := vector.Dot(dpp, p)
  cp    := vector.ContractSymmetricTensor(s.c, p)
  cpp   := vector.Dot(cp, p)
  bp    := vector.Dot(s.b, p)

  s.a += bp + cpp + dppp + epppp
  s.b = vector.Plus(s.b, vector.Plus(vector.Times(2, cp),
    vector.Plus(vector.Times(3, dpp), vector.Times(4, eppp))))
  vector.AddToSymmetricTensor(vector.AddToSymmetricTensor(s.c, vector.SymmetricTensorTimes(3, dp)), 
    vector.SymmetricTensorTimes(6, epp))
  vector.AddToSymmetric3Tensor(s.d, vector.SymmetricTensor3Times(4, ep))
  return s
}

//The next four functions change the coordinates of a polynomial object.
func coordinateShiftLinear(s *linearSurface, A [][]float64) {
  s.b = vector.MatrixMultiply(A, s.b)
}

func coordinateShiftQuadratic(s *quadraticSurface, A [][]float64) {
  s.b = vector.MatrixMultiply(A, s.b)
  s.c = vector.MatrixMultiplySymmetricTensor(A, s.c)
}

//These next two functions are done, but not all the functions they call are done! 
func coordinateShiftCubic(s *cubicSurface, A [][]float64) {
  s.b = vector.MatrixMultiply(A, s.b)
  s.c = vector.MatrixMultiplySymmetricTensor(A, s.c)
  s.d = vector.MatrixMultiplySymmetric3Tensor(A, s.d)
}

func coordinateShiftQuartic(s *quarticSurface, A [][]float64) {
  s.b = vector.MatrixMultiply(A, s.b)
  s.c = vector.MatrixMultiplySymmetricTensor(A, s.c)
  s.d = vector.MatrixMultiplySymmetric3Tensor(A, s.d)
  s.e = vector.MatrixMultiplySymmetric4Tensor(A, s.e)
}

//A general quadratic surface from a central point and a list of vectors
//defining a quadratic form on the coordinates. The vectors do not need
//to satisfy any particular properties, but a set which is all normal
//to one another is the most general possibility. 
//
//   p  - a vector that translates the entire curve. 
//   vp - the vectors which are positive definite in c.
//   vp - the vectors which are negative definite in c.
//   y  - the linear component of the curve before the curve is translated by p.
//   r2 - a linear component that is also independent if the effects of p. 
//
// (using some mixed notation here, but it works.)
//  (x - p) (v_i v_i) (x - p) + y (x - p) + r2 == 0
//  p.v_i v_i.p - 2 x.v_i v_i.p + x.v_i v_i.x + y.x + y.p + r2 == 0
//
//May return nil
func NewQuadraticSurface(p []float64, vp, vn [][]float64, y[] float64, r2 float64) Surface {
  if p == nil || vp == nil || vn == nil || y == nil {
    return nil
  }

  if len(p) != len(y) {
    return nil
  }

  dim := len(p)
  for i := 0; i < len(vp); i++ {
    if vp[i] == nil || len(vp[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vn); i++ {
    if vn[i] == nil || len(vn[i]) != dim {
      return nil
    }
  }

  var c [][]float64 = make([][]float64, dim)

  //Calculate c.
  for i := 0; i < dim; i ++ {
    c[i] = make([]float64, i + 1)
    for j := 0; j <= i; j ++ {
      for k := 0; k < len(vp); k ++ {
        c[i][j] -= vp[k][i] * vp[k][j]
      }
      for k := 0; k < len(vn); k ++ {
        c[i][j] += vn[k][i] * vn[k][j]
      }
    }
  }

  return translateQuadratic(&quadraticSurface{dim, c, vector.Negative(y), r2}, p)
}

//A general cubic surface from a central point and a list of vectors
//defining a quadratic form on the coordinates. The vectors do not need
//to satisfy any particular properties, but a set which is all normal
//to one another is the most general possibility. 
//
//   p  - a vector that translates the entire curve. 
//   vd - vectors used to make d. 
//   vp - the vectors which are positive definite in c.
//   vp - the vectors which are negative definite in c.
//   y  - the linear component of the curve before the curve is translated by p.
//   r2 - a linear component that is also independent if the effects of p. 
//
// (using some mixed notation here, but it works.)
//  (w_i w_i w_i) (x - p) (x - p) (x - p)
//    + (v_i v_i) (x - p) (x - p) + y (x - p) + r3 == 0
//  w_i.x w_i.x w_i.x - 3 w_i.x w_i.x w_i.p + 3 w_i.x w_i.p w_i.p - w_i.p w_i.p w_i.p
//    + p.v_i v_i.p - 2 x.v_i v_i.p + x.v_i v_i.x + y.x + y.p + r3 == 0
//  w_i.x w_i.x w_i.x + (- 3 w_i.x w_i.x w_i.p + x.v_i v_i.x) 
//    + (3 w_i.x w_i.p w_i.p - 2 x.v_i v_i.p + y.x)
//    + (- w_i.p w_i.p w_i.p + p.v_i v_i.p + y.p + r3) == 0
//
//May return nil
func NewCubicSurface(p []float64, vd, vp, vn [][]float64, y[] float64, r3 float64) Surface {
  if p == nil || vp == nil || vn == nil || y == nil || vd == nil{
    return nil
  }

  if len(p) != len(y) {
    return nil
  }

  dim := len(p)
  for i := 0; i < len(vp); i++ {
    if vp[i] == nil || len(vp[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vn); i++ {
    if vn[i] == nil || len(vn[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vd); i++ {
    if vd[i] == nil || len(vd[i]) != dim {
      return nil
    }
  }

  var c [][]float64 = make([][]float64, dim)
  var d [][][]float64 = make([][][]float64, dim)

  //calculate c and d.
  for i := 0; i < dim; i ++ {
    d[i] = make([][]float64, i + 1)
    c[i] = make([]float64, i + 1)
    for j := 0; j <= i; j ++ {
      d[i][j] = make([]float64, j + 1)
      for k := 0; k <= j; k ++ {
        for l := 0; l < len(vd); l ++ {
          d[i][j][k] -= vd[l][i] * vd[l][j] * vd[l][k]
        }
      }
      for k := 0; k < len(vp); k ++ {
        c[i][j] -= vp[k][i] * vp[k][j]
      }
      for k := 0; k < len(vn); k ++ {
        c[i][j] += vn[k][i] * vn[k][j]
      }
    }
  }

  return translateCubic(&cubicSurface{dim, d, c, vector.Negative(y), r3}, p)
}

//A general cubic surface from a central point and a list of vectors
//defining a quadratic form on the coordinates. The vectors do not need
//to satisfy any particular properties, but a set which is all normal
//to one another is the most general possibility. 
//
//   p   - a vector that translates the entire curve. 
//   vqp - the vectors which are positive definite in e.
//   vqp - the vectors which are negative definite in e.
//   vd  - vectors used to make d. 
//   vp  - the vectors which are positive definite in c.
//   vp  - the vectors which are negative definite in c.
//   y   - the linear component of the curve before the curve is translated by p.
//   r2  - a linear component that is also independent if the effects of p. 
//
// (mixed notation again.)
//  (u_i u_i u_i u_i) (x - p) (x - p) (x - p) (x - p) + (w_i w_i w_i) (x - p) (x - p) (x - p)
//    + (v_i v_i) (x - p) (x - p) + y (x - p) + r4 == 0
//  (u_i.x)^4 - 4 (u_i.p) (u_i.x)^3 + 6 (u_i.p)^2 (u_i.x)^2 + 4 (u_i.p)^3 (u_i.x) + (u_i.p)^4 
//    + w_i.x w_i.x w_i.x - 3 w_i.x w_i.x w_i.p + 3 w_i.x w_i.p w_i.p - w_i.p w_i.p w_i.p
//    + p.v_i v_i.p - 2 x.v_i v_i.p + x.v_i v_i.x + y.x + y.p + r4 == 0
//  (u_i.x)^4 + (- 4 (u_i.p) (u_i.x)^3 + w_i.x w_i.x w_i.x)
//    + (6 (u_i.p)^2 (u_i.x)^2 - 3 w_i.x w_i.x w_i.p + x.v_i v_i.x) 
//    + (4 (u_i.p)^3 (u_i.x) + 3 w_i.x w_i.p w_i.p - 2 x.v_i v_i.p + y.x)
//    + ((u_i.p)^4 - w_i.p w_i.p w_i.p + p.v_i v_i.p + y.p + r4) == 0
//
//May return nil
func NewQuarticSurface(p []float64, vqp, vqn, vd, vp, vn [][]float64, y[] float64, r4 float64) Surface {
  if p == nil || vp == nil || vn == nil || y == nil || vd == nil || vqp == nil || vqn == nil {
    return nil
  }

  if len(p) != len(y) {
    return nil
  }

  dim := len(p)
  for i := 0; i < len(vp); i++ {
    if vp[i] == nil || len(vp[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vn); i++ {
    if vn[i] == nil || len(vn[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vd); i++ {
    if vd[i] == nil || len(vd[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vqp); i++ {
    if vqp[i] == nil || len(vqp[i]) != dim {
      return nil
    }
  }
  for i := 0; i < len(vqn); i++ {
    if vqn[i] == nil || len(vqn[i]) != dim {
      return nil
    }
  }

  var c [][]float64 = make([][]float64, dim)
  var d [][][]float64 = make([][][]float64, dim)
  var e [][][][]float64 = make([][][][]float64, dim)

  //calculate c, d, and e.
  for i := 0; i < dim; i ++ {
    c[i] = make([]float64, i + 1)
    d[i] = make([][]float64, i + 1)
    e[i] = make([][][]float64, i + 1)

    for j := 0; j <= i; j ++ {
      d[i][j] = make([]float64, j + 1)
      e[i][j] = make([][]float64, j + 1)

      for k := 0; k <= j; k ++ {
        e[i][j][k] = make([]float64, k + 1)

        for l := 0; l < len(vd); l ++ {
          d[i][j][k] -= vd[l][i] * vd[l][j] * vd[l][k]
        }

        for l := 0; l <= k; l ++ {
          for m := 0; m < len(vqp); m ++ {
            e[i][j][k][l] -= vqp[m][i] * vqp[m][j] * vqp[m][k] * vqp[m][l]
          }
          for m := 0; m < len(vqn); m ++ {
            e[i][j][k][l] += vqn[m][i] * vqn[m][j] * vqn[m][k] * vqn[m][l]
          }
        }
      }

      for k := 0; k < len(vp); k ++ {
        c[i][j] -= vp[k][i] * vp[k][j]
      }
      for k := 0; k < len(vn); k ++ {
        c[i][j] += vn[k][i] * vn[k][j]
      }
    }
  }

  return translateQuartic(&quarticSurface{dim, e, d, c, vector.Negative(y), r4}, p)
}
