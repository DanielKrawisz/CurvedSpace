package vector_test

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/test"
import "github.com/DanielKrawisz/CurvedSpace/vector"

var vec_err float64 = .000001

func TestVectorOperations(t *testing.T) {
  V := []float64{7, 4, 4}
  W := []float64{6, 6, 3}

  if !test.CloseEnough(vector.Dot(V, W), 78, vec_err) {
    t.Error("vector error 1")
  }
  if !test.CloseEnough(vector.Dot(V, V), 81, vec_err) {
    t.Error("vector error 2")
  }
  if !test.CloseEnough(vector.Length(V), 9, vec_err) {
    t.Error("vector error 3")
  }
  if !test.CloseEnough(vector.Length(W), 9, vec_err) {
    t.Error("vector error 4")
  }
  if !test.VectorCloseEnough(vector.Plus(V, W), []float64{13, 10, 7}, vec_err) {
    t.Error("vector error 5")
  }
  if !test.VectorCloseEnough(vector.Minus(V, W), []float64{1, -2, 1}, vec_err) {
    t.Error("vector error 6")
  }
  case7 := vector.LinearSum(2, 3, V, W)
  if !test.VectorCloseEnough(case7, []float64{32, 26, 17}, vec_err) {
    t.Error("vector error 7: ")
  }
  if !test.VectorCloseEnough(vector.Times(-4, V), []float64{-28, -16, -16}, vec_err) {
    t.Error("vector error 8: ")
  }
  if !test.VectorCloseEnough(vector.Negative(V), []float64{28, 16, 16}, vec_err) {
    t.Error("vector error 9: ")
  }
  case10 := vector.Normalize(V)
  if !test.VectorCloseEnough(case10, []float64{-0.77777777778, -0.44444444444, -0.44444444444}, vec_err) {
    t.Error("vector error 10: got ", case10)
  }
}

func TestOrthonormalize(t *testing.T) {
  //Some set trials.
  trials :=
    [][][][]float64{
      [][][]float64{[][]float64{
                      []float64{2, 0},
                      []float64{0, -4}},
                    [][]float64{
                      []float64{1, 0},
                      []float64{0, -1}}},
      [][][]float64{[][]float64{
                      []float64{1, 0},
                      []float64{2, 3}},
                    [][]float64{
                      []float64{1, 0},
                      []float64{0, 1}}},
      [][][]float64{[][]float64{
                      []float64{1, 0},
                      []float64{2, 0},
                      []float64{2, 1}},
                    [][]float64{
                      []float64{1, 0},
                      []float64{0, 1},
                      []float64{0, 0}}},
      [][][]float64{[][]float64{
                      []float64{1, 0, 0},
                      []float64{1, 1, 0},
                      []float64{2, 1, 0},
                      []float64{1, 2, 0},
                      []float64{1, 0, 1}},
                    [][]float64{
                      []float64{1, 0, 0},
                      []float64{0, 1, 0},
                      []float64{0, 0, 1},
                      []float64{0, 0, 0},
                      []float64{0, 0, 0}}}}

  for i, trial := range trials {
    got := vector.Orthonormalize(trial[0])
    if !test.MatrixCloseEnough(got, trial[1], vec_err) {
      t.Error("orthonormalize error trial ", i, ". expected ", trial[1], ", got ", got)
    }
  }

  //five random trials. 
  dim := 4
  for i := 0; i < 5; i ++ {
    o := test.RandFloatMatrix(-5, 5, dim, dim)

    vector.Orthonormalize(o)

    for j := 0; j < dim; j ++ {
      for k := 0; k < j ; k ++ {
        if !test.CloseEnough(vector.Dot(o[j], o[k]), 0, .000001) {
          t.Error("Orthonormalization error: ", o)
        }
      }

      if !test.CloseEnough(vector.Dot(o[j], o[j]), 1, .000001) {
        t.Error("Orthonormalization error.")
      }
    }
  }
}

