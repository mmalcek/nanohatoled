# nanohatoled

## GO lang package for NanoHat OLED display

Display wiki description [NanoHat OLED display](http://wiki.friendlyarm.com/wiki/index.php/NanoHat_OLED)

- Very basic but easy to use
- Allow rotate screen, display up to 6 lines of text in horizontal view mode, basic shapes - line, rectangle
- Buttons controll implemented
- Tested with NanoPi NEO 2

![example screen](https://github.com/mmalcek/nanohatoled/blob/master/nanohat.jpg?raw=true)

### Example

- Start with "go get github.com/mmalcek/nanohatoled"
- Don't forget allow i2c0 in armbian-config (or other OS config)
- Application must run in sudo

```go
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mmalcek/nanohatoled"
)

var testCounter int

func main() {
	nanoHat, err := nanohatoled.Open()
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer nanoHat.Close()
	nanoHat.New(90) // Create new empty image, rotation 0,90,180,270.

	nanoHat.Text(10, 30, "START", true) // X, Y, string to write, color (true-white, false-black)
	nanoHat.Send()                      // Send image to screen
	time.Sleep(2000 * time.Millisecond)
	nanoHat.Clear() // Clear Image buffer and screen

	nanoHat.New(0) // Create new empty image, rotation 0,90,180,270.
	cursorY := 25
	for cursorY >= 0 {
		nanoHat.Rect(3, cursorY, 125, cursorY+12, false) // Clear previous text
		nanoHat.Text(3, cursorY, "Hello World!!!", true) // X, Y, string, color (true-white, false-black)
		nanoHat.Send()
		time.Sleep(50 * time.Millisecond)
		cursorY--
	}

	nanoHat.LineH(0, 25, 128, true) // X, Y, length - Horizontal line, color (true-white, false-black)

	nanoHat.Rect(5, 30, 125, 50, true) // Create frame using rectangles
	nanoHat.Rect(7, 32, 123, 48, false)
	nanoHat.LineV(5, 55, 20, true) // X, Y, length - Vertical line, color (true-white, false-black)
	nanoHat.Text(10, 55, "C:", true)
	nanoHat.Text(60, 55, "Prog:", true)
	nanoHat.LineV(125, 55, 20, true)
	go showProgress(nanoHat)

	// Watch buttons state
	go watchBtn(nanoHat, 0)
	go watchBtn(nanoHat, 1)
	go watchBtn(nanoHat, 2)

	nanoHat.Text(3, 11, "TIME: ", true)
	nanoHat.Pixel(124, 18, true) //X, Y, (0-black, 1White)

	var timeBuffer string
	for {
		if timeBuffer != time.Now().Format("15:04:05") {
			timeBuffer = time.Now().Format("15:04:05")
			nanoHat.Rect(50, 11, 120, 22, false) // Clear time area before showing new
			nanoHat.Text(50, 11, time.Now().Format("15:04:05"), true)
			nanoHat.Send()
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func showProgress(nanoHat *nanohatoled.NanoOled) {
	progress := 0
	for {
		progress = progress + 2
		if progress >= 100 {
			nanoHat.Rect(10, 35, 115, 45, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			progress = 0
		}
		nanoHat.Rect(10, 35, progress+15, 45, true) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
		nanoHat.Rect(105, 55, 120, 61, false)       // Clean progress bar from image
		nanoHat.Text(105, 55, strconv.Itoa(progress), true)
		nanoHat.Send()
		time.Sleep(100 * time.Millisecond)
	}
}

// Watch buttons state and update counter on change
func watchBtn(nanoHat *nanohatoled.NanoOled, button int) {
	for {
		nanoHat.Btn[button].WaitForEdge(-1)
		if button == 0 {
			testCounter++
			nanoHat.Rect(37, 55, 54, 64, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			nanoHat.Text(37, 55, strconv.Itoa(testCounter), true)
		}
		if button == 1 {
			testCounter--
			nanoHat.Rect(37, 55, 54, 64, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			nanoHat.Text(37, 55, strconv.Itoa(testCounter), true)
		}
		if button == 2 {
			testCounter = 0
			nanoHat.Rect(37, 55, 54, 64, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			nanoHat.Text(37, 55, "0", true)
		}
		fmt.Printf("-> %s\n", nanoHat.Btn[button])
	}
}
```

- Example video
[![Nanohatoled](http://img.youtube.com/vi/dXVI3RB2pK0/0.jpg)](https://www.youtube.com/watch?v=f843eVuxJVc "NanoHat OLED example")

- Menu example video
[![Nanohatoled](http://img.youtube.com/vi/dXVI3RB2pK0/0.jpg)](https://www.youtube.com/watch?v=rb1zSpImd_0 "NanoHat OLED menu")
