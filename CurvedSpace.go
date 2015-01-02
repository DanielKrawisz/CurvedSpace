package main

//This will be a program to do ray-tracing over curved spaces. 
//It doesn't do curved spaces yet though. Right now there is
//a numerical differential equation solver, quadratic surfaces
//and some demos.

//Short-term goals. 
//TODO: complete tests for basic surfaces. 
//TODO: complete intersection tests for booleans and higher polynomials.
//TODO: create materials.
//TODO: allow for a refractive index that varies over space and color. 
//TODO: create polyhedra.
//TODO: create spacially bounded objects for faster intersection computations.

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
  "./pathtrace"
  "./distributions"
  //"./BlackHoles"
)

func main() {
  pathtrace_activity_01()
  pathtrace_activity_02()
  pathtrace_activity_03()
}

func createOutputDirectory() {
  src, err := os.Stat("output")
  if err != nil {
    if os.IsNotExist(err) {
      os.Mkdir("output", 0777)
      return 
    } else {
      panic(err)
    }
  }

  if !src.IsDir() {
    fmt.Println("Source is not a directory")
    os.Exit(1)
  }
}

//The purpose of the following demos is not only to show what the
//program can do, but to prototype future features to get an idea
//of how to design them. Thus, some of them show off things that
//this cannot do in general yet. 

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
  createOutputDirectory()
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_01.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  var size_u, size_v int = 640, 480

  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))

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

  png.Encode(file, img)
}

//in this demo, the spheres reflect light and produce a fractal.
func pathtrace_activity_02() {
  createOutputDirectory()
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_02.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  var size_u, size_v int = 800, 600
  var total_pixels = size_u * size_v

  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))

  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}

  reflectance := []pathtrace.RayInteraction{
    pathtrace.NewMirrorReflection(), 
    pathtrace.NewMirrorReflection(),
    pathtrace.NewMirrorReflection(), 
    pathtrace.NewMirrorReflection()}

  colors := [][]float64{[]float64{1, 1, 0}, []float64{1, 0, 1},
    []float64{0, 1, 1}, []float64{1, 1, 1}}
  background := []float64{0,0,0}

  cam_pos := []float64{0, 0, 3}
  cam_dir := []float64{0, 0, -1}
  cam_up  := []float64{0, 1./2., 0}
  cam_right := []float64{1.333/2., 0, 0}

  var ray_pos, ray_dir, col []float64 = make([]float64, 3), make([]float64, 3), make([]float64, 3)

  var depth, minp int = 40, 50
  var maxMeanVariance float64 = .00001
  var absorbed float64 = .5

  var n int = 0
  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      if n % 24000 == 0 {
        fmt.Println("  ", float64(n)/float64(total_pixels), " complete.")
      }

      //Set up the color for this pixel.
      for k := 0; k < 3; k ++ {
        col[k] = 0
      }

      var p int = 0
      var variance_check bool

      for {
        //Set up the variance monitor.
        var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
          distributions.NewSampleStatistics(), 
          distributions.NewSampleStatistics(), 
          distributions.NewSampleStatistics()}

        //Set up the ray.
        var ou float64 = 2*(float64(i - size_u/2) + rand.Float64() - .5) / float64(size_u)
        var ov float64 = 2*(float64(j - size_v/2) + rand.Float64() - .5) / float64(size_v)
        for k := 0; k < 3; k ++ {
          ray_pos[k] = cam_pos[k]
          ray_dir[k] = cam_dir[k] + ov * cam_up[k] + ou * cam_right[k]
        }

        p ++

        var last int = -1
        var selected int

        //follow the ray for k bounces. 
        for k := 0; k < depth; k ++ {
          var u float64 = math.Inf(1)
          var s surface.Surface = nil
          var c []float64
          var r pathtrace.RayInteraction

          //Test for intersection with every object in the scene. 
          for l := 0; l < len(spheres); l ++ {
            if l != last {
              intersection := spheres[l].Intersection(ray_pos, ray_dir)
              for m := 0; m < len(intersection); m ++ {
                if intersection[m] < u && intersection[m] > 0 {
                  u = intersection[m]
                  s = spheres[l]
                  c = colors[l]
                  r = reflectance[l]
                  selected = l
                }
              }
            }
          }

          if s == nil { //The ray has diverged to infinity. 
            for l := 0; l < 3; l ++ {
              col[l] += background[l];
            }
            break
          } else { //The ray has intersected an object. 
            if rand.Float64() > absorbed {
              //TODO The picture would be a lot smoother if the ray were
              //averaged with the glow color rather than randomly assigned
              //to one or the other. 
              for l := 0; l < 3; l ++ {
                ray_pos[l] = ray_pos[l] + u * ray_dir[l]
              }
              last = selected
              ray_dir = r.Interact(ray_dir, surface.SurfaceNormal(s, ray_pos))
            } else {
              for l := 0; l < 3; l ++ {
                col[l] += c[l]
              }
              break
            }
          }
        }

        //Check the variance of the pixel to see if it is low enough.
        if p > minp {
          variance_check = true
          for l := 0; l < 3; l ++ {
            if monitor[l].MeanVariance() > maxMeanVariance {
              variance_check = false
            }
          }
        }
        if variance_check {
          break
        }
      }

      //Generate the pixel. 
      for l := 0; l < 3; l ++ {
        col[l] = math.Min(255 * col[l] / float64(p), 255)
      }

      img.Set(i, j, &color.NRGBA{uint8(col[0]), uint8(col[1]), uint8(col[2]), 255})

      n++
    }
  }

  png.Encode(file, img)
}

