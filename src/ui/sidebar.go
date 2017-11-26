package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"searchdoc/src/data_models"
)

type SideBar struct {
	ResultsChannel        chan []data_models.DocsetElement
	UpdateSelectedElement chan data_models.DocsetElement
	results               *[]data_models.DocsetElement
	Scroll                chan int // let us know when the user scrolled
	searchBarHeight       int
	width                 int
	direction             *int
	cursorPosition        *int
}

func CreateSideBar(resultsChannel chan []data_models.DocsetElement, scroll chan int, searchBarHeight int, width int) SideBar {
	cursorPosition := 0
	direction := 1
	var results []data_models.DocsetElement
	updateChannel := make(chan data_models.DocsetElement)

	go func() {
		// keep waiting for changes and updating the scroll state
		for {
			tempDirection := <-scroll
			if cursorPosition+direction >= -0 && cursorPosition+direction < len(results) {
				direction = tempDirection
				cursorPosition += direction
				updateChannel <- results[cursorPosition]
			}
		}
	}()

	go func() {
		for {
			results = <-resultsChannel
			cursorPosition = 0
			direction = 1
		}
	}()

	return SideBar{resultsChannel, updateChannel, &results, scroll, searchBarHeight, width, &direction, &cursorPosition}
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
		fmt.Fprintln(v, result.Name)
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
