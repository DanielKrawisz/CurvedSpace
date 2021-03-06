package pathtrace

import "math/rand"
import "math"
import "github.com/DanielKrawisz/CurvedSpace/distributions"
import "github.com/DanielKrawisz/CurvedSpace/vector"

//TODO turn these back into objects so that they can trace multiple rays.

//Functions can be mocked out for testing purposes. 
var randomUnitSphereSurfacePoint func() *[3]float64 = distributions.RandomUnitSphereSurfacePoint

var randomNormallyDistributedVector func(int, float64, float64) []float64 =
  distributions.RandomNormallyDistributedVector

//An function that redirects a ray
type Redirection func([]float64, []float64) []float64 

//The Lambertian reflectance algorithm used here works as follows.
//  1. generate a vector v uniformly distributed on the surface of a unit sphere.
//  2. Add this vector to the normal vector. 
//This should generate the correct distribution of vectors. 
func LambertianReflection(direction, normal []float64) []float64 {
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

  return reflect
}

func MirrorReflection(direction, normal []float64) []float64 {
  reflect := make([]float64, len(normal))
  //Find the dot product of the normal with the incoming ray.
  d := vector.Dot(normal, direction)
  //Mirror the ray in the direction of the normal. 
  for l := 0; l < len(normal); l ++ {
    reflect[l] = direction[l] - 2 * normal[l] * d
  }

  return reflect
}

//TODO Try this with my other idea for doing specular reflection.
func SpecularReflection(scatter float64) Redirection {
  return func (direction, normal []float64) []float64 {
    reflect := make([]float64, len(normal))
    //Find the dot product of the normal with the incoming ray.
    d := vector.Dot(normal, direction)

    //Mirror the ray in the direction of the normal. 
    for l := 0; l < len(normal); l ++ {
      reflect[l] = direction[l] - 2 * normal[l] * d
    }

    //Normalize the outgoing ray. 
    vector.Normalize(reflect)

    //Add a random jostling. 
    spec := randomNormallyDistributedVector(len(normal), 0, scatter)
    for l := 0; l < len(normal); l ++ {
      reflect[l] += spec[l]
    }

    //Test find the dot product of the new vector with the normal.
    d = vector.Dot(normal, reflect)

    //If the light ray has gone into the surface, mirror it with
    //the normal vector again. 
    if d < 0 {
      for l := 0; l < len(normal); l ++ {
        reflect[l] = reflect[l] - 2 * normal[l] * d
      }
    }

    return reflect
  }
}

//This refraction does not take into account the way that refraction
//changes with color. 
func BasicRefraction(index float64) Redirection {
  inv := 1/index
  return func(direction, normal []float64) []float64 {
    //Find the dot product of the normal with the incoming ray.
    vector.Normalize(direction)
    c := -vector.Dot(normal, direction)

    //Special cases for whether we are going in or coming out of the object.
    var r, sign float64
    if c > 0 {
      r = inv
      sign = -1
    } else {
      r = index
      sign = 1
    }

    //This determines whether the ray is reflected or transmitted. 
    rad := 1 - r * r * (1 - c * c)

    if rad < 0 {
      return MirrorReflection(direction, normal)
    } else {
      return vector.LinearSum(r, r * c + sign * math.Sqrt(rad), direction, normal)
    }
  }
}

//Scatters a ray in a random direction.
func ScatterRedirector(degree float64) Redirection {
  //Independent of the normal vector given it--this could even be nil! 
  return func (direction, normal []float64) []float64 {

    //Normalize the outgoing ray. 
    vector.Normalize(direction)

    //Add a random jostling. 
    scatter := randomNormallyDistributedVector(len(normal), 0, degree)
    for l := 0; l < len(normal); l ++ {
      direction[l] += scatter[l]
    }

    return direction
  }
}

//TODO Something like oren-nayer
