package main

import (
  "image"
  "image/color"
  "image/png"
  "math"
  "math/rand"
  "os"
  "fmt"
  "./surface"
  "./pathtrace"
  "./distributions"
)

//Some prototype objects which will be programmed for real later.
type ray struct {
  pos, dir, receptor []float64
  //As the ray bounces along, if it encounters glowing objects,
  //these numbers keep track of how the final color should be
  //adjusted to take these earlier interactions into account.
  adj_num []float64
  adj_den float64 
}

type object struct {
  surf surface.Surface
  redirect func(r *ray, normal []float64) *ray
}

func (obj *object) Interact(r *ray) *ray {
  return obj.redirect(r, surface.SurfaceNormal(obj.surf, r.pos))
}

type scene struct {
  obj []object
  background []float64
}

func (sc *scene) TracePath(pos, dir []float64, max_depth int, receptor_tolerance float64) []float64 {
  var last int = - 1

  r := &ray{pos, dir, []float64{1, 1, 1}, []float64{0, 0, 0}, 1}

  //Follow the ray for max_depth bounces. 
  //TODO make each bounce a separate function call.
  for k := 0; k < max_depth; k ++ {
    var u float64 = math.Inf(1)
    var s *object = nil
    var selected int

    //check every shape for intersection. 
    for l := 0; l < len(sc.obj); l ++ {
      if l != last {
        intersection := sc.obj[l].surf.Intersection(r.pos, r.dir)
        for m := 0; m < len(intersection); m ++ {
          if intersection[m] < u && intersection[m] > 0 {
            u = intersection[m]
            s = &sc.obj[l]
            selected = l
          }
        }
      }
    }

    if s == nil { //The ray has diverged to infinity.
      for i := 0; i < 3; i ++ {
        r.receptor[i] = sc.background[i]
      }
      break
    }

    //The ray has interacted with something.
    last = selected
    //Generate the new ray position and surface normal for the interaction.
    for l := 0; l < 3; l ++ {
      r.pos[l] = r.pos[l] + u * r.dir[l]
    }

    r = s.Interact(r)

    //check if we should bother continuing to bounce the ray.
    if r.adj_den <= receptor_tolerance {break}
  }

  for i := 0; i < 3; i ++ {
    r.receptor[i] = r.receptor[i] * r.adj_den + r.adj_num[i]
  }

  return r.receptor
}

type traceFunc func(r *ray, norm []float64) *ray

//This function makes an object glow. 
func Glow(c []float64) traceFunc {
  return func(r *ray, norm []float64) *ray {
      r.adj_den = 0
      r.adj_num = c
      return r
  }
}

//This makes an object glow AND do something else. 
func GlowAverage(c []float64, absorb, transmit float64,
  ref pathtrace.RayInteraction, surf surface.Surface) traceFunc {
  return func(r *ray, norm []float64) *ray {
    for i := 0; i < 3; i ++ {
      r.adj_num[i] += r.adj_den * absorb * c[i]
    }
    r.adj_den *= transmit
    r.dir = ref.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    return r
  }
}

var lambertian pathtrace.RayInteraction = pathtrace.NewLambertianReflection()
var mirror pathtrace.RayInteraction = pathtrace.NewMirrorReflection()

func Lambertian(c []float64, surf surface.Surface) traceFunc {
  return func(r *ray, norm []float64) *ray {
    for l := 0; l < 3; l ++ {
      r.receptor[l] *= c[l]
    }
    r.dir = lambertian.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    return r
  }
}

//TODO These next two really should just send out two rays.
func Shiney(c []float64, shiney float64, p float64, surf surface.Surface) traceFunc {
  reflect := pathtrace.NewSpecularReflection(shiney)

  return func(r *ray, norm []float64) *ray {
    if rand.Float64() > p {
      r.dir = reflect.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    } else {
      r.dir = lambertian.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    }
    return r
  }
}

func Transparent(refraction, shiney, p, q float64, surf surface.Surface) traceFunc {
  refract := pathtrace.NewBasicRefractiveTransmission(refraction)
  reflect := pathtrace.NewSpecularReflection(shiney)
  pq := p + q
  a := p / pq

  return func(r *ray, norm []float64) *ray {
    for i := 0; i < 3; i ++ {
      r.receptor[i] *= pq
    }
    if rand.Float64() > a {
      r.dir = refract.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    } else {
      r.dir = reflect.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
    }
    return r
  }
}

