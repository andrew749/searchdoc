package ui

import (
    "main"

    "fmt"
    "github.com/nsf/termbox-go"
)

type GraphicsContext struct {
    x, y int
    fg, bg termbox.Attribute
}

type Renderable interface {
    render(GraphicsContext)
    children() []*Renderable
}

type BorderStyle uint16
const (
    BorderNone BorderStyle = iota
    BorderSingle
    BorderDouble
)
// Implements Renderable
type Border struct {
    top, bottom, left, right BorderStyle
}

func (b *Border) render(gc GraphicsContext) {
    fmt.Printf("Render border (%v): %v", gc, *b)
}

// Implements Renderable
type Panel struct {
    border Border
    width, height int
}

func (p *Panel) render(gc GraphicsContext) {
    fmt.Printf("Render panel (%v): %v", gc, *p)
}

// Implements Observer and Renderable
type TermboxView struct {
    query string
    matches string
    document string
    root Panel
}

func (v *TermboxView) notify(o main.Observable) {
    // TODO(ajklen)
}

func (v *TermboxView) render(gc GraphicsContext) {
    fmt.Printf("Render view (%v): %v", gc, *p)
}

