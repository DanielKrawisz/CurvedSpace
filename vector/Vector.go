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

func Length(V []float64) float64 {
  return math.Sqrt(Dot(V, V))
}

func Normalize(v []float64) []float64 {
  d := Length(v)
  for j := 0; j < len(v); j ++ {
    v[j] /= d
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
  C = make([]float64, len(A))
  for i := 0; i < len(A); i ++ {
    C[i] = A[i] - B[i]
  }
  return
}

func Negative(v []float64) []float64 {
  p := make([]float64, len(v))
  for i := 0; i < len(v); i ++ {
    p[i] = -v[i]
  }
  return p
}

func Times(a float64, v []float64) []float64 {
  for i := 0; i < len(v); i ++ {
    v[i] *= a
  }
  return v
}

func LinearSum(a, b float64, A, B []float64) (C []float64) {
  C = make([]float64, len(A))
  for i := 0; i < len(A); i ++ {
    C[i] = a * A[i] + b * B[i]
  }
  return
}

//Gram Schmidt process on a bunch of vectors. 
//May not return an orthonormal set! If the set of vectors
//is linearly dependent, some of them will end up as zero
//vectors.
func Orthonormalize(v [][]float64) [][]float64 {
  var d float64
  var swap []float64
  D := len(v)

  for i := 0; i < D; i ++ {
    for j := 0; j < i; j ++ {
      d = Dot(v[i], v[j])

      for k := 0; k < len(v[i]); k ++ {
        v[i][k] -= d * v[j][k]
      }
    }

    d = Dot(v[i], v[i])

    //Zero vectors. 
    if d == 0 {
      for j := 0; j < len(v[i]); j ++ {
        v[i][j] = 0
      }
      //move the zero vector to the back.
      swap = v[i]
      D --
      for j := i; j < D; j ++ {
        v[j] = v[j + 1]
      }
      v[D] = swap
      i --
    } else {
      d = math.Sqrt(d)
      for j := 0; j < len(v[i]); j ++ {
        v[i][j] /= d
      }
    }
  }

  return v
}

//An iterator for calculating the determinate.
type detIterator struct {
  det float64
  v [][]float64
}

func (d *detIterator) Iterate(index []uint, sig int) {
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

  combinatorics.NestedForPermutation(I, uint(len(v)))

  return I.det
}

//An iterator for calculating the generalized cross product.
type crossIterator struct {
  cross []float64
  v [][]float64
}

func (d *crossIterator) Iterate(index []uint, sig int) {
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

  combinatorics.NestedForPermutation(I, uint(len(v) + 1))

  return I.cross
}

func Transpose(m [][]float64) [][]float64 {
  var swap float64
  for i := 0; i < len(m); i ++ {
    for j := 0; j < i; j ++ {
      swap = m[i][j]
      m[i][j] = m[j][i]
      m[j][i] = swap
    }
  }

  return m
}

func elementaryRowOperationPlus(m [][]float64, a, b int, d float64) {
  for i := 0; i < len(m); i ++ {
    m[b][i] = m[b][i] + d * m[a][i]
  }
}

func elementaryRowOperationTimes(m [][]float64, a int, d float64) {
  for i := 0; i < len(m); i ++ {
    m[a][i] *= d
  }
}

func elementaryRowOperationSwap(m [][]float64, a, b int) {
  var swap []float64
  swap = m[a]
  m[a] = m[b]
  m[b] = swap
}

//Invert the matrix m
//May return nil! 
func Inverse(m [][]float64) [][]float64 {
  det := Det(m)
  if det == 0 { return nil }
  dim := len(m)

  v := make([][]float64, dim)
  M := make([][]float64, dim)
  for i := 0; i < dim; i ++ {
    M[i] = make([]float64, dim)
    v[i] = make([]float64, dim)
    for j := 0; j < dim; j ++ {
      M[i][j] = m[i][j]
      if i == j {
        v[i][j] = 1
      } else {
        v[i][j] = 0
      }
    }
  }

  var x float64
  for i := 0; i < dim; i ++ {
    if M[i][i] == 0 {
      for k := i + 1; k < dim; k ++ {
        if M[k][i] != 0 {
          elementaryRowOperationSwap(M, i, k)
          elementaryRowOperationSwap(v, i, k)
          break
        }
      }
    }
    for j := 0; j < dim; j ++ {
      if i == j {
        x = 1 / M[i][i]
        elementaryRowOperationTimes(M, i, x)
        elementaryRowOperationTimes(v, i, x)
      } else {
        x = - M[j][i] / M[i][i]
        elementaryRowOperationPlus(M, i, j, x)
        elementaryRowOperationPlus(v, i, j, x)
      }
    }
  }

  return v
}

func MatrixMultiply(m [][]float64, v []float64) []float64 {
  dim := len(m)
  z := make([]float64, dim)

  for i := 0; i < dim; i ++ {
    z[i] = Dot(m[i], v)
  }

  return z
}
