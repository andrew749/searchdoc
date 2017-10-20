package main

/**
* This package should contain the ui printer that interfaces with the command line
* The goal is to accept input from the user and make a query to the core layer.
 */

type CoreUILayer interface {
	startRepl()
}
