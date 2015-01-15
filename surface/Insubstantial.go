package surface

import "strings"
import "fmt"
import "math/rand"
import "math"
import "../vector"

//The insubstantial "surface" is basically a mist that fills up all of space.
//A ray randomly intersects with it over some distance. 
type insubstantial struct {
  dimension int
  ltau float64
}

func (i *insubstantial) Dimension() int {
  return i.dimension
}

//Everywhere is inside this thing. 
func (i *insubstantial) F(x []float64) float64 {
  return 1
}

func (i *insubstantial) Intersection(x, v []float64) []float64 {
  d := vector.Length(vector.Minus(v, x))

  return []float64{math.Log(1./rand.Float64()) / (i.ltau * d)}
}

func (i *insubstantial) Gradient(x []float64) []float64 {
  return make([]float64, i.dimension)
}

func (i *insubstantial) String() string {
  return strings.Join([]string{"insubstantial{", fmt.Sprint(i.ltau), "}"}, "")
}

func NewInsubstantialSurface(dim int, tau float64) Surface {
  return &insubstantial{dim, math.Log(tau)}
}