//Snap a photo! 
func Snap(sc *scene, cam_func pathtrace.RayFunc, size_u, size_v,
  depth, minp, maxp int, maxMeanVariance float64) *image.NRGBA {
  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))
  var ray_pos, ray_dir []float64

  pix_sum := make([]float64, 3)
  pix := make([]float64, 3)

  var n int = 0
  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      //TODO put some heuristic thing here to report on how complete the image is.

      for k := 0; k < 3; k ++ {
        pix_sum[k] = 0
      }

      var p int = 0
      var variance_check bool

      //Set up the variance monitor.
      var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics()}

      for {
        //Set up the ray.
        ray_pos, ray_dir = cam_func(i, j)

        p ++

        //Trace the path.
        c := sc.TracePath(ray_pos, ray_dir, depth, 1./256.)

        for l := 0; l < 3; l ++ {
          pix_sum[l] += c[l]
          monitor[l].AddVariable(c[l])
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

  return img
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

  //Four spheres in a tetrahedron. 
  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}
  colors := [][]float64{[]float64{1, 1, 0}, []float64{1, 0, 1}, []float64{0, 1, 1}, []float64{1, 1, 1}}
  background := []float64{0,0,0}

  //Create the objects and scene. 
  objects := make([]object, len(spheres))
  for i := 0; i < len(spheres); i ++ {
    objects[i].surf = spheres[i]
    c := colors[i]
    objects[i].redirect = Glow(c)
  }
  scene_1 := &scene{objects, background}

  //get camera function
  cam_pos   := []float64{0,0,3}
  cam_look  := []float64{0,0,0}
  cam_up    := []float64{0,1,0}
  cam_right := []float64{1,0,0}
  cam_func := pathtrace.FlatCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 1.33333, 1.)

  img := Snap(scene_1, cam_func, size_u, size_v, 1, 1, 1, 1)

  png.Encode(file, img)
}

//in this demo, the spheres reflect light and produce a fractal.
func pathtrace_activity_02(smooth_reflection bool) {
  createOutputDirectory()
  fmt.Println("Running path trace activity 02")
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_02.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  var size_u, size_v int = 1600, 1200

  //Set up the scene. Four spheres again, each with a different color.
  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}
  colors := [][]float64{[]float64{1, 1, 0}, []float64{1, 0, 1},
    []float64{0, 1, 1}, []float64{1, 1, 1}}
  background := []float64{0,0,0}

  //These spheres are all perfectly reflective. 
  reflection := pathtrace.NewMirrorReflection()

  //Set up the scene object. 
  objects := make([]object, len(spheres))
  for i := 0; i < len(spheres); i ++ {
    surf := spheres[i]
    objects[i].surf = surf
    c := colors[i]
    var absorb float64 = .5

    if smooth_reflection {
      //This is the much faster and better way of calculating the color. 
      objects[i].redirect = GlowAverage(c, absorb, 1 - absorb, reflection, surf) 
    } else {
      //This mimmics my original way of calculating colors, just to show for comparison.
      objects[i].redirect = func(r *ray, norm []float64) * ray {
        if rand.Float64() > absorb {
          r.dir = reflection.Interact(r.dir, surface.SurfaceNormal(surf, r.pos))
          return r
        } else {
          for i := 0; i < 3; i ++ {
            r.adj_num[i] = c[i]
            r.receptor[i] = 0
          }
          r.adj_den = 0
          return r
        }
      }
    }
  }
  scene_2 := &scene{objects, background}

  //Set up the camera. 
  cam_pos   := []float64{0, 0, 2.6}
  cam_look  := []float64{0, 0, 0}
  cam_up    := []float64{0, 1, 0}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.FlatCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 1.33333/2., 1./2.)

  //Four hundred bounces, 16 rays per pixel. 
  var depth, minp, maxp int = 400, 16, 1000
  //Using the new awy of calculating pixels, there should be almost no variance with each ray.
  var maxMeanVariance float64 = .00001

  img := Snap(scene_2, cam_func, size_u, size_v, depth, minp, maxp, maxMeanVariance)

  png.Encode(file, img)
}

