package test

import "math"
import "math/rand"
import "time"

var seed_set bool = false

func setSeed() {
  rand.Seed( time.Now().UTC().UnixNano())
  seed_set = true
}

func RandInt(min int, max int) int {
  if !seed_set {setSeed()}
  return min + rand.Intn(max-min+1)
}

func RandFloat(min, max float64) float64 {
  if !seed_set {setSeed()}
  return min + (max - min) * rand.Float64()
}

func CloseEnough(a, b, e float64) bool {
  return math.Abs(a - b) < e
}
