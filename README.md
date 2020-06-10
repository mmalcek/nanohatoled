# nanohatoled

## GO lang package for NanoHat OLED display

Display wiki description  [NanoHat OLED display](http://wiki.friendlyarm.com/wiki/index.php/NanoHat_OLED)

* Very simple but easy to use package 
* Allow rotate screen, display up to 6 lines of text in wide mode, basic shapes - line, rectangle
* Written over one night just for basic purposes so please excuse brevity
* Tested with NanoPi NEO 2

![example screen](https://github.com/mmalcek/nanohatoled/blob/master/nanohat.jpg?raw=true)

### Example

```go
package main

import (
	"fmt"
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
	for {
		nanoImg.Text(3, 0, "Hello World!!!") // X, Y, string
		nanoImg.Text(3, 11, "TIME: ")
		go showTime(nanoImg)
		nanoImg.LineH(0, 25, 128)  // X, Y, length
		nanoImg.LineV(10, 25, 20)  // X, Y, length
		nanoImg.LineV(118, 25, 20) // X, Y, length
		progress = progress + 5
		if progress > 110 {
			nanoImg.Rect(10, 30, 115, 40, false) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
			progress = 15
		}
		nanoImg.Rect(10, 30, progress, 40, true) // Xmin, Ymin, Xmax, Ymax, true=White, false=Black color
		nanoImg.Text(15, 50, "Progress...")
		nanoImg.Pixel(5, 30, true) //X, Y, true=White, false=Black color
		nanoImg.Send()             // Send image to screen
		time.Sleep(1000 * time.Millisecond)
	}
}

func showTime(nanoImg *nanohatoled.NanoImg) {
	for {
		nanoImg.Rect(50, 11, 120, 22, false)
		nanoImg.Text(50, 11, time.Now().Format("15:04:05"))
		time.Sleep(500 * time.Millisecond)
		nanoImg.Send()
	}
}

```






