package pathtrace

import "math/rand"
import "github.com/DanielKrawisz/CurvedSpace/surface"

type ColorInteraction func(ray *LightRay) 

//A function to make an object glow.
func Glow(c []float64) ColorInteraction {
  if c == nil {
    return nil
  } else {
    return func(ray *LightRay) {
      ray.Glow(c)
    }
  }
}

//A function to make an object absorb light.
func Absorb(c []float64) ColorInteraction {
  if c == nil {
    return nil
  } else {
    return func(ray *LightRay) {
      ray.Absorb(c)
    }
  }
}

//A function for both.
func GlowAbsorbAverage(glow_color, transmit_color []float64, absorb float64) ColorInteraction {
  if glow_color == nil || transmit_color == nil {
    return nil
  } else {
    return func(ray *LightRay) {
      ray.GlowAbsorbAverage(glow_color, transmit_color, absorb)
    }
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
  glow []float64
}

func (g *glowEmitter) Interact(ray *LightRay) *LightRay {
  ray.Glow(g.glow)
  return ray
}

func (g *glowEmitter) Trace(scene *Scene, ray *LightRay) []float64 {
  ray.Glow(g.glow)
  return ray.DeriveColor()
}

//May return nil.
func NewGlowingObject(color []float64) Interactor {
  if color == nil { return nil }
  return &glowEmitter{color}
}

type lambertianReflector struct {
  surf surface.Surface
  color ColorInteraction
}

func (l *lambertianReflector) Interact(ray *LightRay) *LightRay {
  l.color(ray)
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
  l.color(ray)
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
  l.color(ray)
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
  l.color(ray)
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
  s.color(ray)

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