func TestDot(t *testing.T) {
  cases := [][][]float64{
    [][]float64{[]float64{}, []float64{}},
    [][]float64{[]float64{1}, []float64{1}},
    [][]float64{[]float64{1}, []float64{2}},
    [][]float64{[]float64{3}, []float64{2}},
    [][]float64{[]float64{1, 1}, []float64{1, -1}},
    [][]float64{[]float64{2, 6}, []float64{4, 3}}}

  expected := []float64{0, 1, 2, 6, 0, 26}

  for i := 0; i < len(cases); i ++ {
    test_dot := vector.Dot(cases[i][0], cases[i][1])

    if !test.CloseEnough(test_dot, expected[i], .00001) {
      t.Error("Dot error! test case ", cases[i], "; expected ", expected[i], "; got ", test_dot)
    }
  }
}

func TestDet(t *testing.T) {
  cases := [][][]float64{
      [][]float64{},
      [][]float64{
          []float64{8}},
      [][]float64{
          []float64{2, 0, 0}, 
          []float64{0, 0, 1},
          []float64{0, 1, 0}}, 
      [][]float64{
          []float64{1, 1, 0}, 
          []float64{0, 1, 0},
          []float64{0, 0, 1}}, 
      [][]float64{
          []float64{1, 0, 0}, 
          []float64{0, 0, 0},
          []float64{0, 0, 1}}, 
      [][]float64{
          []float64{1, 0, 0}, 
          []float64{0, 1, 0},
          []float64{1, 0, 0}}, 
      [][]float64{
          []float64{0, 1, 0}, 
          []float64{1, 0, 0},
          []float64{0, 0, 1}},
      [][]float64{ 
          []float64{0.3734838781633063, 0, 0},
          []float64{0, 0.7177204175693372, 0}, 
          []float64{0, 0, 0.9782291177529912}}, 
      [][]float64{ 
          []float64{0.3734838781633063, -0.14242181016587718, 0.7832328557174395},
          []float64{0, 0.7177204175693372, -0.4183506835575348}, 
          []float64{0, 0, 0.9782291177529912}}, 
      [][]float64{ 
          []float64{0.3734838781633063, -0.14242181016587718, 0.7832328557174395},
          []float64{-0.43875143013729945, 0.7177204175693372, -0.4183506835575348}, 
          []float64{-0.9423920649659894, 0.08840190638381973, 0.9782291177529912}}}

  det_exp := []float64{0, 8, -2, 1, 0, 0, -1, 0.262221, 0.262221, 0.658136}

  for i := 0; i < len(cases); i ++ {
    det_test := vector.Det(cases[i])
    if !test.CloseEnough(det_test, det_exp[i], .00001) {
      t.Error("Det error! test case ", cases[i], "; expected ", det_exp[i], "; got ", det_test)
    }
  }
}

func TestCross(t *testing.T) {
  cases := [][][]float64{
      [][]float64{
          []float64{1}}, 
      [][]float64{
          []float64{1, 0},
          []float64{0, 1}}, 
      [][]float64{
          []float64{1, 0, 0}, 
          []float64{0, 1, 0}, 
          []float64{0, 0, 1}}, 
      [][]float64{
          []float64{1, 0, 0, 0}, 
          []float64{0, 1, 0, 0}, 
          []float64{0, 0, 1, 0}, 
          []float64{0, 0, 0, 1}}}

  for i := 0; i < len(cases); i ++ {
    cross := vector.Cross(cases[i][:len(cases[i]) - 1])
    if !test.VectorCloseEnough(cross, cases[i][len(cases[i]) - 1], .00001) {
      t.Error("Cross error: test case ", cases[i], ", got ", cross)
    }
  }
}

