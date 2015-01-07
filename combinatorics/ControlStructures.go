package combinatorics

//Functions for running nested loops. 

//An inner loop to a nested for control structure.
type IterationLoop interface {
  Iterate([]uint, int)
}

//The nested for iterates over any number of indices with the given limits.
//The set of indices is fed to the iterate object with each iteration. 
func NestedFor(i IterationLoop, limit []uint) {
  if limit == nil { return }
  dim := len(limit)
  if dim == 0 { return }
  index := make([]uint, dim)

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

//The permutation loop helps iterate over permutations of numbers. 
type permutationIterationLoop struct {
  i IterationLoop 
  permutation []uint
}

func (al *permutationIterationLoop) Iterate(index []uint, x int) {
  var swap uint
  var dim uint = uint(len(index)) + 1

  var i uint
  //Construct the identity permutation.
  for i = 0; i < dim; i ++ {
    al.permutation[i] = i
  }

  var signature int = 0
  //This loop constructs the permutation from the
  //list of indices. 
  for k, j := range index {
    if index[k] != 0 {
      signature ++
    }

    swap = al.permutation[k]
    al.permutation[k] = al.permutation[j + uint(k)]
    al.permutation[j + uint(k)] = swap
  }

  al.i.Iterate(al.permutation, -2 * (signature & 1) + 1)
}

//Iterates over the permutations in the numbers from 0 to dim - 1.
func NestedForPermutation(i IterationLoop, dim uint) {
  if dim == 0 { return }
  //Special degenerate case. 
  if dim == 1 {
    i.Iterate([]uint{0}, 1)
    return
  }

  //limit is one shorter than dim so that we can skip
  //iterating over a loop that only goes up to 1. It's
  //automatically included. However, that excludes the 
  //degenerate case in which dim == 1, so we handle that
  //separately. 
  limit := make([]uint, dim - 1) 

  for i := uint(0); i < dim - 1; i ++ {
    limit[i] = dim - i
  }

  il := &permutationIterationLoop{i, make([]uint, dim)}

  NestedFor(il, limit) 
}

//Iterates over all ordered lists of length rank that contain
//the numbers 0 - dim with no repetition. Obviously, if
//rank > dim, then there is no iteration. 
func NestedForAsymmetric(il IterationLoop, rank, dim uint) {
  if rank > dim { return }
  if dim == 0 { return }
  //Special degenerate case. 
  if rank == 0 {
    il.Iterate([]uint{}, 0)
    return
  }

  index := make([]uint, rank)
  limit := make([]uint, rank)

  var i, in uint

  for i = 0; i < rank; i ++ {
    limit[i] = i + 1 + dim - rank
    index[i] = i
  }

  for {
    il.Iterate(index, 0)

    in = rank - 1
    for {
      index[in] ++
      if index[in] < limit[in] {
        for in ++; in < rank; in ++ {
          index[in] = index[in - 1] + 1
        }
        break
      } else {
        if in == 0 {
          return
        }
        in --
      }
    }
  }
}

//Iterates over all ordered lists of length rank that contain
//the numbers 0 - dim with repetition. 
func NestedForSymmetric(il IterationLoop, rank, dim uint) {
  if dim == 0 { return }
  //Special degenerate case. 
  if rank == 0 {
    il.Iterate([]uint{}, 1)
    return
  }
  index := make([]uint, rank)

  var i uint
  for i = 0; i < rank; i ++ {
    index[i] = 0
  }

  var in uint

  for {
    il.Iterate(index, int(Permutations(index)))

    //This adjusts the indices.
    in = rank - 1
    for {
      index[in] ++
      if index[in] < dim {
        for in ++; in < rank; in ++ {
          index[in] = index[in - 1]
        }
        break
      } else {
        if in == 0 {
          return
        }
        in --
      }
    }
  }
}
