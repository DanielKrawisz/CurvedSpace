package surface

import "./polynomials"

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

type LinearCurve struct {
  dimension int
  b []float64
  a float64
}

func (s *LinearCurve) Dimension() int {
  return s.dimension
}

func (s *LinearCurve) F(x []float64) float64 {
  var f float64 = 0
  for i := 0; i < s.dimension; i++ {
    f += s.b[i]*x[i]
  }
  return f + s.a
}

func (s *LinearCurve) Interior(x []float64) bool {
  return s.F(x) >= 0
}

func (s *LinearCurve) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  for i := 0; i < s.dimension; i++ {
    z[i] = s.b[i]
  }
  return z
}

//Solving for b (x + u v) + a == 0
func (s *LinearCurve) Intersection(x, v []float64) []float64 {
  var f = s.F(x)
  var q float64 = 0
  for i := 0; i < 0; i ++ {
    q += v[i] * s.b[i]
  }
  if q == 0.0 {
    return []float64{}
  } else {
    return []float64{-f / q}
  }
}

type QuadraticCurve struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  //Thus, all indices must be in descending order. 
  c [][]float64
  b []float64
  a float64
}

func (s *QuadraticCurve) Dimension() int {
  return s.dimension
}

func (s *QuadraticCurve) F(x []float64) float64 {
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

func (s *QuadraticCurve) Interior(x []float64) bool {
  return s.F(x) >= 0
}

func (s *QuadraticCurve) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)

  for i := 0; i < s.dimension; i++ { 
    z[i] += s.b[i]

    for j := 0; j <= i; j++ {
      z[i] += 2 * s.c[i][j] * x[j]
    }

    for j := i; j < s.dimension; j ++ {
      z[i] += 2 * s.c[j][i] * x[j]
    }
  }

  return z
}

//Solving for c (x + u v) (x + u v) + b (x + u v) + a == 0
func (s *QuadraticCurve) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*x[i]

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

type CubicCurve struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  d [][][]float64
  c [][]float64
  b []float64
  a float64
}

func (s *CubicCurve) Dimension() int {
  return s.dimension
}

func (s *CubicCurve) F(x []float64) float64 {
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


func (s *CubicCurve) Interior(x []float64) bool {
  return s.F(x) >= 0
}

func (s *CubicCurve) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  var index [3]int
  var ordering [3]int
  var inverse [3]int
  var swap int

  for index[0] = 0; index[0] < s.dimension; index[0] ++ { 
    ordering[0] = 0
    inverse[0] = 0
    z[index[0]] += s.b[index[0]]

    for index[1] = 0; index[1] < s.dimension; index[1] ++ {
      ordering[1] = 1
      inverse[1] = 1

      z[index[0]] += 2 * s.c[index[ordering[0]]][index[ordering[1]]] * x[index[1]]

      if index[ordering[1]] > index[ordering[2]] { 
        swap = ordering[1]
        ordering[1] = ordering[0]
        ordering[0] = swap

        inverse[1] = ordering[0]
        inverse[0] = ordering[1]
      }

      for index[2] = 0; index[2] < s.dimension; index[2] ++ {
        ordering[2] = 2
        inverse[2] = 2 
        z[index[0]] += 3 * s.d[index[ordering[0]]][index[ordering[1]]][index[ordering[2]]] * x[index[1]] * x[index[2]]

        for inverse[2] > 0 && inverse[2] > index[ordering[inverse[2] - 1]] {
          swap = ordering[inverse[2]]
          ordering[inverse[2]] = ordering[inverse[2] - 1]
          ordering[inverse[2] - 1] = swap

          swap = inverse[ordering[inverse[2]]] 
          inverse[ordering[inverse[2]]] = inverse[2]
          inverse[2] = swap
        }
      }
    }
  }

  return z
}

//Solving for d (x + u v)^3 + c (x + u v)^2 + b (x + u v) + a == 0
func (s *CubicCurve) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv, dvvv, dxxx, dvvx, dvxx float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*x[i]

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

type QuarticCurve struct {
  dimension int
  //These arrays represent symmetric tensors and include only the lower parts of the tensor. 
  e [][][][]float64
  d [][][]float64
  c [][]float64
  b []float64
  a float64
}

func (s *QuarticCurve) Dimension() int {
  return s.dimension
}

func (s *QuarticCurve) F(x []float64) float64 {
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

func (s *QuarticCurve) Interior(x []float64) bool {
  return s.F(x) >= 0
}

func (s *QuarticCurve) Gradient(x []float64) []float64 {
  z := make([]float64, s.dimension)
  var index [4]int
  var ordering [4]int
  var inverse [4]int
  var swap int

  for index[0] = 0; index[0] < s.dimension; index[0] ++ { 
    ordering[0] = 0
    inverse[0] = 0
    z[index[0]] += s.b[index[0]]

    for index[1] = 0; index[1] < s.dimension; index[1] ++ {
      ordering[1] = 1
      inverse[1] = 1

      z[index[0]] += 2 * s.c[index[ordering[0]]][index[ordering[1]]] * x[index[2]]

      if index[ordering[1]] > index[ordering[2]] { 
        swap = ordering[1]
        ordering[1] = ordering[0]
        ordering[0] = swap

        inverse[1] = ordering[0]
        inverse[0] = ordering[1]
      }

      for index[2] = 0; index[2] < s.dimension; index[2] ++ {
        ordering[2] = 2
        inverse[2] = 2
        z[index[0]] += 3 * s.d[index[ordering[0]]][index[ordering[1]]][index[ordering[2]]] * x[index[1]] * x[index[2]]

        for inverse[2] > 0 && inverse[2] >  index[ordering[inverse[2] - 1]] {
          swap = ordering[inverse[2]]
          ordering[inverse[2]] = ordering[inverse[2] - 1]
          ordering[inverse[2] - 1] = swap

          swap = inverse[ordering[inverse[2]]] 
          inverse[ordering[inverse[2]]] = inverse[2]
          inverse[2] = swap
        }

        for index[3] = 0; index[3] < s.dimension; index[3] ++ {
          ordering[3] = 3
          inverse[3] = 3
          z[index[0]] += 4 * s.e[index[ordering[0]]][index[ordering[1]]][index[ordering[2]]][index[ordering[3]]] * x[index[1]] * x[index[2]] * x[index[3]]

          for inverse[3] > 0 && inverse[3] > index[ordering[inverse[3] - 1]] {
            swap = ordering[inverse[3]]
            ordering[inverse[3]] = ordering[inverse[3] - 1]
            ordering[inverse[3] - 1] = swap

            swap = inverse[ordering[inverse[3]]] 
            inverse[ordering[inverse[3]]] = inverse[3]
            inverse[3] = swap
          }
        }
      }
    }
  }

  return z
}

//TODO: must simplify for testing purposes. 
//Solving for e (x + u v)^4 + d (x + u v)^3 + c (x + u v)^2 + b (x + u v) + a == 0
func (s *QuarticCurve) Intersection(x, v []float64) []float64 {
  var cxx, cvx, cvv, bx, bv, dvvv, dxxx, dvvx, dvxx, evvvv, evvvx, evvxx, evxxx, exxxx float64

  for i := 0; i < s.dimension; i ++ {
    bx += s.b[i]*x[i]
    bv += s.b[i]*x[i]

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
        //they should be tested algebraically and be given unit tests. 
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
