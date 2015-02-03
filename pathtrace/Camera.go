package pathtrace

import "math"
import "math/rand"
import "github.com/DanielKrawisz/CurvedSpace/vector"

//TODO allow the cameras to have a focal point. 

func CameraMatrix(pos, look, up, right []float64) [][]float64 {
  return vector.Orthonormalize([][]float64{vector.Minus(look, pos), up, right})
}

func CameraStochastic() float64 {
  return rand.Float64() - .5
}

var camJitter func() float64 = CameraStochastic

func CameraCoordinates(i, j, pix_u, pix_v int, fov_u, fov_v float64) (float64, float64){
  return 2 * fov_u * (float64(i) - float64(pix_u - 1)/2. + camJitter()) / float64(pix_u - 1),
    -2 * fov_v * (float64(j) - float64(pix_v - 1)/2. + camJitter()) / float64(pix_v - 1)
}

type GenerateRay func(int, int) ([]float64, []float64)

//The camera rays are given by evenly-spaced points on a grid on a plane. 
func IsometricCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k] + ov * mtrx[1][k] + ou * mtrx[2][k]
      ray_dir[k] = mtrx[0][k] + ov * mtrx[1][k] + ou * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

//The camera rays are given by evenly-spaced points on a grid on a plane. 
func FlatCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = mtrx[0][k] + ov * mtrx[1][k] + ou * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

//Inverse of the flat camera; the eye is like a plane aimed at a single point. 
//(Simulates a pinhole camera)
func InverseFlatCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -mtrx[0][k] - ov * mtrx[1][k] - ou * mtrx[2][k]
      ray_pos[k] = pos[k] - ray_dir[k]
    }
    return ray_pos, ray_dir
  }
}

//The camera rays are given by points at equal angles around a cylinder
//and equal distances up and down it. 
func CylindricalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = math.Cos(ou) * mtrx[0][k] + ov * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

//Inverse of the cylindrical camera. 
func InverseCylindricalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -math.Cos(ou) * mtrx[0][k] - ov * mtrx[1][k] - math.Sin(ou) * mtrx[2][k]
      ray_pos[k] = pos[k] - ray_dir[k]
    }
    return ray_pos, ray_dir
  }
}

//The camera rays are given by points at equal angles around a cylinder
//and equal distances up and down it. 
func IsometricCylindricalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k] + ov * mtrx[1][k]
      ray_dir[k] = math.Cos(ou) * mtrx[0][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func InverseIsometricCylindricalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = - math.Cos(ou) * mtrx[0][k] - math.Sin(ou) * mtrx[2][k]
      ray_pos[k] = -ray_dir[k] + pos[k] + ov * mtrx[1][k]
    }
    return ray_pos, ray_dir
  }
}

//The camera rays uses points on a sphere that are on equally
//spaced lines of lattitude and longetude. 
func PolarSphericalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    c := math.Cos(ov)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = math.Cos(ou) * c * mtrx[0][k] +
        math.Sin(ov) * mtrx[1][k] + math.Sin(ou) * c * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func InversePolarSphericalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    c := math.Cos(ov)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -(math.Cos(ou) * c * mtrx[0][k] +
        math.Sin(ov) * mtrx[1][k] + math.Sin(ou) * c * mtrx[2][k])
      ray_pos[k] = -ray_dir[k] + pos[k]
    }
    return ray_pos, ray_dir
  }
}

//This camera is a little bit more difficult to make sense of. Basically
//Imagine a sphere and take two planes going through the center of the
//sphere. The planes are rotated respectively around two fixed axes
//perpendicular to one another. The intersections of the two planes gives
//the direction of the camera ray. 
func SphericalCamera(pos []float64, mtrx [][]float64, 
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    cv  := math.Cos(ov)
    cu  := math.Cos(ou)
    sv  := math.Sin(ov)
    su  := math.Sin(ou)
    gz  := math.Sqrt(cu * cu * cv * cv + su * su)

    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = (cu * cv * mtrx[0][k] + su * sv * mtrx[1][k] + cv * su * mtrx[2][k]) / gz
    }
    return ray_pos, ray_dir
  }
}

