package ui

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestSetGetQuery(t *testing.T) {
    controller := UIControllerImpl{}
    query := "some query"
    controller.setQuery(query)
    assert.Equal(t, query, controller.getQuery())
}

func TestTermboxView(t *testing.T) {
    view := TermboxView{}
    v.render(GraphicsContext{})
}
