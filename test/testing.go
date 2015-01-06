package test

//This contains functions that are helpful for doing tests. 

import "math"
import "math/rand"
import "time"
import "fmt"

var seed_set bool = false

func SetSeed(x int64) {
  rand.Seed(x)
  fmt.Println("seed set")
  seed_set = true
}

func setSeed() {
  rand.Seed( time.Now().UTC().UnixNano())
  seed_set = true
}

//Get some random values in a range. 
func RandInt(min int, max int) int {
  if !seed_set {setSeed()}
  return min + rand.Intn(1 + max - min)
}

func RandSign() int {
  if !seed_set {setSeed()}
  return 2 * rand.Intn(2) - 1
}

func RandFloat(min, max float64) float64 {
  if !seed_set {setSeed()}
  return min + (max - min) * rand.Float64()
}

func RandFloatVector(min, max float64, length int) []float64 {
  v := make([]float64, length)
  for i := 0; i < length; i ++ {
    v[i] = RandFloat(min, max)
  }
  return v
}

//Test whether two floats are close enough to one another. 
func CloseEnough(a, b, e float64) bool {
  return math.Abs(a - b) < e
}

//Tests whether two lists are close enough to one another.
func VectorCloseEnough(a, b []float64, e float64) bool {
  if a == nil || b == nil {
    return false
  }
  if len(a) != len(b) {
    return false
  }

  for i := 0; i < len(a); i ++ {
    if !CloseEnough(a[i], b[i], e) {
      return false
    }
  }

  return true
}

//Tests whether a given number is close enough to
//any element in a list. 
func MemberCloseEnough(a float64, list []float64, e float64) (bool, int) {
  for i, elem := range list {
    if CloseEnough(a, elem, e) {
      return true, i
    }
  }
  return false, -1
}

//Tests whether two lists are close enough to one another.
func MatrixCloseEnough(a, b [][]float64, e float64) bool {
  if a == nil || b == nil {
    return false
  }
  if len(a) != len(b) {
    return false
  }

  for i := 0; i < len(a); i ++ {
    if !VectorCloseEnough(a[i], b[i], e) {
      return false
    }
  }

  return true
}
