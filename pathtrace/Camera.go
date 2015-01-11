package pathtrace

import "math"
import "math/rand"
import "../vector"

func CameraMatrix(pos, look, up, right []float64) [][]float64 {
  return vector.Orthonormalize([][]float64{vector.Minus(look, pos), up, right})
}

func camCoordinates(i, j, pix_u, pix_v int, fov_u, fov_v float64) (float64, float64){
  return 2 * fov_u * (float64(i - pix_u/2) + rand.Float64() - .5) / float64(pix_u),
    2 * fov_v * (float64(j - pix_v/2) + rand.Float64() - .5) / float64(pix_v)
}

type RayFunc func(int, int) ([]float64, []float64)

func FlatCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) RayFunc {
  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = mtrx[0][k] + ov * mtrx[1][k] + ou * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func CylindricalCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) RayFunc {
  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      ray_dir[k] = math.Cos(ou) * mtrx[0][k] + ov * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func PolarSphericalCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) RayFunc {
  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      c := math.Cos(ov)
      ray_dir[k] = math.Cos(ou) * c * mtrx[0][k] +
        math.Sin(ou) * c * mtrx[1][k] + math.Sin(ou) * mtrx[2][k]
    }
    return ray_pos, ray_dir
  }
}

func SphericalCamera(pos []float64, mtrx [][]float64, pix_u, pix_v int, fov_u, fov_v float64) RayFunc {
  return func(i, j int) ([]float64, []float64) {
    ray_pos, ray_dir := make([]float64, len(pos)), make([]float64, len(pos))
    var ou, ov float64 = camCoordinates(i, j, pix_u, pix_v, fov_u, fov_v)
    for k := 0; k < 3; k ++ {
      ray_pos[k] = pos[k]
      cv  := math.Cos(ov)
      cu  := math.Cos(ou)
      sv  := math.Sin(ov)
      su  := math.Sin(ov)
      cv2 := cv * cv
      cu2 := cu * cu
      sv2 := sv * sv
      su2 := su * su
      d := math.Sqrt(cu2 * cv2 + cv2 * su2* + cu2 * sv2)
      ray_dir[k] = (cu * cv * mtrx[0][k] + cu * sv * mtrx[1][k] + cv * su * mtrx[2][k]) / d
    }
    return ray_pos, ray_dir
  }
}
