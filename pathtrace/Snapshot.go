package pathtrace

import "image"
import "image/color"
import "math"
import "../distributions"
import "fmt"
import "time"

//Create a section of a photo. 
func snapSegment(scene *Scene, cam_func GenerateRay,
  size_u, v_min, v_max, depth, minp, maxp int, maxMeanVariance float64) [][][]float64 {

  section := make([][][]float64, v_max - v_min)

  var ray_pos, ray_dir []float64

  for i := v_min; i < v_max; i ++ {
    section[i - v_min] = make([][]float64, size_u)

    for j := 0; j < size_u; j ++ {
      pix := make([]float64, 3)

      var p int = 0
      var variance_check bool

      //Set up the variance monitor.
      var monitor []*distributions.SampleStatistics = []*distributions.SampleStatistics{
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics(), 
        distributions.NewSampleStatistics()}

      for {
        //Set up the ray.
        ray_pos, ray_dir = cam_func(j, i)

        //Trace the path.
        c := scene.TracePath(ray_pos, ray_dir, depth, 1./256.)

        p ++
        //iterations ++
        //iteration_count ++

        //Add the pixel to the variance monitor.
        for l := 0; l < 3; l ++ {
          pix[l] += c[l]
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
        pix[l] = math.Min(255 * pix[l] / float64(p), 255)
      }

      section[i - v_min][j] = pix
    }
  }

  return section
}

//A data structure used to pass information over a channel from
//one goroutine to another. 
type image_slice struct {
  v_min, v_max int
  pix [][][]float64
}

//Snap a photo! 
func Snapshot(sceneBuild func() *Scene, cam_func GenerateRay, size_u, size_v, depth, minp, maxp int,
  maxMeanVariance float64, minPercentNotification float64,
  minIterationNotification, routines int) *image.NRGBA {
  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))

  //Write to the screen information about the progress of the picture. 
  /*notify := func() {
    fmt.Println(percent, " complete. ", pixels, " pixels ", iterations, "iterations",
      float64(iterations) / float64(pixels), "iterations per pixel.")
    iteration_count = 0
    percent_monitor = 0
  }*/

  //Set up the wait group and channels. 
  height := 5
  slices := size_v / height
  if size_v % height != 0 { slices ++ }
  ch_in := make(chan []int, slices)
  ch_out := make(chan *image_slice)

  //Set up the parameters for the concurrent evaluation
  //and send them into the channel.
  for i := 0; i < size_v; i += height {
    param := make([]int, 2)
    param[0] = i
    if i + height > size_v {
      param[1] = size_v
    } else {
      param[1] = i + height
    }
    fmt.Println("slice ", i)

    ch_in <- param
  }
  fmt.Println("got this far - ugh", routines)

  //Start running the go routines.
  for i := 0; i < routines; i ++ {

    fmt.Println("about to start a routine.")
    go func(scene *Scene, ch_in chan []int, ch_out chan *image_slice) {
      fmt.Println("++ go routine started.")
      for param := range ch_in {
        fmt.Println("++ about to run slice", param[0], param[1])
        time.Sleep(time.Millisecond * 250)

        ch_out <- &image_slice{param[0], param[1],
          snapSegment(scene, cam_func, size_u, param[0], param[1], depth, minp, maxp, maxMeanVariance)}
        //ch_out <- &image_slice{param[0], param[1], [][][]float64{}}
      }
    } (sceneBuild(), ch_in, ch_out)
  }

  fmt.Println("got this far")

  //Create the image over several go routines. 
  received := 0
  for received < slices {
    select {
    case slice := <-ch_out :
      fmt.Println("got slice ", slice.v_min, slice.v_max)
      for i := slice.v_min; i < slice.v_max; i ++ {
        for j := 0; j < size_u; j ++ {
          pix := slice.pix[i - slice.v_min][j]
          img.Set(j, i, &color.NRGBA{uint8(pix[0]), uint8(pix[1]), uint8(pix[2]), 255})
        }
      }
      received ++
    }
  }

  return img
}

//Snap a photo! This is the old version, without multithreading. 
func SnapshotNoThreads(scene *Scene, cam_func GenerateRay, size_u, size_v,
  depth, minp, maxp int, maxMeanVariance float64,
  minPercentNotification float64, minIterationNotification int) *image.NRGBA {
  img := image.NewNRGBA(image.Rect(0, 0, size_u, size_v))                                           

  slice := snapSegment(scene, cam_func, size_u, 0, size_v, depth, minp, maxp, maxMeanVariance)

  for i := 0; i < size_v; i ++ {
    for j := 0; j < size_u; j ++ {
      pix := slice[i][j]
      img.Set(j, i, &color.NRGBA{uint8(pix[0]), uint8(pix[1]), uint8(pix[2]), 255})
    }
  }

  return img
}
