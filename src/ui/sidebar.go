package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

type SideBar struct {
	ResultsChannel  chan []string
	results         *[]string
	Scroll          chan int // let us know when the user scrolled
	searchBarHeight int
	width           int
	direction       *int
	cursorPosition  *int
}

func CreateSideBar(resultsChannel chan []string, scroll chan int, searchBarHeight int, width int) SideBar {
	cursorPosition := 0
	direction := 1
	var results []string

	go func() {
		// keep waiting for changes and updating the scroll state
		for {
			direction = <-scroll
			cursorPosition += direction
		}
	}()

	go func() {
		for {
			results = <-resultsChannel
		}
	}()

	return SideBar{resultsChannel, &results, scroll, searchBarHeight, width, &direction, &cursorPosition}
}

func (sb SideBar) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	v, err := g.SetView("sidebar", -1, sb.searchBarHeight, sb.width, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.Clear()

	for _, result := range *sb.results {
		fmt.Fprintln(v, result)
	}

	v.Highlight = true
	v.SelBgColor = gocui.ColorGreen
	v.SelFgColor = gocui.ColorBlack

	ox, oy := v.Origin()
	cx, _ := v.Cursor()

	switch *sb.direction {
	// handle a movement up
	case -1:
		if err := v.SetCursor(cx, *sb.cursorPosition); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}

		}
		// handle a movement down
	case 1:
		if err := v.SetCursor(cx, *sb.cursorPosition); err != nil {
			ox, oy := v.Origin()

			// move the view down
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	default:
		fmt.Println("default")
	}

	return nil
}
