package vector

import "testing"
import "../test"

func TestVectorOperations(t *testing.T) {
  //TODO
}

func TestOrthonormalize(t *testing.T) {
  //TODO
}

func TestDet(t *testing.T) {
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
          []float64{1, 0, 0}, 
          []float64{0, 1, 0},
          []float64{1, 0, 0}}, 
      [][]float64{
          []float64{0, 1, 0}, 
          []float64{1, 0, 0},
          []float64{0, 0, 1}}}

  det_exp := []float64{-2, 1, 0, 0, -1}

  for i := 0; i < len(cases); i ++ {
    det_test := Det(cases[i])
    if !test.CloseEnough(det_test, det_exp[i], .00001) {
      t.Error("Det error! test case ", cases[i], "; expected ", det_exp[i], "; got ", det_test)
    }
  }
}

func TestCross(t *testing.T) {
  cases := [][][]float64{
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
    cross := Cross(cases[i][:len(cases[i]) -1])
    if !test.VectorCloseEnough(cross, cases[i][len(cases[i]) - 1], .00001) {
      t.Error("Cross error: test case ", cases[i], ", got ", cross)
    }
  }
}
