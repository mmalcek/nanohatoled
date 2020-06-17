package nanohatoled

//display integration is based on github.com/mdp/monochromeoled project

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/imaging"
	"github.com/pbnjay/pixfont"
	"golang.org/x/exp/io/i2c"
)

const (
	// On or off registers.
	ssd1306DisplayOn  = 0xAf
	ssd1306DisplayOff = 0xAe

	// Scrolling registers.
	ssd1306ActivateScroll                   = 0x2F
	ssd1306DeactivateScroll                 = 0x2E
	ssd1306SetVerticalScrollArea            = 0xA3
	ssd1306RightHorizontalScroll            = 0x26
	ssd1306LeftHorizontalScroll             = 0x27
	ssd1306VerticalAndRightHorizontalScroll = 0x29
	ssd1306VerticalAndLeftHorizontalScroll  = 0x2A
)

// NanoOled - Current display
type NanoOled struct {
	dev *i2c.Device

	w             int    // width of the display
	h             int    // height of the display
	buf           []byte // each pixel is represented by a bit
	rotation      int
	rotationState bool
	image         *image.NRGBA
}

// Init sets up the display for writing
func (nanoOled *NanoOled) Init() (err error) {
	err = nanoOled.dev.Write([]byte{
		0xae,
		0x00 | 0x00, // row offset
		0x10 | 0x00, // column offset
		0xd5, 0x80,
		0xa8, uint8(nanoOled.h - 1),
		0xd3, 0x00, // set display offset to no offset
		0x80 | 0,
		0x8d, 0x14,
		0x20, 0x0,

		0xA0 | 0x1,
		0xC8,
	})
	if err != nil {
		return
	}
	if nanoOled.h == 32 {
		err = nanoOled.dev.Write([]byte{
			0xda, 0x02,
			0x81, 0x8f, // set contrast
		})
	}
	if nanoOled.h == 64 {
		err = nanoOled.dev.Write([]byte{
			0xda, 0x12,
			0x81, 0x7f, // set contrast
		})
	}
	err = nanoOled.dev.Write([]byte{
		0x9d, 0xf1,
		0xdb, 0x40,
		0xa4, 0xa6,

		0x2e,
		0xaf,
	})
	return
}

// Open - Initialize display
func Open() (*NanoOled, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-0"}, 0x3C)
	buf := make([]byte, 128*(64/8)+1)
	buf[0] = 0x40 // start frame of pixel data
	oled := &NanoOled{dev: dev, w: 128, h: 64, buf: buf}
	err = oled.Init()
	if err != nil {
		return nil, err
	}
	return oled, nil
}

// On turns on the display if it is off.
func (nanoOled *NanoOled) On() error {
	return nanoOled.dev.Write([]byte{ssd1306DisplayOn})
}

// Off turns off the display if it is on.
func (nanoOled *NanoOled) Off() error {
	return nanoOled.dev.Write([]byte{ssd1306DisplayOff})
}

// Close - Close display communication
func (nanoOled *NanoOled) Close() error {
	return nanoOled.dev.Close()
}

// New - Create new image for display
func (nanoOled *NanoOled) New(rotation int) {
	nanoOled.rotation = rotation
	nanoOled.rotationState = false
	nanoOled.Clear()
	if rotation == 90 || rotation == 270 {
		nanoOled.image = image.NewNRGBA(image.Rect(0, 0, 64, 128))
	} else {
		nanoOled.image = image.NewNRGBA(image.Rect(0, 0, 128, 64))
	}
}

