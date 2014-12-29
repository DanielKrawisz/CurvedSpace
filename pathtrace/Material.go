package pathtrace

import "math/rand"
import "math"
import "../distributions"

type RayInteraction interface {
  Interact(direction, normal []float64)
}

type noInteraction struct {
  c Color
}

func (m *noInteraction) Interact(direction, normal []float64) {
}

type lambertianReflectance struct {
}

//The Lambertian reflectance algorithm used here works as follows.
//  1. generate a vector v uniformly distributed on the surface of a unit sphere.
//  2. Add this vector to the normal vector. 
//This should generate the correct distribution of vectors. 
func (m *lambertianReflectance) Interact(direction, normal []float64) {

  if len(normal) == 3 {
    var dir *[3]float64
    dir = distributions.RandomUnitSphereSurfacePoint()
    for i := 0; i < len(dir); i ++ {
      direction[i] = normal[i] + (*dir)[i]
    }
  } else {
    var dir []float64
    var r2 float64
    dir = make([]float64, len(normal))

    for {
      for {
        r2 = 0
        for i := 0; i < len(dir); i ++ {
          dir[i] = 2 * rand.Float64() - 1
          r2 += dir[i] * dir[i]
        }

        if r2 >= 0 {
          break
        }
      }

      r2 = math.Sqrt(r2)
      for i := 0; i < len(dir); i ++ {
        direction[i] = normal[i] + dir[i] / r2
      }
    }
  }
}

func NewLambertianReflectance() RayInteraction {
  return &lambertianReflectance{}
}

type mirrorReflectance struct {
}

func (m *mirrorReflectance) Interact(direction, normal []float64) {
  var d float64	
  for l := 0; l < 3; l ++ {
    d += normal[l] * direction[l]
  }
  for l := 0; l < 3; l ++ {
    direction[l] = direction[l] - 2 * normal[l] * d
  }
}

func NewMirrorReflectance() RayInteraction {
  return &mirrorReflectance{}
}

/*type refract struct {
  c Color
  n float64
  scatter float64
}

type phongSpecularReflectance struct {
  scatter float64
}

func (m *phongReflectance) Interact(l LightRay, normal float64) {
}

type orenNayerReflectance struct {
  c Color
}*/

/*type compoundMaterial struct {
  probability []float64
  material []Material
}

func (m *compoundMaterial) Interact(l LightRay, normal float64) {
  p := rand.Float64()
  for i := 0; i < len(m.probability); i ++ {
    if p < m.probability[i] {
      m.material[i].Interact(l, normal)
      return
    }
  }
}

func NewCompoundMaterial(probability []float64, material []Material) Material {
  if probability == nil || material == nil {
    return nil
  }
  if len(probability) != len(material) {
    return nil
  }
  if len(probability) < 1 {return nil}
  if probability[len(probability) - 1] != 1.0 {return nil}

  var p float64 = 0.0
  for i := 0; i < len(probability); i++ {
    if material[i] == nil {return nil}
    if probability[i] <= p {
      return nil
    } else {
      p = probability[i]
    }
  }

  return &compoundMaterial{probability, material}
}*/

