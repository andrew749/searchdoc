package ui

/**
* This package should contain the ui printer that interfaces with the command line
* The goal is to accept input from the user and make a query to the core layer.
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
