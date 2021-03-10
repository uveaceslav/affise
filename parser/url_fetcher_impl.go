package parser

import (
	"io/ioutil"
	"net/http"
)

type urlFetcher struct {
	httpClient *http.Client
}

func NewURLFetcher(httpClient *http.Client) URLFetcher {
	return &urlFetcher{
		httpClient: httpClient,
	}
}

func (uf *urlFetcher) Fetch(uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true

	resp, err := uf.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
