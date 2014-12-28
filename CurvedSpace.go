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
  "math/rand"
  //"bufio"
  "./diffeq"
  "./geometry"
  "./surface"
  //"./BlackHoles"
)

func main() {
  pathtrace_activity_02()
}

//A simple demo of the differential equation solver 
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

//A simple demo of the most basic form of path-tracing. There are four spheres, 
//each with a different color, and they only emit light, but do not reflect it.
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

  var ray_pos, ray_dir []float64 = make([]float64, 3), make([]float64, 3)

  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      //fmt.Println("pixel ", i, ", ", j)

      var ou float64 = 2*float64(i - size_u/2)/float64(size_u)
      var ov float64 = 2*float64(j - size_v/2)/float64(size_v)

      //Set up the ray.
      for k := 0; k < 3; k ++ {
        ray_pos[k] = cam_pos[k]
        ray_dir[k] = cam_dir[k] + ov * cam_up[k] + ou * cam_right[k]
      }

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
      }

      if s == nil {
        c = background
      }

      img.Set(i, j, c)
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

//in this demo, the spheres also reflect light. 
func pathtrace_activity_02() {
  var size_u, size_v int = 640, 480

  img := image.NewNRGBA(image.Rect(0, 0, size_u - 1, size_v - 1))

  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}

  colors := [][]float64{[]float64{1, 1, 0}, []float64{1, 0, 1},
    []float64{0, 1, 1}, []float64{1, 1, 1}}
  background := []float64{0,0,0}

  cam_pos := []float64{0,0,3}
  cam_dir := []float64{0,0,-1}
  cam_up  := []float64{0,1./2.,0}
  cam_right := []float64{1.333/2.,0,0}

  var ray_pos, ray_dir, col []float64 = make([]float64, 3), make([]float64, 3), make([]float64, 3)

  var depth, rpix int = 20, 400
  var absorbed float64 = .3

  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      //fmt.Println("pixel ", i, ", ", j)

      for k := 0; k < 3; k ++ {
        col[k] = 0
      }

      for p := 0; p < rpix; p ++ {
        //Set up the ray.
        var ou float64 = 2*(float64(i - size_u/2) + rand.Float64() - .5)/float64(size_u)
        var ov float64 = 2*(float64(j - size_v/2) + rand.Float64() - .5)/float64(size_v)
        for k := 0; k < 3; k ++ {
          ray_pos[k] = cam_pos[k]
          ray_dir[k] = cam_dir[k] + ov * cam_up[k] + ou * cam_right[k]
        }

        for k := 0; k < depth; k ++ {
          var u float64 = math.Inf(1)
          var s surface.Surface = nil
          var c []float64

          for l := 0; l < len(spheres); l ++ {
            intersection := spheres[l].Intersection(ray_pos, ray_dir)
            for m := 0; m < len(intersection); m ++ {
              if intersection[m] < u && intersection[m] > 0 {
                u = intersection[m]
                s = spheres[l]
                c = colors[l]
              }
            }
          }

          if s == nil {
            for l := 0; l < 3; l ++ {
              col[l] += background[l];
            }
            break
          } else {
            if rand.Float64() > absorbed {
              var d float64 = 0
              //fmt.Println("ray: ", ray_pos, " + u ", ray_dir)
              for l := 0; l < 3; l ++ {
                ray_pos[l] = ray_pos[l] + u * ray_dir[l]
              }
              n := surface.SurfaceNormal(s, ray_pos)
              //fmt.Println("hit sphere ", s.String(), "; u = ", u, "; normal = ", n)
              for l := 0; l < 3; l ++ {
                d += n[l] * ray_dir[l]
              }
              for l := 0; l < 3; l ++ {
                ray_dir[l] = ray_dir[l] - 2 * n[l] * d
              }
              //fmt.Println("new ray: ", ray_pos, " + u ", ray_dir)
              //bio := bufio.NewReader(os.Stdin)
              //_, _, _ = bio.ReadLine()
            } else {
              for l := 0; l < 3; l ++ {
                col[l] += c[l]
              }
              break
            }
          }
        }
      }

      for l := 0; l < 3; l ++ {
        col[l] = math.Min(255 * col[l] / float64(rpix), 255)
      }

      img.Set(i, j, &color.NRGBA{uint8(col[0]), uint8(col[1]), uint8(col[2]), 255})
    }
  }

  //Save the image.
  file, err := os.Create("./output/activity_02.png")
  if err == nil {
    png.Encode(file, img)
  } else {
    fmt.Println("Could not write file")
  }
}
