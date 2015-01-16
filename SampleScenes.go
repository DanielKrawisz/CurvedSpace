package main

import (
  "image/png"
  "./surface"
  "./pathtrace"
  "./vector"
)

//The purpose of the following demos is not only to show what the
//program can do, but to prototype future features to get an idea
//of how to design them. Thus, some of them show off things that
//this cannot do in general yet. 

//A simple demo of the most basic form of path-tracing. There are four spheres, 
//each with a different color, and they only emit light, but do not reflect it.
func pathtrace_activity_01() {

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
  objects := make([]*pathtrace.ExtendedObject, len(spheres))
  for i := 0; i < len(spheres); i ++ {
    objects[i] = pathtrace.NewExtendedObject(spheres[i], pathtrace.NewGlowingObject(colors[i]))
  }
  scene_1 := pathtrace.NewScene(objects, background)

  //get camera function
  cam_pos   := []float64{0,0,3}
  cam_look  := []float64{0,0,0}
  cam_up    := []float64{0,1,0}
  cam_right := []float64{1,0,0}
  cam_func := pathtrace.FlatCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 1.33333, 1.)

  img := pathtrace.Snapshot(scene_1, cam_func, size_u, size_v, 1, 1, 1, 1, 1, 100000)

  file := getHandleToOutputFile("activity 01", "activity_01.png")
  if file == nil {return}

  png.Encode(file, img)
}

//in this demo, the spheres reflect light and produce a fractal.
func pathtrace_activity_02() {

  var size_u, size_v int = 1600, 1200

  white := []float64{1, 1, 1}
  yellow := []float64{1, 1, 0}
  magenta := []float64{1, 0, 1}
  cyan := []float64{0, 1, 1}

  //Set up the scene. Four spheres again, each with a different color.
  spheres := []surface.Surface{
    surface.NewSphere([]float64{1, 0, 0}, .866025),
    surface.NewSphere([]float64{-1./2., 0.866025, 0}, .866025),
    surface.NewSphere([]float64{-1./2., -0.866025, 0}, .866025),
    surface.NewSphere([]float64{0, 0, -0.8556}, .866025)}
  colors := [][]float64{yellow, magenta, cyan, white}
  background := []float64{0,0,0}

  //Set up the scene object. 
  objects := make([]*pathtrace.ExtendedObject, len(spheres))
  for i := 0; i < len(spheres); i ++ {
    objects[i] = pathtrace.NewExtendedObject(spheres[i], 
      pathtrace.NewMirrorReflector(spheres[i], pathtrace.GlowAbsorbAverage(colors[i], white, .5)))
  }
  scene_2 := pathtrace.NewScene(objects, background)

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

  img := pathtrace.Snapshot(scene_2, cam_func, size_u, size_v,
    depth, minp, maxp, maxMeanVariance, .01, 1000000)

  file := getHandleToOutputFile("activity 02", "activity_02.png")
  if file == nil {return}

  png.Encode(file, img)
}

