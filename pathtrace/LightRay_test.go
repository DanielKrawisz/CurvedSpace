package pathtrace

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/vector"
import "github.com/DanielKrawisz/CurvedSpace/test"

var mat_err float64 = .00001

//Just testing some particular conditions for each of the functions. 

func TestDeriveColor(t *testing.T) {
  color := []float64{.1, .2, .3}
  emission := []float64{.4, .5, .6}
  redirected := .7
  ray := &LightRay{0, []float64{}, []float64{}, []float64{4, 5, 6}, color, emission, redirected}

  c := ray.DeriveColor()

  if !test.VectorCloseEnough(c, vector.LinearSum(1, redirected, emission, color), mat_err) {
    t.Error("derive color error")
  }
}

func TestGlow(t *testing.T) {
  rec := []float64{.1, .2, .3}
  em  := []float64{.4, .5, .6}
  red := .7
  ray := &LightRay{0, []float64{}, []float64{}, []float64{4, 5, 6}, rec, em, red}

  ray.Glow([]float64{.4, .7, .9})

  if !test.VectorCloseEnough(rec, []float64{.1, .2, .3}, .00001) {
    t.Error("Glow error 2")
  }
  if !test.VectorCloseEnough(em, []float64{0.428, 0.598, 0.789}, .00001) {
    t.Error("Glow error 3")
  }
  if !test.CloseEnough(ray.redirected, 0, .00001) {
    t.Error("Glow error 4")
  }
}

func TestAbsorb(t *testing.T) {
  rec := []float64{.1, .2, .3}
  em  := []float64{.4, .5, .6}
  red := .7
  ray := &LightRay{0, []float64{}, []float64{}, []float64{4, 5, 6}, rec, em, red}

  ray.Absorb([]float64{.4, .7, .9})

  if !test.VectorCloseEnough(rec, []float64{0.04, 0.14, 0.27}, .00001) {
    t.Error("Absorb error 2; got ", rec)
  }
  if !test.VectorCloseEnough(em, []float64{.4, .5, .6}, .00001) {
    t.Error("Absorb error 3")
  }
  if !test.CloseEnough(ray.redirected, .7, .00001) {
    t.Error("Absorb error 4")
  }
}

func TestGlowAbsorbAverage(t *testing.T) {

  rec := []float64{.1, .2, .3}
  em  := []float64{.4, .5, .6}
  red := .7
  ray := &LightRay{0, []float64{}, []float64{}, []float64{4, 5, 6}, rec, em, red}

  ray.GlowAbsorbAverage([]float64{.4, .7, .9}, []float64{.5, .6, .8}, .3)

  if !test.VectorCloseEnough(rec, []float64{0.05, 0.12, 0.24}, .00001) {
    t.Error("Average error 3: got ", rec)
  }
  if !test.VectorCloseEnough(em, []float64{0.484, 0.647, 0.789}, .00001) {
    t.Error("Average error 4")
  }
  if !test.CloseEnough(ray.redirected, .49, .00001) {
    t.Error("Average error 5")
  }
}
