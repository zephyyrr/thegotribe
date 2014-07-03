package main

import (
	termbox "github.com/nsf/termbox-go"
	tet "github.com/zephyyrr/thegotribe"
	"log"
)

const (
	leftMarker  = termbox.ColorGreen
	rightMarker = termbox.ColorRed
)

func main() {
	termbox.Init()
	defer termbox.Close()

	et, err := tet.Create()

	if err != nil {
		log.Println(err)
		return
	}

	defer et.Close()
	if hb, err := et.Get("heartbeatinterval"); err == nil {
		log.Printf("Heartbeat Interval: %d", hb.(float64))
	}
	termbox.Clear(0, termbox.ColorWhite) //Set bg to white.
	go func() {
		for frame := range et.Frames {
			if frame.TrackingGaze() {
				termbox.Clear(0, termbox.ColorWhite)
				drawPoint(frame.LeftEye.Average, leftMarker)
				drawPoint(frame.RightEye.Average, rightMarker)
				termbox.Flush()
			}
		}
	}()

	done := make(chan struct{})

	go func() {
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventKey && event.Key == termbox.KeyEsc {
				close(done)
			}
		}
	}()
	<-done
	log.Println("Shutting down.")
}

func drawPoint(point tet.Point2D, attributes termbox.Attribute) {
	sizeX, sizeY := termbox.Size()
	x, y := int(float64(sizeX)*(float64(point.X)/float64(1366))),
		int(float64(sizeY)*(float64(point.Y)/float64(768)))
	switch {
	case x < 0:
		x = 0
	case x >= sizeX:
		x = sizeX - 1
	}

	switch {
	case y < 0:
		y = 0
	case y >= sizeY:
		y = sizeY - 1
	}
	termbox.CellBuffer()[sizeX*y+x] = termbox.Cell{'O', attributes, termbox.ColorWhite}
}
