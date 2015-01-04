package combinatorics

import "testing"
import "../test"

//The mock test iterator just counts the right number of iterations have taken place.
type mockTestIterator struct {
  n int
}

func (i *mockTestIterator) Iterate(index []int, x int) {
  i.n ++
}

func TestNestedFor(t *testing.T) {
  var l int = 3
  limit := make([]int, l)

  for a := 0; a < 4; a ++ {
    for i := 0; i < l; i ++ {
      limit[i] = test.RandInt(3, 6)
    }

    var exp int = 1
    for i := 0; i < l; i ++ {
      exp *= limit[i]
    }

    I := &mockTestIterator{0}

    NestedFor(I, limit)

    if I.n != exp {
      t.Error("Nested for error. Limit = ", limit, "; exp = ", exp, "; n = ", I.n)
    }
  }
}

//TODO this test does not ensure that all permutations fed to
//the iteration loop are correct.
func TestNestedForPermutations(t *testing.T) {
  I := &mockTestIterator{0}

  for n := 3; n < 10; n ++ {
    I.n = 0
    NestedForPermutation(I, n)

    if I.n != factorial[n] {
      t.Error("Nested permutation for error: limit ", n, ", count ", I.n)
    }
  }
}
