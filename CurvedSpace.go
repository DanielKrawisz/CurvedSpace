package main

//This will be a program to do ray-tracing over curved spaces. 
//It doesn't do curved spaces yet though. Right now there are
//a numerical differential equation solver, quadratic surfaces,
//solid constructive geometry, and some demos. There are also
//some ways that light rays can interact with a surface. 

//Short-term goals. 
//TODO: update the symmetric tensor contraction functions to use the symmetric permutation loop. 
//TODO: functions to change the coordinates of polynomial objects and to translate polynomial objects.
//TODO: complete tests for basic surfaces. 
//TODO:   complete intersection tests for higher polynomials.
//TODO: make torus surface.
//TODO: work on materials. Materials should have glow, transmissive, and reflective components.
//      Ultimately, each should be able to send out its own light rays. 
//TODO:   need to make scene object first.
//TODO: allow for a refractive index that varies over space and color. 
//TODO: create polyhedra.
//TODO: add a glow mode which randomly assigns some pastel color to each object so as to generate
//      quick test images of a scene. 

//Longer-term goals.
//TODO enable multithreaded computation. Path-tracing can be paralellized! 
//TODO It is very easy to make images that are over- or under-exposed. 
//Allow the program to handle this gracefully. 
//TODO image is very grainy. Make less grainy. 
//TODO allow for solid objects that affect the light day during its entire
//course through it.
//TODO other kinds of boundary conditions: elliptic and hyperbolic geometry!
//TODO curved space with any kind of metric we want. 
//TODO wormholes.  

import (
  "fmt"
  "os"
  //"bufio"
  "./diffeq"
  "./geometry"
  //"./BlackHoles"
)

func main() {
  pathtrace_activity_01()
  pathtrace_activity_02(true)
  pathtrace_activity_03()
  //pathtrace_activity_04()
  //pathtrace_test_01(false)
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
