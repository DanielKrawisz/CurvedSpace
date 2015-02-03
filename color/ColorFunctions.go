package color

import "github.com/DanielKrawisz/CurvedSpace/vector"
import "github.com/DanielKrawisz/CurvedSpace/functions"

//A color is actually a function over the wavelengths of light. 
type Color func([]float64) []float64

//For when we're working in a set color model, like rgb. 
func PresetColor(color []float64) Color {
  return func ([]float64) []float64 {
    return color
  }
}

//These two are actually the same type, but the spherical color function
//is supposed to depend only on the direction of the vector given it.
//There is no real way to enforce that though. 
type SpacialColorFunction func([]float64) Color

type SphericalColorFunction func([]float64) Color

func ConstantColorFunction(f Color) (func([]float64) Color) {
  return func([]float64) Color {
   return f
  }
}

// v - a list of normalized vectors. 
// d - a list of dot products. 
// f - a list of colors.
// 
// may return nil.
func Spotlights(v[][] float64, d[]float64, f[] Color) SphericalColorFunction {
  if v == nil || d == nil || f == nil {return nil}
  if len(v) != len(d) || len(v) != len(f) {return nil}

  return func(direction []float64) Color {
    vector.Normalize(direction)

    return func(receptor []float64) []float64 {

      col := make([]float64, len(receptor)) 
      for i, vec := range v {
        if vector.Dot(vec, direction) > d[i] {
          c := f[i](receptor)

          for j := 0; j < len(col); j ++ {
            col[j] += c[j]
          }
        } 
      }
      return col
    }
  }
}

// v - a position vector that gives the origin point of the checks.
// c - coordinate vectors describing the orientation of the checks.
// a, b - the two colors. 
//
//may return nil.
func Checks(v[] float64, c[][] float64, a, b Color) SpacialColorFunction {
  f := functions.Checks(v, c) 
  if f == nil {return nil}

  return func(position []float64) Color {

    if f(position) > 0 {
      return a
    } else {
      return b
    } 
  }
}
