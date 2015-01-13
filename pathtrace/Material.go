package pathtrace

//From the parameters inside the ray object, derive the color of the ray.
func DeriveColor(receptor, emission []float64, redirected float64) (c []float64) {
  c = make([]float64, 3)
  for i := 0; i < 3; i ++ {
    c[i] = receptor[i] * redirected + emission[i]
  }
  return 
}

type ColorInteraction func(receptor, emission []float64,
  redirected float64) ([]float64, []float64, float64) 

//A function to make an object glow.
func Glow(c []float64) ColorInteraction {
  if c == nil {
    return nil
  } else {
    return func(receptor, emission []float64, redirected float64) ([]float64, []float64, float64) {
      for i := 0; i < 3; i ++ {
        emission[i] = receptor[i] * c[i] * redirected + emission[i]
      }
      return receptor, emission, 0
    }
  }
}

//A function to make an object absorb light.
func Absorb(c []float64) ColorInteraction {
  if c == nil {
    return nil
  } else {
    return func(receptor, emission []float64, redirected float64) ([]float64, []float64, float64) {
      for l := 0; l < 3; l ++ {
        receptor[l] *= c[l]
      }
      return receptor, emission, redirected
    }
  }
}

//A function for both.
func GlowAbsorbAverage(glow_color, transmit_color []float64, absorb float64) ColorInteraction {
  if glow_color == nil || transmit_color == nil {
    return nil
  } else {
    return func(receptor, emission []float64, redirected float64) ([]float64, []float64, float64) {
      for i := 0; i < 3; i ++ {
        emission[i] += redirected * absorb * glow_color[i]
        receptor[i] *= transmit_color[i]
      }
      redirected *= (1 - absorb)
      return receptor, emission, redirected
    }
  } 
}