//A prototype which will eventually show off a variety of materials.
func pathtrace_activity_03() {

  var size_u, size_v int = 1600, 1200

  glow_pink := []float64{1.5, .6, 1.5}
  light_pink := []float64{.94, .75, .86}
  blue := []float64{.3, .7, 1}
  green := []float64{.2, .8, .3}
  orange := []float64{1, .6, .1}
  white := []float64{1, 1, 1}
  light := []float64{2, 2, 2}

  objects := []surface.Surface{
    surface.NewSphere([]float64{0, 0, 26}, 14),
    surface.NewSphere([]float64{0, 0, 1}, 1),
    surface.NewSphere([]float64{2, 0, 1}, 1),
    surface.NewBounding(
      surface.NewSphere([]float64{-2, 0, 1}, 1), 
      surface.NewInsubstantialSurface(3, .001)), 
    surface.NewSphere([]float64{1, 1.73205, 1}, 1),
    surface.NewSphere([]float64{1, -1.73205, 1}, 1),
    surface.NewSphere([]float64{-1, 1.73205, 1}, 1),
    surface.NewSphere([]float64{-1, -1.73205, 1}, 1),
    surface.NewPlaneByPointAndNormal([]float64{0, 0, 0}, []float64{0, 0, 1})}

  scene_3 := pathtrace.NewScene([]*pathtrace.ExtendedObject{
    pathtrace.NewExtendedObject(objects[0],
      pathtrace.NewGlowingObject(light)), 
    pathtrace.NewExtendedObject(objects[1],
      pathtrace.NewMirrorReflector(objects[1], pathtrace.GlowAbsorbAverage(glow_pink, white, .5))), 
    pathtrace.NewExtendedObject(objects[2], 
      pathtrace.NewShineyInteractor(objects[2], pathtrace.Absorb(blue), .1, .2)), 
    pathtrace.NewExtendedObject(objects[3],
      pathtrace.NewScatterTransmitter(pathtrace.GlowAbsorbAverage(white, light_pink, .2), 1.4)), 
    pathtrace.NewExtendedObject(objects[4], 
      pathtrace.NewGlassInteractor(objects[4], pathtrace.Absorb(white), 1.6, .4, .8)), 
    pathtrace.NewExtendedObject(objects[5], 
      pathtrace.NewShineyInteractor(objects[5], pathtrace.Absorb(orange), .2, .3)), 
    pathtrace.NewExtendedObject(objects[6], 
      pathtrace.NewShineyInteractor(objects[6], pathtrace.Absorb(green), .5, .4)), 
    pathtrace.NewExtendedObject(objects[7], 
      pathtrace.NewLambertianReflector(objects[7], pathtrace.Absorb(green))), 
    pathtrace.NewExtendedObject(objects[8],
      pathtrace.NewShineyInteractor(objects[8], pathtrace.Absorb(white), .35, .25))}, []float64{0, 0, 0})

  cam_pos   := []float64{0, 3, 4}
  cam_look  := []float64{0, 0, 0}
  cam_up    := []float64{0, 0, 1}
  cam_right := []float64{-1, 0, 0}
  cam_func := pathtrace.FlatCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, 1.33333 * .85, .85)

  var depth, minp, maxp int = 40, 100, 5000
  var maxMeanVariance float64 = .002

  img := pathtrace.Snapshot(scene_3, cam_func, size_u, size_v,
    depth, minp, maxp, maxMeanVariance, .01, 1000000)

  file := getHandleToOutputFile("activity 03", "activity_03.png")
  if file == nil {return}

  png.Encode(file, img)
}

//A room in which different objects can be set. Something is wrong with this scene. 
func pathtrace_activity_04() {

  //Aspect ratio is (4/3)^3
  var size_u, size_v int = 1536, 648 // 1536, 648 // 768, 324 // 384, 162

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

  torus := surface.NewTorus([]float64{9, 9, 6}, vector.Normalize([]float64{0.32, 1, 0.}), 4.5, 1.5)

  background  := []float64{0, 0, 0}
  white_light := []float64{4, 4, 4}
  white       := []float64{1, 1, 1}
  blue        := []float64{.22, .56, .87}

  //shapes := []surface.Surface{light, room}
  scene_4 := pathtrace.NewScene(
    []*pathtrace.ExtendedObject{
      pathtrace.NewExtendedObject(light,
        pathtrace.NewGlowingObject(white_light)),
      pathtrace.NewExtendedObject(room,
        pathtrace.NewLambertianReflector(room, pathtrace.Absorb(white))), 
      pathtrace.NewExtendedObject(torus, 
        pathtrace.NewShineyInteractor(torus, pathtrace.Absorb(blue), .1, .2))},
    background)

  cam_pos   := []float64{3, 3, 5}
  cam_look  := []float64{room_width, room_width, room_height-2.5}
  cam_up    := []float64{0, 0, 1}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.CylindricalCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), size_u, size_v, .7 * 2.37, .7 * 1)

  var depth, minp, maxp int = 10, 16, 5000
  var maxMeanVariance float64 = .002

  img := pathtrace.Snapshot(scene_4, cam_func, size_u, size_v,
    depth, minp, maxp, maxMeanVariance, .01, 1000000)

  file := getHandleToOutputFile("activity 04", "activity_04.png")
  if file == nil {return}	

  png.Encode(file, img)
}
