package main

type Observer interface {
	notify(o Observable)
}

type Observable interface {
	// notify all observers
	addObserver(o Observer)

	// notify all objects watching this
	notifyObservers()
}

/**
* Implements Observable Interface
 */
type Model struct {
	observers []Observer
}

func (m *Model) addObserver(o Observer) {
	m.observers = append(m.observers, o)
}

func (m *Model) notifyObservers() {
	for _, observer := range m.observers {
		observer.notify(m)
	}
}
