package munsell

import (
	"fmt"
	"math"
	"strconv"
)

type Color int

const (
	Unknown Color = iota
	White
	Black
	Red
	Orange
	Yellow
	Green
	LightBlue
	Blue
	Purple
)

type RGBComponents struct {
	Red, Green, Blue uint8
}

type HSLComponents struct {
	Hue, Saturation, Lightness float64
}

func (c Color) String() string {
	colors := [...]string{"Unknown", "White", "Black", "Red", "Orange", "Yellow", "Green", "Light Blue", "Blue", "Purple"}
	if c < Unknown || c > Purple {
		return fmt.Sprintf("Color(%d)", int(c))
	}
	return colors[c]
}

// Returns a color from Color when passed a hex color code, defaults to Unknown
func GetColorFromHex(hexColor string) (color Color, err error) {
	if hexColor[0] == '#' {
		hexColor = hexColor[1:]
	}
	switch len(hexColor) {
	case 3:
		// shorthand
		rgb := getRGBFromHex(hexColor[0:1] + hexColor[0:1] + hexColor[1:2] + hexColor[1:2] + hexColor[2:3] + hexColor[2:3])
		color = GetColorFromRGB(rgb)
	case 6:
		// standard
		rgb := getRGBFromHex(hexColor)
		color = GetColorFromRGB(rgb)
	case 8:
		// alpha component
		rgb := getRGBFromHex(hexColor[0:6])
		color = GetColorFromRGB(rgb)
	default:
		return Unknown, fmt.Errorf("invalid hex color code: %s", hexColor)
	}
	return color, nil
}

func ensureHexIsValid(hexColor string) bool {
	return true
}

func getRGBFromHex(hexColor string) RGBComponents {
	red, err := strconv.ParseUint(hexColor[0:2], 16, 8)
	if err != nil {
		panic(err)
	}
	green, err := strconv.ParseUint(hexColor[2:4], 16, 8)
	if err != nil {
		panic(err)
	}
	blue, err := strconv.ParseUint(hexColor[4:6], 16, 8)
	if err != nil {
		panic(err)
	}
	return RGBComponents{
		Red:   uint8(red),
		Green: uint8(green),
		Blue:  uint8(blue),
	}
}

func GetColorFromRGB(rgb RGBComponents) Color {
	hsl := getHSLFromRGB(rgb)
	color := matchColorFromHSL(hsl)
	return color
}

func getHSLFromRGB(rgb RGBComponents) HSLComponents {
	// taken from https://stackoverflow.com/a/39147465
	r := float64(rgb.Red)
	g := float64(rgb.Green)
	b := float64(rgb.Blue)

	r /= 255
	g /= 255
	b /= 255

	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)
	c := max - min

	hue := float64(0)

	if c == 0 {
		hue = 0
	} else {
		switch max {
		case r:
			segment := (g - b) / c
			shift := float64(0) // R° / (360° / hex sides)
			if segment < 0 {    // hue > 180, full rotation
				shift = 360 / 60 // R° / (360° / hex sides)
			}
			hue = segment + shift

		case g:
			segment := (b - r) / c
			shift := float64(120 / 60) // G° / (360° / hex sides)
			hue = segment + shift

		case b:
			segment := (r - g) / c
			shift := float64(240 / 60) // B° / (360° / hex sides)
			hue = segment + shift
		}
	}

	hue *= 60 // hue is in [0,6], scale it up

	return HSLComponents{
		Hue:        hue,
		Saturation: c * 100,
		Lightness:  (max + min) * 50,
	}
}

func GetColorFromHSL(hsl HSLComponents) Color {
	color := matchColorFromHSL(hsl)
	return color
}

func matchColorFromHSL(hsl HSLComponents) Color {
	l := math.Floor(hsl.Lightness)
	s := math.Floor(hsl.Saturation)
	h := math.Floor(hsl.Hue)

	if s <= 10 && l >= 90 {
		return White
	} else if l <= 13 {
		return Black
	} else if (s <= 10 && l <= 70) || s == 0 {
		return Black // ("Gray")
	} else if (h >= 0 && h <= 16) || h >= 346 {
		return Red
	} else if h > 16 && h <= 36 {
		if s < 90 {
			return Orange // ("Brown")
		} else {
			return Orange
		}
	} else if h > 36 && h <= 64 {
		if s < 90 {
			return Yellow // ("Brown")
		} else {
			return Yellow
		}
	} else if h > 64 && h <= 165 {
		return Green
	} else if h > 165 && h <= 208 {
		return LightBlue
	} else if h > 208 && h <= 260 {
		return Blue
	} else if h > 260 && h <= 290 {
		return Purple
	} else if h > 290 && h <= 345 {
		return Purple // Pink
	}
	return Unknown
}
