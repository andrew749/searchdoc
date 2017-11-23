package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"searchdoc/src/core"
	"searchdoc/src/data_models"
	"searchdoc/src/html2text"
)

// views
var (
	mainContentView *gocui.View
	searchBar       *gocui.View
	sideBarView     *gocui.View
)

// State variables.
var (
	content      string = ""
	queryResults []string
	language     string
	searchBuffer []rune
)

const (
	searchBarHeight = 1
)

var SearchEditor gocui.Editor = gocui.EditorFunc(simpleEditor)

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if sideBarView != nil {
		// get the cursor so we can determine if we need to scroll the pane as wel
		cx, cy := sideBarView.Cursor()
		if err := sideBarView.SetCursor(cx, cy+1); err != nil {
			ox, oy := sideBarView.Origin()
			// move the view down
			if err := sideBarView.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if sideBarView != nil {
		ox, oy := sideBarView.Origin()
		cx, cy := sideBarView.Cursor()
		if err := sideBarView.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := sideBarView.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

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
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	// exit the program
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	return nil
}

func simpleEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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

	query, err := v.Line(0)

	if err != nil {
		log.Fatal(err)
	}

	results, err := core.Query(query, language)
	queryResults = make([]string, 0)

	queryResults = convertResultsToStringSlice(results)

	// perform query
	SetQueryResults(queryResults)
}

func convertResultsToStringSlice(results []data_models.DocsetElement) []string {
	resultStrs := make([]string, len(results))
	for i, result := range results {
		resultStrs[i] = result.Name
	}
	return resultStrs
}

type SideManager struct{}

func (mgr SideManager) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()

	if v, err := g.SetView("side", -1, searchBarHeight, 30, maxY); err != nil {
		sideBarView = v

		if err != gocui.ErrUnknownView {
			return err
		}

		for _, result := range queryResults {
			fmt.Fprintln(v, result)
		}

		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	return nil
}

type MainManager struct{}

func (mgr MainManager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("main", 30, searchBarHeight, maxX, maxY); err != nil {
		mainContentView = v

		if err != gocui.ErrUnknownView {
			return err
		}
		w, _ := v.Size()
		opts := html2text.Options{PrettyTables: true, MaxLineLength: w - 1}
		text, err := html2text.FromString(content, opts)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(v, text)

		v.Editable = false
		v.Wrap = true
	}

	return nil
}

type SearchManager struct{}

func (mgr SearchManager) Layout(g *gocui.Gui) error {

	maxX, _ := g.Size()

	if v, err := g.SetView("searchBar", -1, -1, maxX, searchBarHeight); err != nil {
		searchBar = v
		fmt.Fprintln(v, string(searchBuffer))

		if err != gocui.ErrUnknownView {
			return err
		}

		if _, err := g.SetCurrentView("searchBar"); err != nil {
			return err
		}

		v.Editable = true
		v.Editor = SearchEditor
	}

	return nil
}

func SetHtmlContent(html string) {
	content = html
}

func SetQueryResults(results []string) {
	queryResults = results
}

func SetLanguage(lang string) {
	language = lang
}

func Init() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = false
	g.SetManager(SearchManager{}, MainManager{}, SideManager{})

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
