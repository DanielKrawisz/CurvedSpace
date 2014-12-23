package diffeq

import "testing"
import "../test"

var err_cd float64 = .000001

func TestCompoundDerivative(t *testing.T) {
  d1 := NewCompoundDerivative(2*6)
  d2 := NewCompoundDerivative(2*9)
  x1 := make([]float64, 2*6)
  x2 := make([]float64, 2*9)
  v1 := make([]float64, 2*6)
  v2 := make([]float64, 2*9)
  vt1 := make([]float64, 2*6)
  vt2 := make([]float64, 2*6)
  vt3 := make([]float64, 2*9)

  if d1 == nil || d2 == nil {
    t.Error("nil compound derivative")
    return
  }

  ipos1 := [][]float64{[]float64{0,0,0}, []float64{0,0,0}}
  ipos2 := [][]float64{[]float64{0,0,0}, []float64{0,0,0}, []float64{0,0,0}}

  p1 := NewNewtonianParticle(3, 2, 0.0, 0.01, ipos1, [][]float64{[]float64{0,0,0}, []float64{0,0,0}})
  p2 := NewNewtonianParticle(3, 3, 0.0, 0.01, ipos2, 
    [][]float64{[]float64{0,0,0}, []float64{0,0,0}, []float64{0,0,0}})

  o1 := NewHarmonicOscillator(p1, []float64{1, 2}, [][][]float64{nil, [][]float64{[]float64{5, .125}}})
  o2 := NewHarmonicOscillator(p1, []float64{3, 3}, [][][]float64{nil, [][]float64{[]float64{1}}})
  o3 := NewHarmonicOscillator(p2, []float64{1, 2, 3},
    [][][]float64{nil , [][]float64{[]float64{1}}, [][]float64{[]float64{1}, []float64{1}}})

  //Test that an empty compound derivative returns zero. 
  if d1.Length() != 0 {
    t.Error("Compound derivative error 0")
  }
  if d2.Length() != 0 {
    t.Error("Compound derivative error 1")
  }

  for i := 0; i < 12; i ++ {
    v1[i] = test.RandFloat(-1, 1)
  }
  d1.DxDs(x1, v1)
  for i := 0; i < 12; i ++ {
    if !test.CloseEnough(v1[i], 0, err_cd) {
      t.Error("Compound derivative error 2")
      break
    }
  }

  //Test that the given derivatives will only fit in the right compound derivatives. 
  d1.Plus(o1)
  d2.Plus(o3)
  if d1.Length() != 1 {
    t.Error("Compound derivative error 3")
  }
  if d2.Length() != 1 {
    t.Error("Compound derivative error 4")
  }

  d1.Plus(o3)
  d2.Plus(o1)
  if d1.Length() != 1 {
    t.Error("Compound derivative error 5")
  }
  if d2.Length() != 1 {
    t.Error("Compound derivative error 6")
  }

  //Test that the derivatives return the correct values.
  o1.DxDs(x1, vt1)
  o3.DxDs(x2, vt3)
  d1.DxDs(x1, v1)
  d2.DxDs(x2, v2)

  for i := 0; i < 12; i ++ {
    if !test.CloseEnough(v1[i], vt1[i], err_cd) {
      t.Error("Compound derivative error 7")
      break
    }
  }
  for i := 0; i < 18; i ++ {
    if !test.CloseEnough(v2[i], vt3[i], err_cd) {
      t.Error("Compound derivative error 8")
      break
    }
  }

  //Test that the compound derivative returns the sum of the two derivatives in it. 
  d1.Plus(o2)
  if d1.Length() != 2 {
    t.Error("Compound derivative error 9")
  }
  o1.DxDs(x1, vt1)
  o2.DxDs(x2, vt2)
  d1.DxDs(x1, v1)

  for i := 0; i < 6; i ++ {
    if !test.CloseEnough(v1[i], vt1[i] + vt2[i], err_cd) {
      t.Error("Compound derivative error 10")
      break
    }
  }

  //Test that the derivative is removed when you add it again. 
  d1.Plus(o2)
  if d1.Length() != 1 {
    t.Error("Compound derivative error 11")
  }

  o1.DxDs(x1, vt1)
  d1.DxDs(x1, v1)

  for i := 0; i < 6; i ++ {
    if !test.CloseEnough(v1[i], vt1[i], err_cd) {
      t.Error("Compound derivative error 12")
      break
    }
  }
}

