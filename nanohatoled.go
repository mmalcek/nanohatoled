package nanohatoled

import (
	"errors"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/mdp/monochromeoled"
	"github.com/pbnjay/pixfont"
	"golang.org/x/exp/io/i2c"
)

// NanoImg - Current display
type NanoImg struct {
	rotation      int
	rotationState bool
	display       *monochromeoled.OLED
	image         *image.NRGBA
}

// Init - Initialize display
func Init() (*NanoImg, error) {
	nanoImg := &NanoImg{}
	var err error
	nanoImg.display, err = monochromeoled.Open(&i2c.Devfs{Dev: "/dev/i2c-0"}, 0x3C, 128, 64)
	if err != nil {
		return nanoImg, errors.New("OLED connection fail")
	}
	return nanoImg, nil
}

// Close - Close display communication
func (nanoImg *NanoImg) Close() {
	nanoImg.display.Close()
}

// New - Create new image for display
func (nanoImg *NanoImg) New(rotation int) {
	nanoImg.rotation = rotation
	nanoImg.rotationState = false
	nanoImg.display.Clear()
	if rotation == 90 || rotation == 270 {
		nanoImg.image = image.NewNRGBA(image.Rect(0, 0, 64, 128))
	} else {
		nanoImg.image = image.NewNRGBA(image.Rect(0, 0, 128, 64))
	}
}

// Send - Send image to display
func (nanoImg *NanoImg) Send() {
	if nanoImg.rotationState == false {
		switch nanoImg.rotation {
		case 90:
			nanoImg.image = imaging.Rotate90(nanoImg.image)
			nanoImg.rotationState = true
		case 180:
			nanoImg.image = imaging.Rotate180(nanoImg.image)
			nanoImg.rotationState = true
		case 270:
			nanoImg.image = imaging.Rotate270(nanoImg.image)
			nanoImg.rotationState = true
		default:
		}
	}
	nanoImg.display.SetImage(0, 0, nanoImg.image)
	nanoImg.display.Draw()
}

// Clear - Clear display
func (nanoImg *NanoImg) Clear() {
	if nanoImg.rotation == 90 || nanoImg.rotation == 270 {
		nanoImg.image = image.NewNRGBA(image.Rect(0, 0, 64, 128))
	} else {
		nanoImg.image = image.NewNRGBA(image.Rect(0, 0, 128, 64))
	}
	nanoImg.display.Clear()
}

// Text - Write text on screen
func (nanoImg *NanoImg) Text(x int, y int, text string) {
	pixfont.DrawString(nanoImg.image, x, y, text, color.White)
}

// LineH - create Horizontal line
func (nanoImg *NanoImg) LineH(x int, y int, length int) {
	px := x
	for px <= length {
		nanoImg.image.Set(px, y, color.White)
		px++
	}
}

// LineV - create Vertical line
func (nanoImg *NanoImg) LineV(x int, y int, length int) {
	py := y
	for py <= length {
		nanoImg.image.Set(x, py, color.White)
		py++
	}
}

// Pixel - create pixel
func (nanoImg *NanoImg) Pixel(x int, y int, pixColor bool) {
	pColor := color.White
	if pixColor == false {
		pColor = color.Black
	}
	nanoImg.image.Set(x, y, pColor)
}

// Rect - create rectangle
func (nanoImg *NanoImg) Rect(MinX int, MinY int, MaxX int, MaxY int, rectColor bool) {
	rColor := color.White
	if rectColor == false {
		rColor = color.Black
	}
	py := MinY
	for py <= MaxY {
		px := MinX
		for px <= MaxX {
			nanoImg.image.Set(px, py, rColor)
			px++
		}
		py++
	}

}
