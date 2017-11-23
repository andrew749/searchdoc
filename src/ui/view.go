package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"searchdoc/src/html2text"
)

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

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	// TODO: fix
	// load a new view
	fmt.Println(l)

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

	// select an option in the list and load the entry
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}

	return nil
}

var (
	mainContentView *gocui.View
	sideBarView     *gocui.View
)

// State variables.
var (
	content      string = ""
	queryResults []string
)

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
		sideBarView = v
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, result := range queryResults {
			fmt.Fprintln(v, result)
		}
	}
	if v, err := g.SetView("main", 30, -1, maxX, maxY); err != nil {
		mainContentView = v
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = false
		v.Wrap = true
		w, _ := v.Size()
		opts := html2text.Options{PrettyTables: true, MaxLineLength: w - 1}
		text, err := html2text.FromString(content, opts)
		if err != nil {
			return err
		}
		fmt.Fprint(v, text)
		if _, err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func SetHtmlContent(html string) {
	content = html
}

func SetQueryResults(results []string) {
	queryResults = results
}

func Init() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = false
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
