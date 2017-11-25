package ui

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"searchdoc/src/html2text"
)

/**
* Renders HTML content in a text format
 */
type DocumentPane struct {
	NewContent      chan string
	Content         *string
	searchBarHeight int
}

func CreateDocumentPane(newContent chan string) DocumentPane {
	content := ""
	go func() {
		// keep waiting for new content
		for {
			content = <-newContent
		}
	}()

	return DocumentPane{newContent, &content, searchBarHeight}
}

func (mgr DocumentPane) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView("main", 30, mgr.searchBarHeight, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	v.Clear()
	w, _ := v.Size()
	opts := html2text.Options{PrettyTables: true, MaxLineLength: w - 1}
	text, err := html2text.FromString(*mgr.Content, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(v, text)

	v.Editable = false
	v.Wrap = true

	return nil
}
