package pathtrace

import "testing"
import "github.com/DanielKrawisz/CurvedSpace/surface"
import "github.com/DanielKrawisz/CurvedSpace/test"

func TestGlowInteractor(t *testing.T) {
  if NewGlowingObject(nil ) != nil { t.Error("glow error 1") }

  glow := NewGlowingObject([]float64{.5, .7, .9})
  if glow == nil {
    t.Error("glow error 2")
    return
  }

  ray := &LightRay{0, []float64{}, []float64{}, []float64{4, 5, 6}, []float64{1, 1, 1}, []float64{1, 1, 1}, 1}
  glow.Interact(ray)
  if !(test.VectorCloseEnough(ray.color, []float64{1, 1, 1}, mat_err) && 
       test.VectorCloseEnough(ray.emission, []float64{1.5, 1.7, 1.9}, mat_err) && 
       test.CloseEnough(ray.redirected, 0, mat_err)) {
    t.Error("glow error 3: ", ray)
  }

  //TODO can't test this yet.
  /*ray = &LightRay{0, []float64{}, []float64{}, []float64{1, 1, 1}, []float64{1, 1, 1}, 1}
  color := glow.Trace(nil, ray)
  if !test.VectorCloseEnough([]float64{.5, .7, .9}, color, mat_err) {
    t.Error("glow error 4")
  }*/
}

//We already know that lambertian and mirror reflection work,
//so these just test whether the objects are created correctly. 
func TestLambertianReflector(t *testing.T) {
  if NewLambertianReflector(nil, Absorb([]float64{0, 0, 0})) != nil {
    t.Error("lambertian error 1") }
  if NewLambertianReflector(surface.NewSphere([]float64{0, 0, 0}, 1), nil) != nil {
    t.Error("lambertian error 2") }
}

func TestMirrorReflector(t *testing.T) {
  if NewMirrorReflector(nil, Absorb([]float64{0, 0, 0})) != nil {
    t.Error("lambertian error 1") }
  if NewMirrorReflector(surface.NewSphere([]float64{0, 0, 0}, 1), nil) != nil {
    t.Error("lambertian error 2") }
}

func TestRedirectorReflector(t *testing.T) {
  if NewRedirectorInteractor(nil, Absorb([]float64{0, 0, 0}), MirrorReflection) != nil {
    t.Error("redirector error 1") }
  if NewRedirectorInteractor(surface.NewSphere([]float64{0, 0, 0}, 1), nil, MirrorReflection) != nil {
    t.Error("redirector error 2") }
  if NewRedirectorInteractor(surface.NewSphere([]float64{0, 0, 0}, 1), Absorb([]float64{0, 0, 0}), nil) != nil {
    t.Error("redirector error 3") }
}

//These next two we can test by mocking. 
func TestMultipleReflector(t *testing.T) {
  //TODO
}
