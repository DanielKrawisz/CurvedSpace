package surface

import "strings"
import "fmt"
import "math/rand"
import "math"
import "github.com/DanielKrawisz/CurvedSpace/vector"

var insubstantialRand func() float64 = rand.Float64

var log2 float64 = math.Log(2)

//The insubstantial "surface" is basically a mist that fills up all of space.
//A ray randomly intersects with it over some distance.
type insubstantial struct {
	dimension int
	tau       float64
}

func (i *insubstantial) Dimension() int {
	return i.dimension
}

//Everywhere is inside this thing.
func (i *insubstantial) F(x []float64) float64 {
	return 1
}

func (i *insubstantial) Intersection(x, v []float64) []float64 {
	d := vector.Length(v)

	return []float64{math.Log(1./insubstantialRand()) * i.tau / (log2 * d)}
}

func (i *insubstantial) Translate(x []float64) Surface {
	return i
}

func (i *insubstantial) CoordinateShift(x [][]float64) Surface {
	return i
}

func (i *insubstantial) Gradient(x []float64) []float64 {
	return make([]float64, i.dimension)
}

func (i *insubstantial) String() string {
	return strings.Join([]string{"insubstantial{", fmt.Sprint(i.tau), "}"}, "")
}

func NewInsubstantialSurface(dim int, tau float64) Surface {
	return &insubstantial{dim, tau}
}
