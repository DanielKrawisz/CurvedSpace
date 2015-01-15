package pathtrace

import "testing"
import "../vector"
import "../test"

var mat_err float64 = .00001

//Just testing some particular conditions for each of the functions. 

func TestDeriveColor(t *testing.T) {
  receptor := []float64{.1, .2, .3}
  emission := []float64{.4, .5, .6}
  redirected := .7
  c := DeriveColor(receptor, emission, redirected)

  if !test.VectorCloseEnough(c, vector.LinearSum(1, redirected, emission, receptor), mat_err) {
    t.Error("derive color error")
  }
}

func TestGlow(t *testing.T) {
  if Glow(nil) != nil {
    t.Error("Glow error 1")
  }
  rec_in := []float64{.1, .2, .3}
  em_in  := []float64{.4, .5, .6}
  red_in := .7

  glow := Glow([]float64{.4, .7, .9})

  rec_out, em_out, red_out := glow(rec_in, em_in, red_in)
  if !test.VectorCloseEnough(rec_out, []float64{.1, .2, .3}, .00001) {
    t.Error("Glow error 2")
  }
  if !test.VectorCloseEnough(em_out, []float64{0.428, 0.598, 0.789}, .00001) {
    t.Error("Glow error 3")
  }
  if !test.CloseEnough(red_out, 0, .00001) {
    t.Error("Glow error 4")
  }
}

func TestAbsorb(t *testing.T) {
  if Absorb(nil) != nil {
    t.Error("Absorb error 1")
  }
  rec_in := []float64{.1, .2, .3}
  em_in  := []float64{.4, .5, .6}
  red_in := .7

  absorb := Absorb([]float64{.4, .7, .9})

  rec_out, em_out, red_out := absorb(rec_in, em_in, red_in)
  if !test.VectorCloseEnough(rec_out, []float64{0.04, 0.14, 0.27}, .00001) {
    t.Error("Absorb error 2")
  }
  if !test.VectorCloseEnough(em_out, []float64{.4, .5, .6}, .00001) {
    t.Error("Absorb error 3")
  }
  if !test.CloseEnough(red_out, .7, .00001) {
    t.Error("Absorb error 4")
  }
}

func TestGlowAbsorbAverage(t *testing.T) {
  if GlowAbsorbAverage(nil, []float64{0, 0}, 1) != nil {
    t.Error("Glow absorb average error 1")
  }
  if GlowAbsorbAverage([]float64{0, 0}, nil, 1) != nil {
    t.Error("Glow absorb average error 2")
  }

  rec_in := []float64{.1, .2, .3}
  em_in  := []float64{.4, .5, .6}
  red_in := .7

  average := GlowAbsorbAverage([]float64{.4, .7, .9}, []float64{.5, .6, .8}, .3)

  rec_out, em_out, red_out := average(rec_in, em_in, red_in)
  if !test.VectorCloseEnough(rec_out, []float64{0.05, 0.12, 0.24}, .00001) {
    t.Error("Average error 3: got ", rec_out)
  }
  if !test.VectorCloseEnough(em_out, []float64{0.484, 0.647, 0.789}, .00001) {
    t.Error("Average error 4")
  }
  if !test.CloseEnough(red_out, .49, .00001) {
    t.Error("Average error 5")
  }
}
