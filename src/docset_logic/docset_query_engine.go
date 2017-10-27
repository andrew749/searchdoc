package main

/**
* Interface a query engine conforms to
 */
type DocsetQueryEngine interface {
	GetIndicesForLanguage(language string) Docset
}

/**
* Concrete implementation of Query engine
 */
type DocsetQueryEngineImpl struct {
}

func (engine *DocsetQueryEngineImpl) GetIndicesForLanguage(language string) Docset {
	return DocsetForLanguage(language)
}
