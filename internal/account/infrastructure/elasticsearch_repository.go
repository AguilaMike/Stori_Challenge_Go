package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/ports"
)

type ElasticsearchAccountRepository struct {
	client *elastic.Client
	index  string
}

func NewElasticsearchAccountRepository(client *elastic.Client, index string) ports.AccountQueryRepository {
	return &ElasticsearchAccountRepository{
		client: client,
		index:  index,
	}
}

func (r *ElasticsearchAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	result, err := r.client.Get().
		Index(r.index).
		Id(id.String()).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, nil
	}

	var account domain.Account
	err = json.Unmarshal(result.Source, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *ElasticsearchAccountRepository) List(ctx context.Context, limit, offset int64) ([]*domain.Account, error) {
	searchResult, err := r.client.Search().
		Index(r.index).
		From(int(offset)).
		Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var accounts []*domain.Account
	for _, hit := range searchResult.Hits.Hits {
		var account domain.Account
		err := json.Unmarshal(hit.Source, &account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}

func (r *ElasticsearchAccountRepository) Search(ctx context.Context, query string) ([]*domain.Account, error) {
	searchQuery := elastic.NewMultiMatchQuery(query, "nickname").
		Fuzziness("AUTO")

	searchResult, err := r.client.Search().
		Index(r.index).
		Query(searchQuery).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var accounts []*domain.Account
	for _, hit := range searchResult.Hits.Hits {
		var account domain.Account
		err := json.Unmarshal(hit.Source, &account)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &account)
	}
	return accounts, nil
}
