package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/olivere/elastic/v7"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/domain"
	"github.com/AguilaMike/Stori_Challenge_Go/internal/transaction/ports"
)

type ElasticsearchTransactionRepository struct {
	client *elastic.Client
	index  string
	nats   *nats.NatsClient
}

func NewElasticsearchTransactionRepository(client *elastic.Client, nc *nats.NatsClient, index string) ports.TransactionQueryRepository {
	repo := &ElasticsearchTransactionRepository{
		client: client,
		index:  index,
		nats:   nc,
	}
	repo.subscribeToEvents()
	return repo
}

func (r *ElasticsearchTransactionRepository) subscribeToEvents() {
	r.nats.Subscribe("transaction.created", r.handleTransactionCreated)
}

func (r *ElasticsearchTransactionRepository) handleTransactionCreated(data []byte) {
	var transaction domain.Transaction
	if err := json.Unmarshal(data, &transaction); err != nil {
		log.Printf("Error unmarshaling transaction: %v", err)
		return
	}

	_, err := r.client.Index().
		Index(r.index).
		Id(transaction.ID.String()).
		BodyJson(transaction).
		Do(context.Background())
	if err != nil {
		log.Printf("Error indexing transaction: %v", err)
	}
}

func (r *ElasticsearchTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
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

	var transaction domain.Transaction
	err = json.Unmarshal(result.Source, &transaction)
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *ElasticsearchTransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int64) ([]*domain.Transaction, error) {
	query := elastic.NewTermQuery("account_id", accountID.String())
	searchResult, err := r.client.Search().
		Index(r.index).
		Query(query).
		//From(int(offset)).
		//Size(int(limit)).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var transactions []*domain.Transaction
	for _, hit := range searchResult.Hits.Hits {
		var transaction domain.Transaction
		err := json.Unmarshal(hit.Source, &transaction)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func (r *ElasticsearchTransactionRepository) GetSummary(ctx context.Context, accountID uuid.UUID) (*domain.TransactionSummary, error) {
	query := elastic.NewTermQuery("account_id", accountID.String())

	aggs := elastic.NewTermsAggregation().Field("type")
	aggs.SubAggregation("total", elastic.NewSumAggregation().Field("amount"))
	aggs.SubAggregation("avg", elastic.NewAvgAggregation().Field("amount"))
	aggs.SubAggregation("count", elastic.NewValueCountAggregation().Field("_id"))
	aggs.SubAggregation("by_year_month", elastic.NewDateHistogramAggregation().Field("input_date").CalendarInterval("month").Format("yyyy-MM"))

	searchResult, err := r.client.Search().
		Index(r.index).
		Query(query).
		Aggregation("transactions", aggs).
		Size(0).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	summary := &domain.TransactionSummary{
		Monthly: make(map[string]*domain.TransactionMonthly),
	}

	if agg, found := searchResult.Aggregations.Terms("transactions"); found {
		for _, bucket := range agg.Buckets {
			transactionType := bucket.Key.(string)

			if totalAgg, found := bucket.Sum("total"); found {
				if transactionType == "credit" {
					summary.TotalBalance += *totalAgg.Value
				} else if transactionType == "debit" {
					summary.TotalBalance -= *totalAgg.Value
				}
			}

			if avgAgg, found := bucket.Avg("avg"); found {
				if transactionType == "credit" {
					summary.AverageCredit = *avgAgg.Value
				} else if transactionType == "debit" {
					summary.AverageDebit = *avgAgg.Value
				}
			}

			if countAgg, found := bucket.ValueCount("count"); found {
				summary.TotalCount += int(*countAgg.Value)
			}
		}
	}

	return summary, nil
}
