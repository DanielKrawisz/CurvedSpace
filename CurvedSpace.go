package main

//This will be a program to do ray-tracing over curved spaces. 

//TODO: complete tests for basic surfaces. 
//TODO: complete intersection tests for booleans and higher polynomials.
//TODO: complete bounce monitor. 
//TODO: create new default system to run with bounces. 
//TODO: create solver that can be used for path tracing.
//TODO: create luminescent materials.
//TODO: create transparent materials.
//TODO: create reflective materials.
//TODO: create refractive materials.
//TODO: allow for a refractive index that varies over time. 
//TODO: create polyhedra.

import (
  "fmt"
  "image"
  "image/color"
  "image/png"
  "math"
  "os"
  "./diffeq"
  "./geometry"
  "./surface"
  //"./BlackHoles"
)

func main() {
  pathtrace_activity_01()
}

func diffeq_activity_01() {
  var v geometry.CoordinatePoint = nil

  fmt.Println(v)

  var dimension int = 2
  var particles int = 3
  np := diffeq.NewSpringSystem(dimension, particles, 
    [][]float64{
      []float64{0, 0},
      []float64{0, 1},
      []float64{0, 2}},
    [][]float64{
      []float64{0, 0},
      []float64{1, 0},
      []float64{-1, 0}}, 
    []float64{1, 1, 1},
    [][][]float64{nil, 
      [][]float64{
        []float64{1}}, 
      [][]float64{
        []float64{0}, []float64{1}}}, .1, 10,
    diffeq.NewRungeKuttaSolverMethodDormandPrince(2 * dimension * particles, .000001))

  if np != nil {np.Run()}
}

func pathtrace_activity_01() {
  var size_u, size_v int = 640, 480

  img := image.NewNRGBA(image.Rect(0, 0, size_u - 1, size_v - 1))

  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}
  colors := []color.Color{&color.NRGBA{255, 255, 0, 255}, &color.NRGBA{255, 0, 255, 255},
    &color.NRGBA{0, 255, 255, 255}, &color.NRGBA{255, 255, 255, 255}}
  background := &color.NRGBA{0,0,0,255}

  cam_pos := []float64{0,0,3}
  cam_dir := []float64{0,0,-1}
  cam_up  := []float64{0,1,0}
  cam_right := []float64{1.333,0,0}

  var depth int = 3

  var ray_pos, ray_dir []float64 = make([]float64, 3), make([]float64, 3)

  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      //fmt.Println("pixel ", i, ", ", j)

      var ou float64 = 2*float64(i - size_u/2)/float64(size_u)
      var ov float64 = 2*float64(j - size_u/2)/float64(size_u)

      //Set up the ray.
      for k := 0; k < 3; k ++ {
        ray_pos[k] = cam_pos[k]
        ray_dir[k] = cam_dir[k] + ou * cam_up[k] + ov * cam_right[k]
      }

      for k := 0; k < depth; k ++ {
        var u float64 = math.Inf(1)
        var s surface.Surface = nil
        var c color.Color

        for l := 0; l < len(spheres); l ++ {
          intersection := spheres[l].Intersection(ray_pos, ray_dir)
          for m := 0; m < len(intersection); m ++ {
            if intersection[m] < u && intersection[m] > 0 {
              u = intersection[m]
              s = spheres[l]
              c = colors[l]
            }
          }

          if s == nil {
            c = background
            break
          }
        }

        img.Set(i, j, c)
      }
    }
  }

  //Save the image.
  file, err := os.Create("./output/activity_01.png")
  if err == nil {
    png.Encode(file, img)
  } else {
    fmt.Println("Could not write file")
  }
}