//A prototype which will eventually show off a variety of materials.
func pathtrace_activity_03() {
  createOutputDirectory()
  fmt.Println("Running path trace activity 03")
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_03.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  var size_u, size_v int = 800, 600

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
  green := []float64{.2, .8, .3}
  orange := []float64{.6, 1, .1}
  white := []float64{1, 1, 1}
  light := []float64{1.6, 1.6, 1.6}

  background := []float64{0, 0, 0}

  objects := make([]object, len(shapes))
  for i := 0; i < len(shapes); i ++ {
    objects[i].surf = shapes[i]
  }

  //This object is just a light. 
  objects[0].redirect = Glow(light)
  //Glows and also reflects. 
  objects[1].redirect = GlowAverage(pink, .3, 1, mirror, shapes[1])
  objects[2].redirect = Shiney(blue, .2, .1, shapes[2])
  objects[3].redirect = Lambertian(green, shapes[3])
  objects[4].redirect = Transparent(1.6, .1, .5, 1, shapes[4])
  objects[5].redirect = Lambertian(orange, shapes[5])
  objects[6].redirect = Lambertian(white, shapes[6])
  objects[7].redirect = Lambertian(pink, shapes[7])
  objects[8].redirect = Lambertian(white, shapes[8])

  scene_3 := &scene{objects, background}

  cam_pos   := []float64{0, -3, 4}
  cam_look  := []float64{0, 0, 0}
  cam_up    := []float64{0, 1., 0}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.FlatCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 1.33333, 1)

  var depth, minp, maxp int = 40, 100, 5000
  var maxMeanVariance float64 = .0001

  img := Snap(scene_3, cam_func, size_u, size_v, depth, minp, maxp, maxMeanVariance)

  png.Encode(file, img)
}

//A room in which different objects can be set. Something is wrong with this scene. 
func pathtrace_activity_04() {
  createOutputDirectory()
  fmt.Println("Running path trace activity 04")
  //Check if the file can be written. 
  file, err := os.Create("./output/activity_04.png")
  if err != nil {
    fmt.Println("Could not write file: ", err.Error())
    return 
  }

  //Aspect ratio is (4/3)^3
  var size_u, size_v int = 1536, 648
  var total_pixels = size_u * size_v

  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))

  var outer_dim float64   = 30
  var room_width float64  = 18
  var room_height float64 = 12

  var box func([]float64, [][]float64) surface.Surface = surface.NewParallelpipedByCornerAndEdges
  var sub func(surface.Surface, surface.Surface) surface.Surface = surface.NewSubtraction
  var add func(surface.Surface, surface.Surface) surface.Surface = surface.NewAddition

  room :=
    sub(box([]float64{
        (room_width - outer_dim)/2, (room_width - outer_dim)/2, (room_height - outer_dim)/2},
       [][]float64{[]float64{outer_dim, 0, 0}, []float64{0, outer_dim, 0}, []float64{0, 0, outer_dim}}),
      add(add(box([]float64{0, 0, 0},
          [][]float64{
            []float64{room_width, 0, 0},
            []float64{0, room_width, 0},
            []float64{0, 0, room_height - 2}}),
        box([]float64{0, 0, room_height - 1.75}, 
            [][]float64{[]float64{room_width, 0, 0}, []float64{0, room_width, 0}, []float64{0, 0, 1.75}})),
          box([]float64{1, 1, room_height - 3}, 
              [][]float64{
                []float64{room_width - 2, 0, 0}, 
                []float64{0, room_width - 2, 0}, 
                []float64{0, 0, room_height - 1}})))

  light := sub(
      box([]float64{-1, -1, room_height - 1.7},
        [][]float64{
          []float64{room_width + 2, 0, 0}, 
          []float64{0, room_width + 2, 0},
          []float64{0, 0, .1}}),
      box([]float64{.9, .9, room_height - 3},
        [][]float64{
          []float64{room_width - 1.8, 0, 0},
          []float64{0, room_width - 1.8, 0},
          []float64{0, 0, 2}}))

  shapes := []surface.Surface{light, room}

  lambertian := pathtrace.NewLambertianReflection()

  background := []float64{0, 0, 0}
  var light_intensity float64 = 4

  cam_pos   := []float64{2, 2, 6}
  cam_look  := []float64{room_width, room_width, room_height}
  cam_up    := []float64{0, 0, 1}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.CylindricalCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 2.37, 1)

  var ray_pos, ray_dir, pix, pix_sum []float64 =
    make([]float64, 3), make([]float64, 3), make([]float64, 3), make([]float64, 3)

  var depth, minp, maxp int = 40, 100, 5000
  var maxMeanVariance float64 = .001

  var n int = 0
  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {
      if n % 2700 == 0 {
        fmt.Println("  ", float64(n)/float64(total_pixels), " complete.")
      }

      for k := 0; k < 3; k ++ {
        pix_sum[k] = 0
      }

      var p int = 0
      var variance_check bool

      //Set up the variance monitor.
      var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics()}

      for {
        //Set up the ray.
        //Set up the ray.
        ray_pos, ray_dir = cam_func(i, j)

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
                pix[l] = receptor[l] * light_intensity
                pix_sum[l] += pix[l]
                monitor[l].AddVariable(pix[l])
              }
              break
            } else if selected == 1 {  
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