//Inverse of the spherical camera. This can be used to simulate an eye-shaped
//pinhole camera. 
func InverseSphericalCamera(pos []float64, mtrx [][]float64, 
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = CameraCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    cv  := math.Cos(ov)
    cu  := math.Cos(ou)
    sv  := math.Sin(ov)
    su  := math.Sin(ou)
    gz  := math.Sqrt(cu * cu * cv * cv + su * su)

    for k := 0; k < 3; k ++ {
      ray_dir[k] = -(cu * cv * mtrx[0][k] + su * sv * mtrx[1][k] + cv * su * mtrx[2][k]) / gz
      ray_pos[k] = pos[k] - ray_dir[k]
    }
    return ray_pos, ray_dir
  }
}

//I don't think anybody's ever used a toroidial camera before. 
//The input format for this toroidial camera is kind of weird,
//but it's not so bad once you get the hang of it! 
//
// pos = center point of the torus.
//   R = two vectors giving the major axes of the torus. 
//   r = two more vectors giving the minor axes
//
//(only three dimensional)
func ToroidialCamera(pos []float64, R [][]float64, r [][]float64, pix_u, pix_v int) GenerateRay {

  if pos == nil || R == nil || r == nil { return nil }
  if len(R) != 2 || len(r) != 2 { return nil }

  for i := 0; i < 2; i ++ {
    if r[i] == nil || R[i] == nil { return nil }
    if len(R[i]) != 3 || len(r[i]) != 3 {return nil}
  }

  quad := vector.Dot(R[0], R[0]) * vector.Dot(R[1], R[1])
  norm := math.Sqrt(quad)
  R2 := vector.Cross([][]float64{R[0], R[1]})

  R0r0 := vector.Dot(R[0], r[0]) / norm
  R1r0 := vector.Dot(R[1], r[0]) / norm
  R2r0 := vector.Dot(R2, r[0]) / quad

  return func(i, j int) ([]float64, []float64) {
    var ou, ov float64 = CameraCoordinates(i, j, pix_u + 1, pix_v + 1, math.Pi, math.Pi)

    cu := math.Cos(ou)
    su := math.Sin(ou)
    cv := math.Cos(ov)
    sv := math.Sin(ov)

    return vector.Plus(pos, vector.LinearSum(cu, su, R[0], R[1])), 
      vector.LinearSum(sv, cv, r[1],
        vector.LinearSum(R2r0, 1, R2,
          vector.LinearSum(cu, su, 
           vector.LinearSum(R0r0, R1r0, R[0], R[1]), 
           vector.LinearSum(R0r0, -R1r0, R[1], R[0]))))
  }
}

func InverseToroidialCamera(pos []float64, R [][]float64, r [][]float64, pix_u, pix_v int) GenerateRay {

  if pos == nil || R == nil || r == nil { return nil }
  if len(R) != 2 || len(r) != 2 { return nil }

  for i := 0; i < 2; i ++ {
    if r[i] == nil || R[i] == nil { return nil }
    if len(R[i]) != 3 || len(r[i]) != 3 {return nil}
  }

  quad := vector.Dot(R[0], R[0]) * vector.Dot(R[1], R[1])
  norm := math.Sqrt(quad)
  R2 := vector.Cross([][]float64{R[0], R[1]})

  R0r0 := vector.Dot(R[0], r[0]) / norm
  R1r0 := vector.Dot(R[1], r[0]) / norm
  R2r0 := vector.Dot(R2, r[0]) / quad

  return func(i, j int) ([]float64, []float64) {
    var ou, ov float64 = CameraCoordinates(i, j, pix_u + 1, pix_v + 1, math.Pi, math.Pi)

    cu := math.Cos(ou)
    su := math.Sin(ou)
    cv := math.Cos(ov)
    sv := math.Sin(ov)

    ray_dir := 
      vector.LinearSum(-sv, -cv, r[1],
        vector.LinearSum(R2r0, 1, R2,
          vector.LinearSum(cu, su, 
           vector.LinearSum(R0r0, R1r0, R[0], R[1]), 
           vector.LinearSum(R0r0, -R1r0, R[1], R[0]))))

    return vector.Minus(vector.Plus(pos, vector.LinearSum(cu, su, R[0], R[1])), ray_dir), ray_dir
  }
}

//Another crazy concept. (requires polygonal surface first)
/*func TrapezoidalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      c := math.Cos(ov)
      ray_dir[k] = math.Cos(ou) * c * mtrx[0][k] -
        math.Sin(ou) * c * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}*/
