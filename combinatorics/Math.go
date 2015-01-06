package combinatorics

//The first 20 factorials, which is all that can fit in a 64-bit int.
var factorial []uint64 = []uint64{
  1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800, 39916800, 
  479001600, 6227020800, 87178291200, 1307674368000, 20922789888000, 
  355687428096000, 6402373705728000, 121645100408832000, 2432902008176640000}

//the classic factorial. 
func Factorial(n uint) uint64 {
  if n > 20 {
    return 0
  } else {
    return factorial[n]
  }
}

//Binomial coefficients.
func Binomial(d, n uint) uint64 {
  if n > d {
    return 0
  } else if d <= 20 {
    //In this case we can just use the list of factorials directly.
    return factorial[d] / (factorial[n] * factorial[d - n])

  } else if d <= 62 {
    //There are some binomial coefficients that can fit in a uint64, 
    //but can't be calculated with factorials that can fit, so these
    //are computed iteratively. 
    var b uint64 = 1
    var p uint64
    var l uint
    if 2 * n > d {
      p = uint64(n + 1)
      l = d - n
    } else {
      p = uint64(d - n + 1)
      l = n
    }

    for i := uint(1); i <= l; i ++ {
      b *= p
      b /= uint64(i)
      p ++
    }
    return b

  } else {
    return 0
  }
}

func Figurate(r, n uint) uint64 {
  if r == 0 || n == 0 { return 0 }
  return Binomial(n + r - 1, r)
}

//function must take a SORTED list. 
func Permutations(list []uint) uint64 {
  b := factorial[len(list)]
  n := uint64(1)
  for i := 1; i < len(list); i ++ {
    if list[i - 1] == list[i] {
      n ++ 
      b /= n
    } else {
      n = 1
    }
  }

  return b
}

//Go does not have an integer power function
//apparently, so here is one.
func Power(n int, p uint) int64 {
  var s int64 = int64(n)
  var r int64 = 1
  var q uint = p

  for q > 0 {
    if q & 1 == 1 {
      r *= s
    }

    s = s*s

    q = q >> 1
  }

  return r
}
