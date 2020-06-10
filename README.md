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

```go
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mmalcek/nanohatoled"
)

func main() {

	nanoImg, err := nanohatoled.Init()
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer nanoImg.Close()
	nanoImg.New(0) // Create new empty image, rotation 0,90,180,270,
	progress := 15
	nanoImg.Text(3, 0, "Hello World!!!") // X, Y, string
	nanoImg.Text(3, 11, "TIME: ")
	go showTime(nanoImg)
	nanoImg.LineH(0, 25, 128) // X, Y, length

	nanoImg.Rect(5, 30, 125, 50, true) // Create frame using rectangles
	nanoImg.Rect(7, 32, 123, 48, false)

	nanoImg.LineV(5, 55, 20)   // X, Y, length
	nanoImg.LineV(125, 55, 20) // X, Y, length
	for {
		progress = progress + 5
		if progress > 110 {
			nanoImg.Rect(10, 35, 115, 45, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			progress = 15
		}
		nanoImg.Rect(10, 35, progress, 45, true) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
		nanoImg.Text(10, 55, "Progress:")
		nanoImg.Rect(95, 55, 119, 61, false) // Clean progress from image
		nanoImg.Text(95, 55, strconv.Itoa(progress))
		nanoImg.Pixel(25, 45, true) //X, Y, true=White, false=Black color
		nanoImg.Send()              // Send image to screen
		time.Sleep(1000 * time.Millisecond)
	}
}

func showTime(nanoImg *nanohatoled.NanoImg) {
	for {
		nanoImg.Rect(50, 11, 120, 22, false) // Clear time area before showing new
		nanoImg.Text(50, 11, time.Now().Format("15:04:05"))
		time.Sleep(200 * time.Millisecond)
		nanoImg.Send()
	}
}

```
