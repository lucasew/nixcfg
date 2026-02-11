package api

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"

	"workspaced/pkg/api"
	"workspaced/pkg/config"
)

var ErrDriverNotFound = api.ErrDriverNotFound

// LAB represents a color in CIELAB color space
type LAB struct {
	L, A, B float64 // L: lightness [0-100], A: green-red, B: blue-yellow
}

// Polarity represents color scheme preference
type Polarity int

const (
	PolarityAny Polarity = iota
	PolarityDark
	PolarityLight
)

// Options configures palette extraction
type Options struct {
	Polarity   Polarity
	ColorCount int // 16 for base16, 24 for base24
	MaxSamples int // Limit pixels to sample (0 = all)
}

// Driver extracts color palettes from images
type Driver interface {
	Extract(ctx context.Context, img image.Image, opts Options) (*config.PaletteConfig, error)
	Name() string
}

// Color utility functions

// RGBToLAB converts RGB color to LAB for perceptual distance calculations
// Based on Stylix Data/Colour.hs rgb2lab function
func RGBToLAB(c color.Color) LAB {
	r, g, b, _ := c.RGBA()

	// Normalize to [0, 1]
	rf := float64(r) / 65535.0
	gf := float64(g) / 65535.0
	bf := float64(b) / 65535.0

	// Apply gamma correction (sRGB)
	rf = gammaCorrect(rf)
	gf = gammaCorrect(gf)
	bf = gammaCorrect(bf)

	// Convert to XYZ (D65 illuminant)
	x := rf*0.4124564 + gf*0.3575761 + bf*0.1804375
	y := rf*0.2126729 + gf*0.7151522 + bf*0.0721750
	z := rf*0.0193339 + gf*0.1191920 + bf*0.9503041

	// Normalize by reference white (D65)
	x /= 0.95047
	y /= 1.00000
	z /= 1.08883

	// Convert XYZ to LAB
	x = labF(x)
	y = labF(y)
	z = labF(z)

	l := 116.0*y - 16.0
	a := 500.0 * (x - y)
	bVal := 200.0 * (y - z)

	return LAB{L: l, A: a, B: bVal}
}

// gammaCorrect applies sRGB gamma correction
func gammaCorrect(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

// labF is the f(t) function used in XYZ to LAB conversion
func labF(t float64) float64 {
	delta := 6.0 / 29.0
	if t > delta*delta*delta {
		return math.Pow(t, 1.0/3.0)
	}
	return t/(3.0*delta*delta) + 4.0/29.0
}

// LABToRGB converts LAB back to RGB
// Based on Stylix Data/Colour.hs lab2rgb function
func LABToRGB(lab LAB) color.RGBA {
	// Convert LAB to XYZ
	fy := (lab.L + 16.0) / 116.0
	fx := lab.A/500.0 + fy
	fz := fy - lab.B/200.0

	x := labFInv(fx) * 0.95047
	y := labFInv(fy) * 1.00000
	z := labFInv(fz) * 1.08883

	// Convert XYZ to RGB (D65 illuminant)
	r := x*3.2404542 + y*-1.5371385 + z*-0.4985314
	g := x*-0.9692660 + y*1.8760108 + z*0.0415560
	b := x*0.0556434 + y*-0.2040259 + z*1.0572252

	// Apply inverse gamma correction
	r = gammaInverse(r)
	g = gammaInverse(g)
	b = gammaInverse(b)

	// Clamp and convert to 8-bit
	return color.RGBA{
		R: uint8(clamp(r*255.0, 0, 255)),
		G: uint8(clamp(g*255.0, 0, 255)),
		B: uint8(clamp(b*255.0, 0, 255)),
		A: 255,
	}
}

// labFInv is the inverse of labF
func labFInv(t float64) float64 {
	delta := 6.0 / 29.0
	if t > delta {
		return t * t * t
	}
	return 3.0 * delta * delta * (t - 4.0/29.0)
}

// gammaInverse applies inverse sRGB gamma correction
func gammaInverse(v float64) float64 {
	if v <= 0.0031308 {
		return v * 12.92
	}
	return 1.055*math.Pow(v, 1.0/2.4) - 0.055
}

// clamp restricts a value to [min, max]
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// DeltaE calculates perceptual color distance (CIE76)
func DeltaE(c1, c2 LAB) float64 {
	dl := c1.L - c2.L
	da := c1.A - c2.A
	db := c1.B - c2.B
	return math.Sqrt(dl*dl + da*da + db*db)
}

// ToHex converts color.Color to hex string (without #)
func ToHex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// Lightness extracts lightness value from color
func Lightness(c color.Color) float64 {
	lab := RGBToLAB(c)
	return lab.L
}

// SampleImage extracts unique colors from image, limited by maxSamples
func SampleImage(img image.Image, maxSamples int) []color.RGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	totalPixels := width * height

	// Calculate sampling stride
	stride := 1
	if maxSamples > 0 && totalPixels > maxSamples {
		stride = int(math.Ceil(float64(totalPixels) / float64(maxSamples)))
	}

	// Use map to track unique colors
	colorMap := make(map[uint32]color.RGBA)

	for y := bounds.Min.Y; y < bounds.Max.Y; y += stride {
		for x := bounds.Min.X; x < bounds.Max.X; x += stride {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()

			// Skip transparent pixels
			if a == 0 {
				continue
			}

			// Convert to 8-bit RGBA
			rgba := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: 255,
			}

			// Use packed uint32 as key for uniqueness
			key := uint32(rgba.R)<<16 | uint32(rgba.G)<<8 | uint32(rgba.B)
			colorMap[key] = rgba
		}
	}

	// Convert map to slice
	colors := make([]color.RGBA, 0, len(colorMap))
	for _, c := range colorMap {
		colors = append(colors, c)
	}

	return colors
}
