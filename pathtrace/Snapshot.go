package pathtrace

import "image"
import "image/color"
import "math"
import "fmt"
import "../distributions"

//Snap a photo! 
func Snapshot(scene *Scene, cam_func GenerateRay, size_u, size_v,
  depth, minp, maxp int, maxMeanVariance float64,
  minPercentNotification float64, minIterationNotification int) *image.NRGBA {
  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))
  var ray_pos, ray_dir []float64

  pix_sum := make([]float64, 3)
  pix := make([]float64, 3)

  var pixels, iterations, iteration_count, total_pixels int = 0, 0, 0, size_u * size_v
  var percent, percent_monitor float64 = 0, 0

  //Write to the screen information about the progress of the picture. 
  notify := func() {
    fmt.Println(percent, " complete. ", pixels, " pixels ", iterations, "iterations",
      float64(iterations) / float64(pixels), "iterations per pixel.")
    iteration_count = 0
    percent_monitor = 0
  }

  for i := 0; i < size_u; i ++ {
    for j := 0; j < size_v; j ++ {

      for k := 0; k < 3; k ++ {
        pix_sum[k] = 0
      }

      var p int = 0
      var variance_check bool

      //Set up the variance monitor.
      var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics()}

      for {
        //Set up the ray.
        ray_pos, ray_dir = cam_func(i, j)

        //Trace the path.
        c := scene.TracePath(ray_pos, ray_dir, depth, 1./256.)

        p ++
        iterations ++
        iteration_count ++

        if iteration_count == minIterationNotification {notify()}

        //Add the pixel to the variance monitor.
        for l := 0; l < 3; l ++ {
          pix_sum[l] += c[l]
          monitor[l].AddVariable(c[l])
        }

        if p > maxp {
          break
        }

        //Check variance
        if p > minp {
          variance_check = true
          for l := 0; l < 3; l ++ {
            if monitor[l].MeanVariance() > maxMeanVariance {
              variance_check = false
            }
          }
          if variance_check {
            break
          }
        }
      }

      //Generate the pixel. 
      for l := 0; l < 3; l ++ {
        pix[l] = math.Min(255 * pix_sum[l] / float64(p), 255)
      }

      img.Set(i, j, &color.NRGBA{uint8(pix[0]), uint8(pix[1]), uint8(pix[2]), 255})

      pixels ++
      percent = float64(pixels)/float64(total_pixels)
      percent_monitor += .1/float64(total_pixels)
      if percent_monitor > minPercentNotification {notify()}
    }
  }

  return img
}
