package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/account/ports"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
)

type ElasticsearchAccountRepository struct {
	client *elastic.Client
	index  string
	nats   *nats.NatsClient
}

func NewElasticsearchAccountRepository(client *elastic.Client, nc *nats.NatsClient, index string) ports.AccountQueryRepository {
	repo := &ElasticsearchAccountRepository{
		client: client,
		index:  index,
		nats:   nc,
	}
	repo.subscribeToEvents()
	return repo
}

func (r *ElasticsearchAccountRepository) subscribeToEvents() {
	r.nats.Subscribe("account.created", r.handleAccountCreated)
	r.nats.Subscribe("account.updated", r.handleAccountUpdated)
	r.nats.Subscribe("account.deleted", r.handleAccountDeleted)
}

func (r *ElasticsearchAccountRepository) handleAccountCreated(data []byte) {
	var account domain.Account
	if err := json.Unmarshal(data, &account); err != nil {
		log.Printf("Error unmarshaling account: %v", err)
		return
	}

	_, err := r.client.Index().
		Index(r.index).
		Id(account.ID.String()).
		BodyJson(account).
		Do(context.Background())
	if err != nil {
		log.Printf("Error indexing account: %v", err)
	}
}

func (r *ElasticsearchAccountRepository) handleAccountUpdated(data []byte) {
	var account domain.Account
	if err := json.Unmarshal(data, &account); err != nil {
		log.Printf("Error unmarshaling account: %v", err)
		return
	}

	_, err := r.client.Update().
		Index(r.index).
		Id(account.ID.String()).
		Doc(account).
		Do(context.Background())
	if err != nil {
		log.Printf("Error updating account: %v", err)
	}
}

func (r *ElasticsearchAccountRepository) handleAccountDeleted(data []byte) {
	var payload struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Printf("Error unmarshaling account ID: %v", err)
		return
	}

	_, err := r.client.Delete().
		Index(r.index).
		Id(payload.ID).
		Do(context.Background())
	if err != nil {
		log.Printf("Error deleting account: %v", err)
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
