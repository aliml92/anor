package typesense

import "github.com/aliml92/go-typesense/typesense"

type Searcher struct {
	client *typesense.Client
}

func NewSearcher(client *typesense.Client) *Searcher {
	return &Searcher{
		client: client,
	}
}
