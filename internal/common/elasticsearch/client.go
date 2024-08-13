package elasticsearch

import (
	"github.com/olivere/elastic/v7"
)

func NewElasticsearchClient(url string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
}
