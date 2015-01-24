package pathtrace

import "math/rand"
import "../surface"

//This might be updated to be more of an interface or whatever. 
type LightRay struct {
  //The number of steps the ray has taken.
  depth int
  //The ray. 
  position, direction []float64
  //The wavelengths of light that this ray tracks.
  receptor []float64
  //The intensities of those wavelengths. 
  color []float64
  //As the ray bounces along, if it encounters glowing objects,
  //these numbers keep track of how the final color should be
  //adjusted to take these earlier interactions into account.
  emission []float64
  redirected float64 
}

func (r *LightRay) Trace(u float64) {
  for i := 0; i < 3; i ++ {
    r.position[i] = r.position[i] + u * r.direction[i]
  }
}

//Something that interacts with a ray.
type Interactor interface {
  //Interact with the object and return a new ray. 
  Interact(ray *LightRay) *LightRay
  //Some kinds of interactions are more efficient if an object can send out several rays.
  //This can reduce the variance of a pixel. 
  //Trace(scene *Scene, ray *LightRay) []float64 TODO
}

//An object that just glows.
type glowEmitter struct {
  glow ColorInteraction
}

func (g *glowEmitter) Interact(ray *LightRay) *LightRay {
  ray.color, ray.emission, ray.redirected = g.glow(ray.color, ray.emission, ray.redirected)
  return ray
}

func (g *glowEmitter) Trace(scene *Scene, ray *LightRay) []float64 {
  return DeriveColor(g.glow(ray.color, ray.emission, ray.redirected))
}

//May return nil.
func NewGlowingObject(color []float64) Interactor {
  if color == nil { return nil }
  return &glowEmitter{Glow(color)}
}

type lambertianReflector struct {
  surf surface.Surface
  color ColorInteraction
}

func (l *lambertianReflector) Interact(ray *LightRay) *LightRay {
  ray.color, ray.emission, ray.redirected = l.color(ray.color, ray.emission, ray.redirected)
  ray.direction = LambertianReflection(ray.direction, surface.SurfaceNormal(l.surf, ray.position))
  return ray
}

/*func (l *lambertianReflector) Trace(scene *Scene, ray *LightRay, u float64) []float64 {
  return DeriveColor(g.glow(ray.color, ray.emission, ray.redirected))
}*/

//May return nil.
func NewLambertianReflector(surf surface.Surface, color ColorInteraction) Interactor {
  if surf == nil || color == nil {return nil}
  return &lambertianReflector{surf, color}
}

type mirrorReflector struct {
  surf surface.Surface
  color ColorInteraction
}

func (l *mirrorReflector) Interact(ray *LightRay) *LightRay {
  ray.color, ray.emission, ray.redirected = l.color(ray.color, ray.emission, ray.redirected)
  ray.direction = MirrorReflection(ray.direction, surface.SurfaceNormal(l.surf, ray.position))
  return ray
}

//May return nil.
func NewMirrorReflector(surf surface.Surface, color ColorInteraction) Interactor {
  if surf == nil || color == nil {return nil}
  return &mirrorReflector{surf, color}
}

type redirectorInteractor struct {
  surf surface.Surface
  color ColorInteraction
  redirect Redirection
}

func (l *redirectorInteractor) Interact(ray *LightRay) *LightRay {
  ray.color, ray.emission, ray.redirected = l.color(ray.color, ray.emission, ray.redirected)
  ray.direction = l.redirect(ray.direction, surface.SurfaceNormal(l.surf, ray.position))
  return ray
}

//May return nil
func NewRedirectorInteractor(surf surface.Surface, color ColorInteraction, redirect Redirection) Interactor {
  if surf == nil || color == nil || redirect == nil {return nil}
  return &redirectorInteractor{surf, color, redirect}
}

//May return nil.
func NewBasicRefractiveTransmitor(surf surface.Surface, color ColorInteraction, index float64) Interactor {
  return NewRedirectorInteractor(surf, color, BasicRefraction(index))
}

//May return nil.
func NewSpecularReflector(surf surface.Surface, color ColorInteraction, scatter float64) Interactor {
  return NewRedirectorInteractor(surf, color, SpecularReflection(scatter))
}

type scatterInteractor struct {
  color ColorInteraction
  redirect Redirection
}

func (l *scatterInteractor) Interact(ray *LightRay) *LightRay {
  ray.color, ray.emission, ray.redirected = l.color(ray.color, ray.emission, ray.redirected)
  ray.direction = l.redirect(ray.direction, nil)
  return ray
}

//May return nil.
func NewScatterTransmitter(color ColorInteraction, degree float64) Interactor {
  if color == nil {return nil}
  return &scatterInteractor{color, ScatterRedirector(degree)}
}

type multipleInteractor struct {
  surf surface.Surface
  color ColorInteraction
  //The strengths of the various available interactions.
  probabilities []float64
  factors []float64
  redirects []Redirection
}

func (s *multipleInteractor) Interact(ray *LightRay) *LightRay {
  spin := rand.Float64()
  ray.color, ray.emission, ray.redirected = s.color(ray.color, ray.emission, ray.redirected)

  for i, p := range s.probabilities {
    if spin < p {
      for j := 0; j < 3; j ++ {
        ray.color[j] *= s.factors[i]
      }
      ray.direction = s.redirects[i](ray.direction, surface.SurfaceNormal(s.surf, ray.position))
      break 
    } else {
      spin -= p
    }
  }

  return ray
}

/*func (s *surfaceInteractor) Trace(scene *Scene, ray *LightRay, u float64) []float64 {
  ray.color, ray.emission, ray.redirected = s.color(ray.color, ray.emission, ray.redirected)

  
}*/

//May return nil.
func NewMultipleInteractor(surf surface.Surface, color ColorInteraction,
  probabilities []float64, factors []float64, redirects []Redirection) Interactor {

  if surf == nil || color == nil || probabilities == nil || redirects == nil || factors == nil {
    return nil
  }

  if len(probabilities) != len(redirects) || len(factors) != len(redirects) {
    return nil
  }

  return &multipleInteractor{surf, color, probabilities, factors, redirects}
}

func NewShineyInteractor(surf surface.Surface, color ColorInteraction, p, scatter float64) Interactor {
  return NewMultipleInteractor(surf, color,  []float64{1 - p, p}, []float64{1, 1}, 
    []Redirection{SpecularReflection(scatter), LambertianReflection})
}

func NewGlassInteractor(surf surface.Surface, color ColorInteraction, index, a, b float64) Interactor {
  ab := a + b
  return NewMultipleInteractor(surf, color,  []float64{a / ab, b / ab}, []float64{ab, ab}, 
    []Redirection{BasicRefraction(index), MirrorReflection})
}
