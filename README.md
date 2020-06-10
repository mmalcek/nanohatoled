# nanohatoled

## GO lang package for NanoHat OLED display

Display wiki description [NanoHat OLED display](http://wiki.friendlyarm.com/wiki/index.php/NanoHat_OLED)

- Very basic but easy to use
- Allow rotate screen, display up to 6 lines of text in horizontal view mode, basic shapes - line, rectangle
- Written over one night just for basic purposes so please excuse brevity
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

func main() {
	nanoHat, err := nanohatoled.Init()
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer nanoHat.Close()
	nanoHat.New(0) // Create new empty image, rotation 0,90,180,270.

	nanoHat.Text(10, 30, "INITIALIZING") // X, Y, string - write string to image
	nanoHat.Send()                       // Send image to screen
	time.Sleep(1000 * time.Millisecond)
	nanoHat.Clear() // Clear Image buffer and screen

	cursorY := 25
	for cursorY >= 0 {
		nanoHat.Rect(3, cursorY, 125, cursorY+12, false) // Clear previous text
		nanoHat.Text(3, cursorY, "Hello World!!!")       // X, Y, string
		nanoHat.Send()
		time.Sleep(50 * time.Millisecond)
		cursorY--
	}

	nanoHat.LineH(0, 25, 128) // X, Y, length - Horizontal line

	nanoHat.Rect(5, 30, 125, 50, true) // Create frame using rectangles
	nanoHat.Rect(7, 32, 123, 48, false)
	nanoHat.LineV(5, 55, 20) // X, Y, length - Vertical line
	nanoHat.Text(10, 55, "Progress:")
	nanoHat.LineV(125, 55, 20)
	go showProgress(nanoHat)

	nanoHat.Text(3, 11, "TIME: ")
	nanoHat.Pixel(124, 18, true) //X, Y, true=White, false=Black color. Single pixel

	for {
		nanoHat.Rect(50, 11, 120, 22, false) // Clear time area before showing new
		nanoHat.Text(50, 11, time.Now().Format("15:04:05"))
		nanoHat.Send()
		time.Sleep(300 * time.Millisecond)
	}
}

func showProgress(nanoHat *nanohatoled.NanoImg) {
	progress := 0
	for {
		progress = progress + 2
		if progress >= 100 {
			nanoHat.Rect(10, 35, 115, 45, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			progress = 0
		}
		nanoHat.Rect(10, 35, progress+15, 45, true) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
		nanoHat.Rect(95, 55, 119, 61, false)        // Clean progress bar from image
		nanoHat.Text(95, 55, strconv.Itoa(progress))
		nanoHat.Send()
		time.Sleep(100 * time.Millisecond)
	}
}

```

[![Nanohatoled](http://img.youtube.com/vi/dXVI3RB2pK0/0.jpg)](http://www.youtube.com/watch?v=dXVI3RB2pK0 "NanoHat OLED")
