package docset_logic

import (
	data_models "searchdoc/src/data_models"
)

/**
* Interface a query engine conforms to
 */
type DocsetQueryEngine interface {
	GetIndicesForLanguage(language string) data_models.Docset
	GetDownloadedDocsets() []string
	GetDownloadableDocsets() []string
	DownloadDocset(language string) bool
	LoadDocumentationData(language string, url string) []byte
}

/**
* Concrete implementation of Query engine
 */
type DocsetQueryEngineImpl struct {
}

func (engine DocsetQueryEngineImpl) GetDownloadedDocsets() []string {
	return GetAvailableDocsets()
}

func (engine DocsetQueryEngineImpl) GetDownloadableDocsets() []string {
	res := make([]string, 0)
	for _, v := range GetDocsetFeeds() {
		res = append(res, v.Name)
	}

	return res
}

func (engine DocsetQueryEngineImpl) LoadDocumentationData(language string, url string) []byte {
	return LoadDocumentationUrl(language, url)
}

func (engine DocsetQueryEngineImpl) GetIndicesForLanguage(language string) data_models.Docset {
	return DocsetForLanguage(language)
}

func (engine DocsetQueryEngineImpl) DownloadDocset(language string) bool {
	//TODO: make more efficient by using some map and not a linear search
	feeds := GetDocsetFeeds()

	for _, v := range feeds {
		if v.Name == language {
			DownloadDocset(v.Urls[0])
			return true
		}
	}

	return false
}

var _queryEngine = DocsetQueryEngineImpl{}

func GetQueryEngine() DocsetQueryEngine {
	return _queryEngine
}
