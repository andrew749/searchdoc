package docset_logic

import (
	data_models "searchdoc/src/data_models"
)

/**
* Interface a query engine conforms to
 */
type DocsetQueryEngine interface {
	GetIndicesForLanguage(language string) data_models.Docset
	GetDocsets() []string
}

/**
* Concrete implementation of Query engine
 */
type DocsetQueryEngineImpl struct {
}

func (engine *DocsetQueryEngineImpl) GetIndicesForLanguage(language string) data_models.Docset {
	return DocsetForLanguage(language)
}

func (engine *DocsetQueryEngineImpl) GetDocsets() []string {
	return GetAvailableDocsets()
}
