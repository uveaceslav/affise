package parser

type URLFetcher interface {
	Fetch(uri string) ([]byte, error)
}
