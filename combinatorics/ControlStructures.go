package combinatorics

//Functions for running nested loops. 

//The first 20 factorials. 
var factorial []int = []int{
  1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800, 39916800, 
  479001600, 6227020800, 87178291200, 1307674368000, 20922789888000, 
  355687428096000, 6402373705728000, 121645100408832000, 2432902008176640000}

type IterationLoop interface {
  Iterate([]int, int)
}

func NestedFor(i IterationLoop, limit []int) {
  dim := len(limit)
  index := make([]int, dim)

  var in int

  for {
    i.Iterate(index, 0)

    in = 0
    for {
      index[in] ++
      if index[in] == limit[in] {
        index[in] = 0
      } else {
        break
      }
      in ++
      if in == len(limit) {
        return
      }
    }
  }
}

type permutationIterationLoop struct {
  i IterationLoop 
  permutation []int
  signature int
}

func (al *permutationIterationLoop) Iterate(index []int, x int) {
  var swap int
  dim := len(index) + 1

  for i := 0; i < dim; i ++ {
    al.permutation[i] = i
  }

  for i, j := range index {
    swap = al.permutation[i]
    al.permutation[i] = al.permutation[j + i]
    al.permutation[j + i] = swap
  }

  al.signature = -1 * al.signature

  al.i.Iterate(al.permutation, al.signature)
}

func NestedForPermutation(i IterationLoop, dim int) {
  limit := make([]int, dim - 1) 

  for i := dim - 2; i >= 0; i -- {
    limit[i] = dim - i
  }

  il := &permutationIterationLoop{i, make([]int, dim), -1}

  NestedFor(il, limit) 
}

//TODO permute over asymmetric permutations and symmetric ones. 

/*func NestedForSymmetric(IterationLoop(), dim int) {
  index := make([]float64, dim)
  index := make([]float64, dim)

  var in int

  for {
    Iterate(index, 0)

    in = 0
    for {
      index[in] ++
      if index[in] == limit[in] {
        index[in] = 0
      } else {
        break
      }
      in ++
      if in == len(limit) {
        return
      }
    }
  }
}*/
