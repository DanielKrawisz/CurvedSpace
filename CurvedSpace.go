package main

//This will be a program to do ray-tracing over curved spaces. 

//TODO: complete tests for basic surfaces. 
//TODO: complete grad tests for polynomials and booleans. 
//TODO: complete intersection tests for booleans and polynomials.
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
  "./diffeq"
  "./geometry"
  //"./BlackHoles"
)

func main() {
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
