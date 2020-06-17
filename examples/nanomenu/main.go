package main

import (
	"fmt"
	"time"

	"github.com/mmalcek/nanohatoled"
)

var menu struct {
	cursor  int
	submenu bool
}

func main() {
	nanoHat, err := nanohatoled.Open()
	if err != nil {
		fmt.Println("error: ", err)
	}
	defer nanoHat.Close()

	nanoHat.New(0)                      // Create new empty image, rotation 0,90,180,270.
	err = nanoHat.Image("marioimg.png") // Load image from path
	if err != nil {
		fmt.Println("err Load image:", err)
	}

	nanoHat.Send() // Send image to screen
	time.Sleep(2000 * time.Millisecond)

	nanoHat.Rect(55, 0, 128, 64, false)   // Clear rect on screen
	nanoHat.Text(55, 10, "LOADING", true) // X, Y, string, color (true-white, false=black)

	for i := 0; i <= 60; i = i + 3 { // Progress bar
		nanoHat.Rect(55, 20, i+55, 30, true)
		nanoHat.Send()
		time.Sleep(50 * time.Millisecond)
	}

	nanoHat.Clear()
	mainMenu(nanoHat, false, false, false)

	// Watch buttons state
	go watchBtn(nanoHat, 0)
	go watchBtn(nanoHat, 1)
	go watchBtn(nanoHat, 2)

	for { // Keep app running
		//fmt.Println("myProgram")
		time.Sleep(5000 * time.Millisecond)
	}
}

func mainMenu(nanoHat *nanohatoled.NanoOled, btn1 bool, btn2 bool, btn3 bool) {
	if menu.submenu == false {
		if btn1 == true {
			menu.cursor--
			nanoHat.LineH(2, ((menu.cursor+1)*12)+9, 80, false)
		}
		if btn2 == true {
			menu.cursor++
			nanoHat.LineH(2, ((menu.cursor-1)*12)+9, 80, false)
		}
		if menu.cursor < 0 {
			menu.cursor = 3
			nanoHat.LineH(2, ((menu.cursor+1)*12)+9, 80, false)
		}
		if menu.cursor > 3 {
			menu.cursor = 0
			nanoHat.LineH(2, ((menu.cursor-1)*12)+9, 80, false)
		}

		nanoHat.LineH(2, ((menu.cursor * 12) + 9), 80, true)
		nanoHat.Text(15, 0, "Menu 1", true)
		nanoHat.Text(15, 12, "Menu 2", true)
		nanoHat.Text(15, 24, "Menu 3", true)
		nanoHat.Text(15, 36, "Menu 4", true)

		nanoHat.Text(0, 53, "Up", true)
		nanoHat.Text(50, 53, "Down", true)
		nanoHat.Text(110, 53, "OK", true)
		nanoHat.Send()

		if btn3 == true {
			if menu.cursor == 3 {
				nanoHat.Clear()
				nanoHat.Text(15, 20, "Menu4-OPEN", true)
				nanoHat.Send()
				menu.submenu = true
				return
			}
		}
	} else {
		if btn3 == true {
			menu.submenu = false
			nanoHat.Clear()
			mainMenu(nanoHat, true, false, false)
		}
	}
}

// Watch buttons state and call menu on change
func watchBtn(nanoHat *nanohatoled.NanoOled, button int) {
	for {
		nanoHat.Btn[button].WaitForEdge(-1)
		if button == 0 {
			mainMenu(nanoHat, true, false, false)
		}
		if button == 1 {
			mainMenu(nanoHat, false, true, false)
		}
		if button == 2 {
			mainMenu(nanoHat, false, false, true)
		}
		fmt.Printf("-> %s\n", nanoHat.Btn[button])
	}
}
