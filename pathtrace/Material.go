package pathtrace

import "math/rand"
import "math"
import "../distributions"
import "../vector"

//Functions can be mocked out for testing purposes. 
var randomUnitSphereSurfacePoint func() *[3]float64 = distributions.RandomUnitSphereSurfacePoint

var randomNormallyDistributedVector func(int, float64, float64) []float64 =
  distributions.RandomNormallyDistributedVector

/*type Color interface {
}*/

type RayInteraction interface {
  Interact(direction, normal []float64) []float64
}

type lambertianReflection struct {
}

//The Lambertian reflectance algorithm used here works as follows.
//  1. generate a vector v uniformly distributed on the surface of a unit sphere.
//  2. Add this vector to the normal vector. 
//This should generate the correct distribution of vectors. 
func (m *lambertianReflection) Interact(direction, normal []float64) []float64 {

  reflect := make([]float64, len(normal))
  if len(normal) == 3 {
    var dir *[3]float64
    dir = randomUnitSphereSurfacePoint()
    for i := 0; i < len(dir); i ++ {
      reflect[i] = normal[i] + (*dir)[i]
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
        reflect[i] = normal[i] + dir[i] / r2
      }
    }
  }

  return reflect
}

func NewLambertianReflection() RayInteraction {
  return &lambertianReflection{}
}

type mirrorReflection struct {
}

func (m *mirrorReflection) Interact(direction, normal []float64) []float64 {
  reflect := make([]float64, len(normal))
  //Find the dot product of the normal with the incoming ray.
  d := vector.Dot(normal, direction)
  //Mirror the ray in the direction of the normal. 
  for l := 0; l < len(normal); l ++ {
    reflect[l] = direction[l] - 2 * normal[l] * d
  }

  return reflect
}

func NewMirrorReflection() RayInteraction {
  return &mirrorReflection{}
}

//TODO
type specularReflection struct {
  scatter float64
}

func (m *specularReflection) Interact(direction, normal []float64) []float64 {
  var d float64	
  reflect := make([]float64, len(normal))
  //Find the dot product of the normal with the incoming ray.
  for l := 0; l < len(normal); l ++ {
    d += normal[l] * direction[l]
  }

  //Mirror the ray in the direction of the normal. 
  for l := 0; l < len(normal); l ++ {
    reflect[l] = direction[l] - 2 * normal[l] * d
    d += reflect[l] * reflect[l]
  }

  //Normalize the outgoing ray. 
  d = math.Sqrt(d)
  for l := 0; l < len(normal); l ++ {
    reflect[l] /= d
  }

  //Add a random jostling. 
  spec := randomNormallyDistributedVector(len(normal), 0, m.scatter)
  for l := 0; l < len(normal); l ++ {
    reflect[l] += spec[l]
  }

  //Test find the dot product of the new vector with the normal.
  d = 0
  for l := 0; l < len(normal); l ++ {
    d += normal[l] * reflect[l]
  }

  //If the light ray has gone into the surface, mirror it with
  //the normal vector again. 
  if d < 0 {
    for l := 0; l < len(normal); l ++ {
      reflect[l] = reflect[l] - 2 * normal[l] * d
    }
  }

  return reflect
}

func NewSpecularReflection(scatter float64) RayInteraction {
  return &specularReflection{scatter}
}

//This refraction does not take into account the way that refraction
//changes with color. 
type basicRefractiveTransmission struct {
  index float64
}

func (m *basicRefractiveTransmission) Interact(direction, normal []float64) []float64 {
  var d float64	
  //Find the dot product of the normal with the incoming ray.
  for l := 0; l < len(normal); l ++ {
    d += normal[l] * direction[l]
  }

  d = (m.index - 1) * d / m.index
  //Mirror the ray in the direction of the normal. 
  reflect := make([]float64, len(normal))
  for l := 0; l < len(normal); l ++ {
    reflect[l] = direction[l] - normal[l] * d
  }

  return reflect
}

func NewBasicRefractiveTransmission(index float64) RayInteraction {
  return &specularReflection{index}
}

//TODO
type orenNayerReflection struct {
}

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

