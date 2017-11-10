package ui

/**
* This package contains the ui printer that interfaces with the command line.
* It accepts a query from user input and uses it to update the core service.
 */

type UIController interface {
	setQuery(string)
	getQuery() string
}

type UIControllerImpl struct {
	query string
}

func (c *UIControllerImpl) setQuery(q string) {
	c.query = q
}

func (c *UIControllerImpl) getQuery() string {
	return c.query
}

