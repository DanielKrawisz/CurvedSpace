package main

import "fmt"
import "image/png"
import "github.com/DanielKrawisz/CurvedSpace/surface"
import "github.com/DanielKrawisz/CurvedSpace/pathtrace"
import "github.com/DanielKrawisz/CurvedSpace/color"

func test_scene_01() {

  scene := func(i int) *pathtrace.Scene {
    ground := surface.NewPlaneByPointAndNormal([]float64{0, 0, 0}, []float64{0, 0, 1})
    small_sphere := surface.NewSphere([]float64{0, 0, 2}, 2)
    small_box    := surface.NewParallelpipedByCornerAndEdges([]float64{-2, -2, 0},
                      [][]float64{[]float64{4, 0, 0}, []float64{0, 4, 0}, []float64{0, 0, 4}})
    big_sphere   := surface.NewSphere([]float64{-5, 0, 3}, 3)
    big_box      := surface.NewParallelpipedByCornerAndEdges([]float64{-8, -3, 0},
                      [][]float64{[]float64{6, 0, 0}, []float64{0, 6, 0}, []float64{0, 0, 6}})

    background := color.ConstantColorFunction(color.PresetColor([]float64{0,0,0}))
    light := []float64{4, 4, 4}
    white := []float64{.8, .8, .8}
    blue := []float64{0, .5, .8}
    //pink := []float64{1, .8, .6}

    switch i {
    case 0 : 
      return pathtrace.NewScene([]*pathtrace.ExtendedObject{
        pathtrace.NewExtendedObject(ground,
          pathtrace.NewLambertianReflector(ground, pathtrace.Absorb(white))), 
        pathtrace.NewExtendedObject(small_sphere, pathtrace.NewGlowingObject(light)),
        pathtrace.NewExtendedObject(big_sphere,
          pathtrace.NewLambertianReflector(big_sphere, pathtrace.Absorb(blue)))}, background); 
    case 1 : 
      return pathtrace.NewScene([]*pathtrace.ExtendedObject{
        pathtrace.NewExtendedObject(ground,
          pathtrace.NewLambertianReflector(ground, pathtrace.Absorb(white))), 
        pathtrace.NewExtendedObject(small_sphere, pathtrace.NewGlowingObject(light)),
        pathtrace.NewExtendedObject(big_box,
          pathtrace.NewLambertianReflector(big_box, pathtrace.Absorb(blue)))}, background);  
    case 2 : 
      return pathtrace.NewScene([]*pathtrace.ExtendedObject{
        pathtrace.NewExtendedObject(ground,
          pathtrace.NewLambertianReflector(ground, pathtrace.Absorb(white))), 
        pathtrace.NewExtendedObject(small_box, pathtrace.NewGlowingObject(light)),
        pathtrace.NewExtendedObject(big_sphere,
          pathtrace.NewLambertianReflector(big_sphere, pathtrace.Absorb(blue)))}, background);
    case 3 : 
      return pathtrace.NewScene([]*pathtrace.ExtendedObject{
        pathtrace.NewExtendedObject(ground,
          pathtrace.NewLambertianReflector(ground, pathtrace.Absorb(white))), 
        pathtrace.NewExtendedObject(small_box, pathtrace.NewGlowingObject(light)),
        pathtrace.NewExtendedObject(big_box,
          pathtrace.NewLambertianReflector(big_box, pathtrace.Absorb(blue)))}, background); 
    }

    return nil
  }

  var pix_u, pix_v int = 640, 480

  cam_pos   := []float64{0, -5, 3.75}
  cam_look  := []float64{-1, 0, 2}
  cam_up    := []float64{0, 0, 1}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.CylindricalCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), pix_u, pix_v, 1.333, 1)

  var depth, minp, maxp int = 40, 16, 100
  var maxMeanVariance float64 = .01

  for i := 0; i < 4; i++ {
    i := i

    file := getHandleToOutputFile(fmt.Sprint("test scene 01-", i), fmt.Sprint("test_scene_01-", i,".png"))
    if file == nil {continue}

    img := pathtrace.Snapshot(func() *pathtrace.Scene {return scene(i)},
      cam_func, pix_u, pix_v, depth, minp, maxp, maxMeanVariance, .05, 1000000, 8)

    png.Encode(file, img)
  }
}
