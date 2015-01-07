package combinatorics

import "testing"

func TestFactorial(t *testing.T) {
  var i uint
  var exp uint64
  for i = 0; i < 25; i++ {
    f := Factorial(i)

    if i > 20 {
      exp = 0
    } else {
      exp = factorial[i]
    }

    if f != exp {
      t.Error("factorial error for input ", i, ". Expected ", exp, " got ", f)
    }
  }
}

func TestBinomial(t *testing.T) {
  test_cases := [][]uint{
                  []uint{19, 11}, []uint{10, 2}, []uint{4, 1}, []uint{4, 14}, []uint{10, 6},
                  []uint{3, 11}, []uint{18, 15}, []uint{8, 4}, []uint{6, 5}, []uint{20, 5},
                  []uint{20, 20}, []uint{7, 5}, []uint{16, 8}, []uint{19, 20}, []uint{16, 11},
                  []uint{2, 0}, []uint{1, 2}, []uint{6, 4}, []uint{7, 8}, []uint{3, 10},
                  []uint{50, 47}, []uint{28, 9}, []uint{37, 16}, []uint{59, 7}, []uint{66, 29},
                  []uint{28, 26}, []uint{33, 10}, []uint{33, 16}, []uint{64, 2}, []uint{28, 4},
                  []uint{50, 1}, []uint{58, 15}, []uint{31, 2}, []uint{30, 10}, []uint{27, 0},
                  []uint{28, 17}, []uint{49, 38}, []uint{31, 0}, []uint{30, 20}, []uint{33, 22},
                  []uint{70, 30}, []uint{62, 31}}

  expected := []uint64{75582, 45, 4, 0, 210, 0, 816, 70, 6, 15504, 1, 21, 12870, 0, 4368, 1, 0, 15, 0, 0,
     19600, 6906900, 12875774670, 341149446, 0, 378, 92561040, 1166803110, 0, 20475,
     50, 29752626158640, 465, 30045015, 1, 21474180, 29135916264, 1, 30045015, 193536720, 0,
     465428353255261088}

  for i, test_case := range test_cases {
    exp := expected[i] 

    test := Binomial(test_case[0], test_case[1])

    if exp != test {
      t.Error("Binomial error: input ", test_case, " expected ", exp, " got ", test)
    }
  }
}

func TestFigurate(t *testing.T) {
  test_cases := [][]uint{
                  []uint{14, 28}, []uint{23, 13}, []uint{28, 28}, []uint{10, 16}, []uint{1, 11},
                  []uint{20, 24}, []uint{21, 30}, []uint{26, 11}, []uint{21, 27}, []uint{9, 25},
                  []uint{8, 18}, []uint{1, 24}, []uint{28, 28}, []uint{28, 27}, []uint{22, 4},
                  []uint{4, 4}, []uint{1, 2}, []uint{20, 8}, []uint{10, 25}, []uint{15, 16}, 
                  []uint{0, 4}, []uint{3, 0}}

  expected := []uint64{35240152720, 834451800, 3824345300380220, 3268760, 11, 960566918220, 67327446062800,
                254186856, 12551759587422, 38567100, 1081575, 24, 3824345300380220, 1877405874732108, 2300,
                35, 2, 888030, 131128140, 155117520, 1, 0}

  for i, test_case := range test_cases {
    exp := expected[i] 

    test := Figurate(test_case[0], test_case[1])

    if exp != test {
      t.Error("Figurate error: input ", test_case, " expected ", exp, " got ", test)
    }
  }
}

func TestPermutations(t *testing.T) {
  test_cases := [][]uint{
                  []uint{1, 2, 3, 4, 5}, []uint{1, 2, 3, 4, 4}, []uint{1, 2, 2, 4, 4}, 
                  []uint{1, 2, 3, 3, 3}, []uint{1, 1, 4, 4, 4}, []uint{1, 7, 7, 7, 7}, 
                  []uint{8, 8, 8, 8, 8}}

  expected := []uint64{120, 60, 30, 20, 10, 5, 1}

  for i, test_case := range test_cases {
    exp := expected[i] 

    test := Permutations(test_case)

    if exp != test {
      t.Error("Figurate error: input ", test_case, " expected ", exp, " got ", test)
    }
  }
}

func TestPower(t *testing.T) {
  test_cases := [][]int{
                  []int{1, 6}, []int{-2, 10}, []int{2, 0}, []int{7, 8}, []int{2, 12},
                  []int{-1, 14}, []int{-2, 14}, []int{10, 14}, []int{-5, 1}, []int{-10, 5},
                  []int{-6, 5}, []int{8, 9}, []int{-5, 8}, []int{5, 4}, []int{-9, 2},
                  []int{9, 5}, []int{-9, 10}, []int{-10, 3}, []int{0, 4}, []int{4, 7}}

  expected := []int64{1, 1024, 1, 5764801, 4096, 1, 16384, 100000000000000, -5,
                -100000, -7776, 134217728, 390625, 625, 81, 59049, 3486784401, -1000, 0, 16384}

  for i, test_case := range test_cases {
    exp := expected[i] 

    test := Power(test_case[0], uint(test_case[1]))

    if exp != test {
      t.Error("Binomial error: input ", test_case, " expected ", exp, " got ", test)
    }
  }
}
