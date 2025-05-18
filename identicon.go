package identicon

import (
	"crypto/md5"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
)

// GenerateIdenticon creates a GitHub style identicon from the
// provided integer and writes the PNG image to filePath.
func GenerateIdenticon(param int, filePath string) error {
	hash := md5.Sum([]byte(strconv.Itoa(param)))
	id := newIdenticon(hash[:])
	img := id.image()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

type identicon struct {
	source []byte
	size   int
}

func newIdenticon(src []byte) *identicon {
	return &identicon{source: src, size: 420}
}

func mapRange(value, vmin, vmax, dmin, dmax uint32) float32 {
	return float32(dmin) + float32(value-vmin)*float32(dmax-dmin)/float32(vmax-vmin)
}

func (i *identicon) foreground() color.RGBA {
	h := uint32((uint32(i.source[12]) << 8) | uint32(i.source[13]))
	s := uint32(i.source[14])
	l := uint32(i.source[15])

	hue := mapRange(h, 0, 4095, 0, 360)
	sat := mapRange(s, 0, 255, 0, 20)
	lum := mapRange(l, 0, 255, 0, 20)

	return hslToRGB(hue, 65.0-sat, 75.0-lum)
}

func (i *identicon) pixels() [25]bool {
	nib := newNibbler(i.source)
	var pixels [25]bool
	for col := 2; col >= 0; col-- {
		for row := 0; row < 5; row++ {
			ix := col + row*5
			mirrorCol := 4 - col
			mirrorIx := mirrorCol + row*5
			val, ok := nib.next()
			if !ok {
				val = 0
			}
			paint := val%2 == 0
			pixels[ix] = paint
			pixels[mirrorIx] = paint
		}
	}
	return pixels
}

func (i *identicon) image() *image.RGBA {
	pixelSize := 70
	spriteSize := 5
	margin := pixelSize / 2

	img := image.NewRGBA(image.Rect(0, 0, i.size, i.size))

	background := color.RGBA{240, 240, 240, 255}
	for y := 0; y < i.size; y++ {
		for x := 0; x < i.size; x++ {
			img.Set(x, y, background)
		}
	}

	fg := i.foreground()
	pixels := i.pixels()

	for row := 0; row < spriteSize; row++ {
		for col := 0; col < spriteSize; col++ {
			if pixels[row*spriteSize+col] {
				x := col*pixelSize + margin
				y := row*pixelSize + margin
				rect(img, x, y, x+pixelSize, y+pixelSize, fg)
			}
		}
	}

	return img
}

func rect(img *image.RGBA, x0, y0, x1, y1 int, c color.RGBA) {
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			img.Set(x, y, c)
		}
	}
}

// hslToRGB converts an HSL color value to RGB.
func hslToRGB(h, s, l float32) color.RGBA {
	hue := h / 360.0
	sat := s / 100.0
	lum := l / 100.0

	var b float32
	if lum <= 0.5 {
		b = lum * (sat + 1.0)
	} else {
		b = lum + sat - lum*sat
	}
	a := lum*2.0 - b

	r := hueToRGB(a, b, hue+1.0/3.0)
	g := hueToRGB(a, b, hue)
	bl := hueToRGB(a, b, hue-1.0/3.0)

	return color.RGBA{
		R: uint8(math.Round(float64(r * 255.0))),
		G: uint8(math.Round(float64(g * 255.0))),
		B: uint8(math.Round(float64(bl * 255.0))),
		A: 255,
	}
}

func hueToRGB(a, b, h float32) float32 {
	hue := h
	if hue < 0.0 {
		hue += 1.0
	} else if hue > 1.0 {
		hue -= 1.0
	}

	if hue < 1.0/6.0 {
		return a + (b-a)*6.0*hue
	}
	if hue < 1.0/2.0 {
		return b
	}
	if hue < 2.0/3.0 {
		return a + (b-a)*(2.0/3.0-hue)*6.0
	}
	return a
}

type nibbler struct {
	bytes []byte
	idx   int
	half  bool
}

func newNibbler(b []byte) *nibbler {
	return &nibbler{bytes: b}
}

func (n *nibbler) next() (byte, bool) {
	if n.idx >= len(n.bytes) {
		return 0, false
	}
	val := n.bytes[n.idx]
	if n.half {
		n.idx++
		n.half = false
		return val & 0x0f, true
	}
	n.half = true
	return (val & 0xf0) >> 4, true
}
