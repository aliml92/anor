package typesense

import "github.com/aliml92/go-typesense/typesense"

type TsSearcher struct {
	tsClient *typesense.Client
}

func NewSearcher(tsClient *typesense.Client) *TsSearcher {
	return &TsSearcher{
		tsClient: tsClient,
	}
}
