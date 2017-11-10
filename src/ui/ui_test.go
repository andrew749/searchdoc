package ui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetGetQuery(t *testing.T) {
	controller := UIControllerImpl{}
	query := "some query"
	controller.setQuery(query)
	assert.Equal(t, query, controller.getQuery())
}