func TestTranspose(t *testing.T) {
  dim := 4
  for o := 0; o < 5; o ++ {
    m := make([][]float64, dim)
    for i := 0; i < dim; i ++ {
      m[i] = make([]float64, dim)
    }

    test_diag_ind := []int{test.RandInt(0, dim - 1), test.RandInt(0, dim - 1)}

    test_corner_ind := make([][]int, 3)
    for i := 0; i < 3; i ++ {
      var a, b int
      for {
        a = test.RandInt(0, dim - 1)
        b = test.RandInt(0, dim - 1)
        if a != b {
          break
        }
      }
      if a < b {
        test_corner_ind[i] = []int{b, a}
      } else {
        test_corner_ind[i] = []int{a, b}
      }
    }

    //Fill up the bottom half and diagonal with random values.
    for i := 0; i < dim; i ++ {
      for j := 0; j <= i; j ++ {
        m[i][j] = test.RandFloat(-10, 10)
      }
    }

    test_diag := make([]float64, 2)
    test_corner := make([]float64, 3)

    for i := 0; i < 2; i ++ {
      test_diag[i] = m[test_diag_ind[i]][test_diag_ind[i]]
    }
    for i := 0; i < 3; i ++ {
      test_corner[i] = m[test_corner_ind[i][0]][test_corner_ind[i][1]]
    }

    vector.Transpose(m)

    for i := 0; i < 2; i ++ {
      if m[test_diag_ind[i]][test_diag_ind[i]] != test_diag[i] {
        t.Error("Transpose error 0!")
      }
    }
    for i := 0; i < 3; i ++ {
      if m[test_corner_ind[i][0]][test_corner_ind[i][1]] != 0 {
        t.Error("Transpose error 1! ", test_corner_ind[i])
      }
      if m[test_corner_ind[i][1]][test_corner_ind[i][0]] != test_corner[i] {
        t.Error("Transpose error 2!")
      }
    }
  }
}

func TestInverse(t *testing.T) {
  cases := [][][]float64{
      [][]float64{
          []float64{2, 0, 0}, 
          []float64{0, 0, 1},
          []float64{0, 1, 0}}, 
      [][]float64{
          []float64{1, 1, 0}, 
          []float64{0, 1, 0},
          []float64{0, 0, 1}}, 
      [][]float64{
          []float64{1, 0, 0}, 
          []float64{0, 0, 0},
          []float64{0, 0, 1}}, 
      [][]float64{ 
          []float64{0.3734838781633063, -0.14242181016587718, 0.7832328557174395},
          []float64{-0.43875143013729945, 0.7177204175693372, -0.4183506835575348}, 
          []float64{-0.9423920649659894, 0.08840190638381973, 0.9782291177529912}}, 
      [][]float64{ 
          []float64{1.9612548100315585, 0.2216927362203317}, 
          []float64{-0.23968396947452258, 2.0228470649129853}}}

  expected := [][][]float64{
      [][]float64{
          []float64{.5, 0, 0}, 
          []float64{0, 0, 1},
          []float64{0, 1, 0}}, 
      [][]float64{
          []float64{1, -1, 0}, 
          []float64{0, 1, 0},
          []float64{0, 0, 1}}, 
      nil, 
      [][]float64{ 
          []float64{1.122986684105179, 0.3168956366712079, -0.7636110475042883},
          []float64{1.251184769390772, 1.6766521910844983, -0.284739458963841}, 
          []float64{0.9687777680217666, 0.15376835616213544, 0.31235066873550765}},
      [][]float64{
          []float64{0.5031388736442349, -0.05514120940319543},
          []float64{0.05961613437008792, 0.48781915012990057}}}

  for i := 0; i < len(cases); i ++ {
    inv_test := vector.Inverse(cases[i])
    if inv_test == nil {
      if expected[i] != nil {
        t.Error("inverse error, case ", i)
      }
    } else {
      if ! test.MatrixCloseEnough(inv_test, expected[i], .00001) {
        t.Error("inverse error, case ", i, "; expected ", expected[i], ", got ", inv_test)
      }
    }
  }
}

func TestMatrixMultiply(t *testing.T) {
  
}
