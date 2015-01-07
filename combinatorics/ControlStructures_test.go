package combinatorics

import "testing"
import "../test"

//The mock test iterator just counts the right number of iterations have taken place.
type mockTestIterator struct {
  n uint64
  m uint64
}

func (i *mockTestIterator) Iterate(index []uint, x int) {
  i.n ++
  i.m += uint64(x)
}

func TestNestedFor(t *testing.T) {
  //First test the degenerate case. 
  degenerate := [][]uint{nil, []uint{}}
  for _, limit := range degenerate {

    I := &mockTestIterator{0, 0}

    NestedFor(I, limit)
    if I.n != uint64(0) {
      t.Error("degenerate nested for error. Limit = ", limit)
    }
  }

  //A set of random cases.
  for a := 0; a < 4; a ++ {
    var l int = 3
    limit := make([]uint, l)
    for i := 0; i < l; i ++ {
      limit[i] = uint(test.RandInt(3, 6))
    }

    var exp uint = 1
    for i := 0; i < l; i ++ {
      exp *= limit[i]
    }

    I := &mockTestIterator{0, 0}

    NestedFor(I, limit)

    if I.n != uint64(exp) {
      t.Error("Nested for error. Limit = ", limit, "; exp = ", exp, "; n = ", I.n)
    }
  }
}

//TODO these next three tests do not ensure that all permutations
//fed to the iteration loop are correct.
func TestNestedForPermutations(t *testing.T) {
  for n := uint(1); n < 8; n ++ {
    I := &mockTestIterator{0, 0}
    NestedForPermutation(I, n)

    if I.n != factorial[n] {
      t.Error("Nested permutation for error: limit ", n, ", count ", I.n)
    }
  }
}

//This tests to make sure that the control structure 
//iterates the correct number of times. 
func TestNestedForAsymmetric(t *testing.T) {
  for rank := uint(0); rank < 6; rank ++ {
    for dim := uint(0); dim < 6; dim ++ {
      I := &mockTestIterator{0, 0}
      NestedForAsymmetric(I, rank, dim)

      var exp uint64
      if dim == 0 {
        exp = 0
      } else {
        exp = Binomial(dim, rank)
      }

      if I.n != exp {
        t.Error("Nested asymmetric permutation for error: rank ", 
          rank, ", dim ", dim, ", expected ", exp, " count ", I.n)
      }
    }
  }
}

func TestNestedForSymmetric(t *testing.T) {
  for rank := uint(0); rank < 6; rank ++ {
    for dim := uint(0); dim < 4; dim ++ {
      I := &mockTestIterator{0, 0}
      NestedForSymmetric(I, rank, dim)

      var exp_m int64
      var exp_n uint64
      if dim != 0 {
        exp_m = Power(int(dim), rank)
        exp_n = Figurate(rank, dim)
      } else {
        exp_m = 0
        exp_n = 0
      }

      if I.m != uint64(exp_m) {
        t.Error("Nested symmetric permutation for error type 1: rank ",
          rank, ", dim ", dim, ", expected ", exp_m, " count ", I.m)
      }
      if I.n != uint64(exp_n) {
        t.Error("Nested symmetric permutation for error type 2: rank ",
          rank, ", dim ", dim, ", expected ", exp_n, " count ", I.n)
      }
    }
  }
}
