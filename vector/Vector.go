package vector

import "math"

//TODO tests. 
//TODO There are some places around this program that should be using
//these functions, but don't. Check around to fix that. 

func Dot(A, B []float64) (d float64) {
  for i := 0; i < len(A); i ++ {
    d += A[i] * B[i]
  }
  return
}

func Times(a float64, v []float64) []float64 {
  for i := 0; i < len(v); i ++ {
    v[i] *= a
  }
  return v
}

func Plus(A, B []float64) (C []float64) {
  C = make([]float64, len(A))
  for i := 0; i < len(A); i ++ {
    C[i] = A[i] + B[i]
  }
  return
}

func Minus(A, B []float64) (C []float64) {
  for i := 0; i < len(A); i ++ {
    C[i] = A[i] - B[i]
  }
  return
}

//Gram Schmidt process on a bunch of vectors. 
func Orthonormalize(v [][]float64) {
  var d float64
  for i := 0; i < len(v); i ++ {
    for j := 0; j < i; j ++ {
      d = Dot(v[i], v[j])

      for k := 0; k < len(v[i]); k ++ {
        v[i][k] -= d * v[j][k]
      }
    }

    d = Dot(v[i], v[i])

    if d == 0 {
      for j := 0; j < len(v); j ++ {
        v[i][j] = 0
      }
    } else {
      d = math.Sqrt(d)
      for j := 0; j < len(v); j ++ {
        v[i][j] /= d
      }
    }
  }
}
