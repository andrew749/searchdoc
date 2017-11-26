package ui

import (
	"log"

	"github.com/jroimartin/gocui"
	"searchdoc/src/core"
	"searchdoc/src/data_models"
)

// views
var (
	searchBar    SearchBar
	sideBar      SideBar
	documentPane DocumentPane
)

// State variables.
var (
	language          string
	searchBuffer      []rune
	lastDisplayedPage string
)

const (
	searchBarHeight = 1
)

func scrollUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		if oy <= 0 {
			return nil
		}
		v.SetOrigin(ox, oy-1)
	}
	return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		_, h := v.Size()
		// Disable infinite scrolling.
		if _, err := v.Line(h); err != nil {
			return nil
		}
		v.SetOrigin(ox, oy+1)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	// scroll the content pane down
	if err := g.SetKeybinding("", gocui.KeyCtrlJ, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	// scroll the content pane up
	if err := g.SetKeybinding("", gocui.KeyCtrlK, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	// navigate the option list
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursor(1)); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursor(-1)); err != nil {
		return err
	}

	// exit the program
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func cursor(direction int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		sideBar.Scroll <- direction
		return nil
	}
}

func Search(query string, language string) []data_models.DocsetElement {
	results, err := core.Query(query, language)

	if err != nil {
		log.Fatal(err)
	}

	return results
}

func convertResultsToStringSlice(results []data_models.DocsetElement) []string {
	resultStrs := make([]string, len(results))
	for i, result := range results {
		resultStrs[i] = result.Name
	}
	return resultStrs
}

func Init(lang string) {
	// set the language for this searching instance
	language = lang

	// create a new gui
	g, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		log.Panicln(err)
	}

	defer g.Close()

	g.Cursor = false

	// initialize views
	searchBar = CreateSearchBar(make(chan string))
	sideBar = CreateSideBar(make(chan []string), make(chan int), searchBarHeight, 30)
	documentPane = CreateDocumentPane(make(chan string))

	go func() {
		for {
			searchResults := Search(<-searchBar.DidChange, language)
			sideBarItems := convertResultsToStringSlice(searchResults)

			// update the sidebar items
			sideBar.ResultsChannel <- sideBarItems

			// update the document pane
			if len(searchResults) > 0 {
				data, err := core.GetDocContent(searchResults[*sideBar.cursorPosition], language)

				if err != nil {
					log.Fatal(err)
				}

				documentPane.NewContent <- data
			}
		}
	}()
	go func() {

	}()

	// setup views
	g.SetManager(searchBar, documentPane, sideBar)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
