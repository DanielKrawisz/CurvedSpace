package pathtrace

import "math"
import "math/rand"
import "../vector"

//TODO allow the cameras to have a focal point. 

func CameraMatrix(pos, look, up, right []float64) [][]float64 {
  return vector.Orthonormalize([][]float64{vector.Minus(look, pos), up, right})
}

func camCoordinates(i, j, pix_u, pix_v int, fov_u, fov_v float64) (float64, float64){
  return 2 * fov_u * (float64(i - pix_u/2) + rand.Float64() - .5) / float64(pix_u),
    2 * fov_v * (float64(j - pix_v/2) + rand.Float64() - .5) / float64(pix_v)
}

type GenerateRay func(int, int) ([]float64, []float64)

//The camera rays are given by evenly-spaced points on a grid on a plane. 
func IsometricCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k] - ov * mtrx[1][k] + ou * mtrx[2][k]
      ray_dir[k] = mtrx[0][k] - ov * mtrx[1][k] + ou * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

//The camera rays are given by evenly-spaced points on a grid on a plane. 
func FlatCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = mtrx[0][k] - ov * mtrx[1][k] + ou * mtrx[2][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -mtrx[0][k] + ov * mtrx[1][k] - ou * mtrx[2][k]
      ray_pos[k] = pos[k] - ray_pos[k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = math.Cos(ou) * mtrx[0][k] - ov * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -math.Cos(ou) * mtrx[0][k] + ov * mtrx[1][k] - math.Sin(ou) * mtrx[2][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k] - ov * mtrx[1][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = - math.Cos(ou) * mtrx[0][k] - math.Sin(ou) * mtrx[2][k]
      ray_pos[k] = -ray_dir[k] + pos[k] - ov * mtrx[1][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    c := math.Cos(ov)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = math.Cos(ou) * c * mtrx[0][k] -
        math.Sin(ou) * c * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func InversePolarSphericalCamera(pos []float64, mtrx [][]float64,
  pix_u, pix_v int, fov_u, fov_v float64) GenerateRay {
  if pos == nil || mtrx == nil { return nil }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    c := math.Cos(ov)
    for k := 0; k < 3; k ++ {
      ray_dir[k] = -math.Cos(ou) * c * mtrx[0][k] +
        math.Sin(ou) * c * mtrx[1][k] - math.Sin(ou) * mtrx[2][k]
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    cv  := math.Cos(ov)
    cu  := math.Cos(ou)
    sv  := math.Sin(ov)
    su  := math.Sin(ov)
    cv2 := cv * cv
    cu2 := cu * cu
    sv2 := sv * sv
    su2 := su * su
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      d := math.Sqrt(cu2 * cv2 + cv2 * su2* + cu2 * sv2)
      ray_dir[k] = (cu * cv * mtrx[0][k] - cu * sv * mtrx[1][k] + cv * su * mtrx[2][k]) / d
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
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    cv  := math.Cos(ov)
    cu  := math.Cos(ou)
    sv  := math.Sin(ov)
    su  := math.Sin(ov)
    cv2 := cv * cv
    cu2 := cu * cu
    sv2 := sv * sv
    su2 := su * su
    for k := 0; k < 3; k ++ {
      d := math.Sqrt(cu2 * cv2 + cv2 * su2* + cu2 * sv2)
      ray_dir[k] = -(cu * cv * mtrx[0][k] - cu * sv * mtrx[1][k] + cv * su * mtrx[2][k]) / d
      ray_pos[k] = -ray_dir[k] + pos[k]
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
//   r = two sets of two vectors giving the minor axes corresponding to each major one.
//
//(only three dimensional)
func ToroidialCamera(pos []float64, R [][]float64, r [][][]float64, pix_u, pix_v int) GenerateRay {

  if pos == nil || R == nil || r == nil { return nil }
  if len(R) != 2 || len(r) != 4 { return nil }

  for i := 0; i < 2; i ++ {
    if r[i] == nil || R[i] == nil { return nil }
    if len(R[i]) != 3 {return nil}
    for j := 0; j < 2 ; j ++ {
      if r[i][j] == nil { return nil }
      if len(r[i][j]) != 3 {return nil}
    }
  }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u + 1, pix_v + 1, math.Pi, math.Pi)

    Cou := math.Cos(ou)
    Sou := math.Sin(ou)
    Cov := math.Cos(ov)
    Sov := math.Sin(ov)

    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k] + Cou * R[0][k] + Sou * R[1][k]
      ray_dir[k] = Cou * (Cov * r[0][0][k] - Sov * r[0][1][k]) + Sou * (Cov * r[1][0][k] - Sov * r[1][1][k])
    }
    return ray_pos, ray_dir
  }
}

func InverseToroidialCamera(pos []float64, R [][]float64, r [][][]float64, pix_u, pix_v int) GenerateRay {

  if pos == nil || R == nil || r == nil { return nil }
  if len(R) != 2 || len(r) != 4 { return nil }

  for i := 0; i < 2; i ++ {
    if r[i] == nil || R[i] == nil { return nil }
    if len(R[i]) != 3 {return nil}
    for j := 0; j < 2 ; j ++ {
      if r[i][j] == nil { return nil }
      if len(r[i][j]) != 3 {return nil}
    }
  }

  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u + 1, pix_v + 1, math.Pi, math.Pi)

    Cou := math.Cos(ou)
    Sou := math.Sin(ou)
    Cov := math.Cos(ov)
    Sov := math.Sin(ov)

    for k := 0; k < 3; k ++ {
      ray_dir[k] = -Cou * (Cov * r[0][0][k] - Sov * r[0][1][k]) -
        Sou * (Cov * r[1][0][k] - Sov * r[1][1][k])
      ray_pos[k] = pos[k] + Cou * R[0][k] + Sou * R[1][k] - ray_dir[k]
    }
    return ray_pos, ray_dir
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
