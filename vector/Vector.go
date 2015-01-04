package vector

import "math"
import "../combinatorics"

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

//An iterator for calculating the determinate.
type detIterator struct {
  det float64
  v [][]float64
}

func (d *detIterator) Iterate(index []int, sig int) {
  var elem float64 = float64(sig)
  for i := 0; i < len(index); i ++ {
    elem *= d.v[i][index[i]]
  }
  d.det += elem
}

//The determinate of a bunch of vectors. 
//Assumes an n * n matrix has been given.
func Det(v [][]float64) float64 {

  I := &detIterator{0, v}

  combinatorics.NestedForPermutation(I, len(v))

  return I.det
}

//An iterator for calculating the generalized cross product.
type crossIterator struct {
  cross []float64
  v [][]float64
}

func (d *crossIterator) Iterate(index []int, sig int) {
  dim := len(index)
  var elem float64 = float64(sig)
  for i := 0; i < dim - 1; i ++ {
    elem *= d.v[i][index[i]]
  }
  d.cross[index[dim - 1]] += elem
}

//The cross product of vectors. 
//Assumes an (n - 1) * n matrix has been given. 
func Cross(v [][]float64) []float64 {
  I := &crossIterator{make([]float64, len(v) + 1), v}

  combinatorics.NestedForPermutation(I, len(v) + 1)

  return I.cross
}
