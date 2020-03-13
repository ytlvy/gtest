package gif

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"time"
)

// Maker : gif image maker
type Maker struct {
}

//调色板 第一个为背景色 第二个为前景色
var pallette = []color.Color{color.RGBA{0xdd, 0xff, 0, 0xff}, color.RGBA{0xff, 0, 0, 0xff}}

const (
	whiteIndex = 0
	blackIndex = 1
)

// MakeGif make a gif
func (g *Maker) MakeGif(out io.Writer) {
	rand.Seed(time.Now().UTC().UnixNano())
	g.lissajous(out)
}

func (g *Maker) lissajous(out io.Writer) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		size    = 100
		nframes = 64
		delay   = 8
	)

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, pallette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*size+0.5), size+int(y*size+0.5), blackIndex)
		}
		phase += 0.1

		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}

	gif.EncodeAll(out, &anim)
}
