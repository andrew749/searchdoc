package ui

import (
    "fmt"
    "log"

	"searchdoc/src/html2text"
    "github.com/jroimartin/gocui"
)

func nextView(g *gocui.Gui, v *gocui.View) error {
    if v == nil || v.Name() == "side" {
        _, err := g.SetCurrentView("main")
        return err
    }
    _, err := g.SetCurrentView("side")
    return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy+1); err != nil {
            ox, oy := v.Origin()
            if err := v.SetOrigin(ox, oy+1); err != nil {
                return err
            }
        }
    }
    return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        cx, cy := v.Cursor()
        if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
            if err := v.SetOrigin(ox, oy-1); err != nil {
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
            return nil;
        }
        v.SetOrigin(ox, oy-1)
    }
    return nil
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        ox, oy := v.Origin()
        _, h := v.Size();
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

    maxX, maxY := g.Size()
    if v, err := g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        fmt.Fprintln(v, l)
        if _, err := g.SetCurrentView("msg"); err != nil {
            return err
        }
    }
    return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
    if err := g.DeleteView("msg"); err != nil {
        return err
    }
    if _, err := g.SetCurrentView("side"); err != nil {
        return err
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("side", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyCtrlSpace, gocui.ModNone, nextView); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'j', gocui.ModNone, scrollDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("main", 'k', gocui.ModNone, scrollUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
        return err
    }
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        return err
    }
    if err := g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
        return err
    }
    if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
        return err
    }

    return nil
}

// State variables.
var content string = ""
var queryResults []string

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("side", -1, -1, 30, maxY); err != nil {
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
