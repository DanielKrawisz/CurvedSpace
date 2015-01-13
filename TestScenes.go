package main

import "os"
import "fmt"
import "image/png"
import "./surface"
import "./pathtrace"

func test_scene_01() {
  createOutputDirectory()

  ground := surface.NewPlaneByPointAndNormal([]float64{0, 0, 0}, []float64{0, 0, 1})
  small_sphere := surface.NewSphere([]float64{0, 0, 2}, 2)
  small_box    := surface.NewParallelpipedByCornerAndEdges([]float64{-2, -2, 0},
                    [][]float64{[]float64{4, 0, 0}, []float64{0, 4, 0}, []float64{0, 0, 4}})
  big_sphere   := surface.NewSphere([]float64{-5, 0, 3}, 3)
  big_box      := surface.NewParallelpipedByCornerAndEdges([]float64{-8, -3, 0},
                    [][]float64{[]float64{6, 0, 0}, []float64{0, 6, 0}, []float64{0, 0, 6}})

  background := []float64{0, 0, 0}
  light := []float64{4, 4, 4}
  white := []float64{.8, .8, .8}
  blue := []float64{0, .5, .8}
  //pink := []float64{1, .8, .6}

  scenes := []*scene{
    &scene{[]*object{
      &object{ground, Lambertian(white, ground)}, 
      &object{small_sphere, Glow(light)},
      &object{big_sphere, Lambertian(blue, big_sphere)}}, background}, 
    &scene{[]*object{
      &object{ground, Lambertian(white, ground)}, 
      &object{small_sphere, Glow(light)},
      &object{big_box, Lambertian(blue, big_box)}}, background}, 
    &scene{[]*object{
      &object{ground, Lambertian(white, ground)}, 
      &object{small_box, Glow(light)},
      &object{big_sphere, Lambertian(blue, big_sphere)}}, background}, 
    &scene{[]*object{
      &object{ground, Lambertian(white, ground)}, 
      &object{small_box, Glow(light)},
      &object{big_box, Lambertian(blue, big_box)}}, background}}

  var pix_u, pix_v int = 640, 480

  cam_pos   := []float64{0, -5, 3.75}
  cam_look  := []float64{-1, 0, 2}
  cam_up    := []float64{0, 0, 1}
  cam_right := []float64{1, 0, 0}
  cam_func := pathtrace.CylindricalCamera(cam_pos,
    pathtrace.CameraMatrix(cam_pos, cam_look, cam_up, cam_right), pix_u, pix_v, 1.333, 1)

  var depth, minp, maxp int = 40, 16, 100
  var maxMeanVariance float64 = .01

  for i, sc := range scenes {
    fmt.Println("Running path trace test scene 01 -", i)
    //Check if the file can be written. 
    file, err := os.Create(fmt.Sprint("./output/test_scene_01-", i,".png"))
    if err != nil {
      fmt.Println("Could not write file: ", err.Error())
      continue 
    }

    img := Snap(sc, cam_func, pix_u, pix_v, depth, minp, maxp, maxMeanVariance, .05, 1000000)

    png.Encode(file, img)
  }
}
