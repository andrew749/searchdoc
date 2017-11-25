package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

type SearchBar struct {
	DidChange chan string // channel to send back events on
}

func CreateSearchBar(didChange chan string) SearchBar {
	return SearchBar{
		didChange,
	}
}

func (mgr SearchBar) Layout(g *gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView("searchBar", -1, -1, maxX, searchBarHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err := g.SetCurrentView("searchBar"); err != nil {
			return err
		}

		fmt.Fprintln(v, string(searchBuffer))

		v.Editable = true
		v.Editor = gocui.EditorFunc(simpleEditor(mgr.DidChange))
	}

	return nil
}

// simple function currying to bind the channel to this function so we can have feedback
func simpleEditor(didChange chan string) func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	return func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
		// check the key inputted and make sure it is a character
		switch {
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
			searchBuffer = append(searchBuffer, ch)
		case key == gocui.KeySpace:
			v.EditWrite(' ')
			searchBuffer = append(searchBuffer, ch)
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		}

		// get the value stored in the search bar
		query, err := v.Line(0)

		if err != nil {
			log.Fatal(err)
		}

		// tell a listener that we got new keystroke
		didChange <- query
	}
}