// Send - draws an image on the display buffer starting from x, y.
func (nanoOled *NanoOled) Send() error {
	if nanoOled.rotationState == false {
		switch nanoOled.rotation {
		case 90:
			nanoOled.image = imaging.Rotate90(nanoOled.image)
			nanoOled.rotationState = true
		case 180:
			nanoOled.image = imaging.Rotate180(nanoOled.image)
			nanoOled.rotationState = true
		case 270:
			nanoOled.image = imaging.Rotate270(nanoOled.image)
			nanoOled.rotationState = true
		default:
		}
	}

	imgW := nanoOled.image.Bounds().Dx()
	imgH := nanoOled.image.Bounds().Dy()

	endX := 0 + imgW
	endY := 0 + imgH

	if endX >= nanoOled.w {
		endX = nanoOled.w
	}
	if endY >= nanoOled.h {
		endY = nanoOled.h
	}

	var imgI, imgY int
	for i := 0; i < endX; i++ {
		imgY = 0
		for j := 0; j < endY; j++ {
			r, g, b, _ := nanoOled.image.At(imgI, imgY).RGBA()
			var v byte
			if r+g+b > 0 {
				v = 0x1
			}
			if err := nanoOled.setPixel(i, j, v); err != nil {
				return err
			}
			imgY++
		}
		imgI++
	}
	nanoOled.draw()
	return nil
}

// Clear - Clear display and image
func (nanoOled *NanoOled) Clear() error {
	if nanoOled.rotation == 90 || nanoOled.rotation == 270 {
		nanoOled.image = image.NewNRGBA(image.Rect(0, 0, 64, 128))
	} else {
		nanoOled.image = image.NewNRGBA(image.Rect(0, 0, 128, 64))
	}
	for i := 1; i < len(nanoOled.buf); i++ {
		nanoOled.buf[i] = 0
	}
	return nanoOled.draw()
}

func (nanoOled *NanoOled) setPixel(x, y int, v byte) error {
	if x >= nanoOled.w || y >= nanoOled.h {
		return fmt.Errorf("(x=%v, y=%v) is out of bounds on this %vx%v display", x, y, nanoOled.w, nanoOled.h)
	}
	if v > 1 {
		return fmt.Errorf("value needs to be either 0 or 1; given %v", v)
	}
	i := 1 + x + (y/8)*nanoOled.w
	if v == 0 {
		nanoOled.buf[i] &= ^(1 << uint((y & 7)))
	} else {
		nanoOled.buf[i] |= 1 << uint((y & 7))
	}
	return nil
}

// draw draws the intermediate pixel buffer on the display.
func (nanoOled *NanoOled) draw() error {
	if err := nanoOled.dev.Write([]byte{
		0xa4,     // write mode
		0x40 | 0, // start line = 0
		0x21, 0, uint8(nanoOled.w),
		0x22, 0, 7,
	}); err != nil { // the write mode
		return err
	}
	return nanoOled.dev.Write(nanoOled.buf)
}

// Text - Write text to image
func (nanoOled *NanoOled) Text(x int, y int, text string) {
	pixfont.DrawString(nanoOled.image, x, y, text, color.White)
}

// Pixel - create pixel in image
func (nanoOled *NanoOled) Pixel(x int, y int, pixColor bool) {
	rColor := color.White
	if pixColor == false {
		rColor = color.Black
	}
	nanoOled.image.Set(x, y, rColor)
}

// LineH - create Horizontal line in image
func (nanoOled *NanoOled) LineH(x int, y int, length int) {
	px := x
	for px <= (length + x) {
		nanoOled.image.Set(px, y, color.White)
		px++
	}
}

// LineV - create Vertical line in image
func (nanoOled *NanoOled) LineV(x int, y int, length int) {
	py := y
	for py <= (length + y) {
		nanoOled.image.Set(x, py, color.White)
		py++
	}
}

// Rect - create rectangle in image
func (nanoOled *NanoOled) Rect(MinX int, MinY int, MaxX int, MaxY int, rectColor bool) {
	rColor := color.White
	if rectColor == false {
		rColor = color.Black
	}
	py := MinY
	for py <= MaxY {
		px := MinX
		for px <= MaxX {
			nanoOled.image.Set(px, py, rColor)
			px++
		}
		py++
	}
}
