package pathtrace

import "math"
import "github.com/DanielKrawisz/CurvedSpace/surface"
import "github.com/DanielKrawisz/CurvedSpace/color"

type InteractionFunction func([]float64) Interactor

//An object that exists in a scene. 
type ExtendedObject struct {
  surf surface.Surface
  interactor InteractionFunction
}

func NewExtendedObject(surf surface.Surface, interactor Interactor) *ExtendedObject {
  if surf == nil || interactor == nil { return nil }

  return &ExtendedObject{surf, func ([]float64) Interactor { return interactor }}
}

func NewTexturedExtendedObject(surf surface.Surface, interactor InteractionFunction) *ExtendedObject {
  if surf == nil || interactor == nil { return nil }

  return &ExtendedObject{surf, interactor}
}

//A set of objects of which a picture can be taken. 
//TODO allow for more complex backgrounds than just single colors.
type Scene struct {
  objects []*ExtendedObject
  background color.SphericalColorFunction
}

func NewScene(objects []*ExtendedObject, background color.SphericalColorFunction) *Scene {
  if objects == nil || background == nil { return nil }

  return &Scene{objects, background}
}

//Traces a light ray through a scene. 
func (scene *Scene) TracePath(pos, dir []float64, max_depth int, receptor_tolerance float64) []float64 {
  var last int = - 1

  ray := &LightRay{0, pos, dir, []float64{4, 5, 6}, []float64{1, 1, 1}, []float64{0, 0, 0}, 1}

  var u float64
  var s Interactor
  var selected int

  //Follow the ray for max_depth bounces. 
  //TODO make each bounce a separate function call.
  for ray.depth = 0; ray.depth < max_depth; ray.depth ++ {
    u = math.Inf(1)
    selected = -1

    //check every shape for intersection. 
    for l, object := range scene.objects {
      if l != last { //Except not the last one, since the ray is right on the surface.
        intersection := object.surf.Intersection(ray.position, ray.direction)

        //An object can return several intersection parameters, so we have to check each one.
        for m := 0; m < len(intersection); m ++ {
          if intersection[m] < u && intersection[m] > 0 {
            u = intersection[m]
            selected = l
          }
        }
      }
    }

    if selected == -1 { //The ray has diverged to infinity.
      last = -1
      bg := scene.background(ray.direction)(ray.receptor)
      for i := 0; i < 3; i ++ {
        ray.color[i] *= bg[i]
      }
      break
    }

    //The ray has interacted with something.
    ray.Trace(u)
    s = scene.objects[selected].interactor(ray.position)
    last = selected

    //Interact with the object that the ray intersected first.
    ray = s.Interact(ray)

    //check if we should bother continuing to bounce the ray.
    if ray.redirected <= receptor_tolerance {break}
  }

  return ray.DeriveColor()
}