//A prototype which will eventually show off a variety of materials.
//TODO enable multithreaded computation. Path-tracing can be paralellized! 
//TODO It is very easy to make images that are over- or under-exposed. 
//Allow the program to handle this gracefully. 
//TODO image is very grainy. Need to find out why. There has to be something
//more efficient than just increasing the number of rays! 
func pathtrace_activity_03() {
  createOutputDirectory()
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_03.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  var size_u, size_v int = 800, 600
  var total_pixels = size_u * size_v

  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))

  //TODO some indirect lighting would be nice. 
  shapes := []surface.Surface{
    surface.NewSphere([]float64{0, 0, 26}, 14),
    surface.NewSphere([]float64{0, 0, 1}, 1),
    surface.NewSphere([]float64{2, 0, 1}, 1),
    surface.NewSphere([]float64{-2, 0, 1}, 1),
    surface.NewSphere([]float64{1, 1.73205, 1}, 1),
    surface.NewSphere([]float64{1, -1.73205, 1}, 1),
    surface.NewSphere([]float64{-1, 1.73205, 1}, 1),
    surface.NewSphere([]float64{-1, -1.73205, 1}, 1),
    surface.NewPlaneByPointAndNormal([]float64{0, 0, 0}, []float64{0, 0, 1})}

  pink := []float64{1, .4, 1}
  blue := []float64{.3, .7, 1}

  lambertian := pathtrace.NewLambertianReflection()
  mirror := pathtrace.NewMirrorReflection()
  specular := pathtrace.NewSpecularReflection(.1)
  refract := pathtrace.NewBasicRefractiveTransmission(1.6)

  background := []float64{0,0,0}

  cam_pos := []float64{0, 0, 5}
  cam_dir := []float64{0, 0, -1}
  cam_up  := []float64{0, 1., 0}
  cam_right := []float64{1.333,0,0}

  var ray_pos, ray_dir, pix, pix_sum []float64 =
    make([]float64, 3), make([]float64, 3), make([]float64, 3), make([]float64, 3)

  var depth, minp, maxp int = 40, 100, 5000
  var maxMeanVariance float64 = .0001

  var n int = 0
  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      if n % 480 == 0 {
        fmt.Println("  ", float64(n)/float64(total_pixels), " complete.")
      }

      for k := 0; k < 3; k ++ {
        pix_sum[k] = 0
      }

      var p int = 0
      var variance_check bool

      for {
        //Set up the variance monitor.
        var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
          distributions.NewSampleStatistics(), 
          distributions.NewSampleStatistics(), 
          distributions.NewSampleStatistics()}

        //Set up the ray.
        var ou float64 = 2*(float64(i - size_u/2) + rand.Float64() - .5)/float64(size_u)
        var ov float64 = 2*(float64(j - size_v/2) + rand.Float64() - .5)/float64(size_v)
        for k := 0; k < 3; k ++ {
          ray_pos[k] = cam_pos[k]
          ray_dir[k] = cam_dir[k] + ov * cam_up[k] + ou * cam_right[k]
        }
        receptor := []float64{1, 1, 1} 

        p ++
        var last int = -1

        //Follow the ray for k bounces. 
        for k := 0; k < depth; k ++ {
          var u float64 = math.Inf(1)
          var s surface.Surface = nil
          var selected int

          //check every shape for intersection. 
          for l := 0; l < len(shapes); l ++ {
            if l != last {
            intersection := shapes[l].Intersection(ray_pos, ray_dir)
              for m := 0; m < len(intersection); m ++ {
                if intersection[m] < u && intersection[m] > 0 {
                  u = intersection[m]
                  s = shapes[l]
                  selected = l
                }
              }
            }
          }

          if s == nil { //The ray has diverged to infinity. 
            for l := 0; l < 3; l ++ {
              pix[l] = receptor[l] * background[l]
              pix_sum[l] += pix[l]
              monitor[l].AddVariable(pix[l])
            }
            break
          } else { //The ray has interacted with something.
            last = selected
            //Generate the new ray position and surface normal for the interaction.
            for l := 0; l < 3; l ++ {
              ray_pos[l] = ray_pos[l] + u * ray_dir[l]
            }
            normal := surface.SurfaceNormal(s, ray_pos)

            if selected == 0 { //This object is just a light. 
              for l := 0; l < 3; l ++ {
                pix[l] = receptor[l] * 2
                pix_sum[l] += pix[l]
                monitor[l].AddVariable(pix[l])
              }
              break
            } else if selected == 1 { //Glows and also reflects. 
              if rand.Float64() > .3 {
                for l := 0; l < 3; l ++ {
                  pix[l] = receptor[l] * 1 * pink[l]
                  pix_sum[l] += pix[l]
                  monitor[l].AddVariable(pix[l])
                }
                break
              } else {
                for l := 0; l < 3; l ++ {
                  receptor[l] *= 2
                }
                ray_dir = mirror.Interact(ray_dir, normal)
              }
            } else if selected == 2 { //Has a specular and a diffuse component.
              if rand.Float64() > .5 {
                for l := 0; l < 3; l ++ {
                  pix[l] = receptor[l] * 1.8 * blue[l]
                  pix_sum[l] += pix[l]
                  monitor[l].AddVariable(pix[l])
                }
                ray_dir = lambertian.Interact(ray_dir, normal)
              } else {
                for l := 0; l < 3; l ++ {
                  receptor[l] *= 2
                }
                ray_dir = specular.Interact(ray_dir, normal)
              }
            } else if selected == 3 { //Only has lambertian reflection. 
              lambertian.Interact(ray_dir, normal)
            } else if selected == 4 { //Partly transparent. 
              for l := 0; l < 3; l ++ {
                //TODO need to handle this sort of thing in its own object.
                //AN object that transmits AND reflects light needs to either
                //Send a ray in both directions or double the intensity. 
                receptor[l] *= 2 
              }
              if rand.Float64() > .5 {
                ray_dir = mirror.Interact(ray_dir, normal)
              } else {
                ray_dir = refract.Interact(ray_dir, normal)
              }
            } else if selected == 5 {
              ray_dir = lambertian.Interact(ray_dir, normal)
            } else if selected == 6 {
              ray_dir = lambertian.Interact(ray_dir, normal)
            } else if selected == 7 {
              ray_dir = lambertian.Interact(ray_dir, normal)
            } else if selected == 8 {
              ray_dir = lambertian.Interact(ray_dir, normal)
            } 
          }
        }

        if p > maxp {
          break
        }
        //Check variance
        if p > minp {
          variance_check = true
          for l := 0; l < 3; l ++ {
            if monitor[l].MeanVariance() > maxMeanVariance {
              variance_check = false
            }
          }
          if variance_check {
            break
          }
        }
      }

      //Generate the pixel. 
      for l := 0; l < 3; l ++ {
        pix[l] = math.Min(255 * pix_sum[l] / float64(p), 255)
      }

      img.Set(i, j, &color.NRGBA{uint8(pix[0]), uint8(pix[1]), uint8(pix[2]), 255})

      n++
    }
  }

  png.Encode(file, img)
}
