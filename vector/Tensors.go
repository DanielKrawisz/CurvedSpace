package vector

import "sort"

//Functions for contracting tensors. These functions are inefficient and should
//not be used in the main loop of a program. Eventually more efficient versions
//will be created, however.
func ContractSymmetric4Tensor(t [][][][]float64, x []float64 ) [][][]float64 {
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

func ContractSymmetric3Tensor(t [][][]float64, x []float64 ) [][]float64 {
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
func ContractSymmetricTensor(t [][]float64, x []float64 ) []float64 {
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

func MatrixMultiplySymmetricTensor(m [][]float64, v [][]float64) [][]float64 {
  dim := len(v)
  z := make([][]float64, dim)
  for i := 0; i < dim; i ++ {
    z[i] = make([]float64, i + 1)

    for j := 0; j < i; j ++ {
      z[i][j] = 0

      for k := 0; k < dim; k ++ {
        for l := 0; l < k; l ++ {
          z[i][j] += v[k][l] * (m[i][k] * m[j][l] + m[i][l] * m[j][k])
        }

        //The case where k == l. 
        z[i][j] += v[k][k] * m[i][k] * m[j][k]
      }
    }

    //The case where i == j.
    z[i][i] = 0

    for k := 0; k < dim; k ++ {
      for l := 0; l < k; l ++ {
        z[i][i] += v[k][l] * (m[i][k] * m[i][l] + m[i][l] * m[i][k])
      }

      //The case where k == l. 
      z[i][i] += v[k][k] * m[i][k] * m[i][k]
    }
  }

  return z
}

func MatrixMultiplySymmetric3Tensor(m [][]float64, v [][][]float64) [][][]float64 {
  return nil // TODO
}

func MatrixMultiplySymmetric4Tensor(m [][]float64, v [][][][]float64) [][][][]float64 {
  return nil // TODO
}

//Add one tensor to another. 
func AddToSymmetricTensor(v, b [][]float64) {
  dim := len(b)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      v[i][j] += b[i][j]
    }
  }
}

func AddToSymmetric3Tensor(v, b [][][]float64) {
  dim := len(b)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      for k := 0; k <= j; k ++ {
        v[i][j][k] += b[i][j][k]
      }
    }
  }
}

func AddToSymmetric4Tensor(v, b [][][][]float64) {
  dim := len(b)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      for k := 0; k <= j; k ++ {
        for l := 0; l <= k; l ++ {
          v[i][j][k][l] += b[i][j][k][l]
        }
      }
    }
  }
}

//Some functions to multiply a tensor by a scalar. 
func SymmetricTensorTimes(a float64, v [][]float64) [][]float64 {
  dim := len(v)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      v[i][j] *= a
    }
  }

  return v
}

func SymmetricTensor3Times(a float64, v [][][]float64) [][][]float64 {
  dim := len(v)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      for k := 0; k <= j; k ++ {
        v[i][j][k] *= a
      }
    }
  }

  return v
}

func SymmetricTensor4Times(a float64, v [][][][]float64) [][][][]float64 {
  dim := len(v)
  for i := 0; i < dim; i ++ {
    for j := 0; j <= i; j ++ {
      for k := 0; k <= j; k ++ {
        for l := 0; l <= k; l ++ {
          v[i][j][k][l] *= a
        }
      }
    }
  }

  return v
}
