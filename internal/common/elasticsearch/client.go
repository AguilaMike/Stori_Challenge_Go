package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/olivere/elastic/v7"
)

func NewElasticsearchClient(url string) (*elastic.Client, error) {
	esClient, err := connectToElasticsearch(url)
	if err != nil {
		log.Fatalf("Failed to connect to Elasticsearch: %v", err)
	}

	// Configurar índices
	if err := SetupElasticsearchIndex(esClient, "accounts"); err != nil {
		log.Fatalf("Failed to setup Elasticsearch index 'accounts': %v", err)
	}
	if err := SetupElasticsearchIndex(esClient, "transactions"); err != nil {
		log.Fatalf("Failed to setup Elasticsearch index 'transactions': %v", err)
	}

	return esClient, err
}

func connectToElasticsearch(url string) (*elastic.Client, error) {
	var esClient *elastic.Client
	var err error

	for i := 0; i < 5; i++ {
		esClient, err = elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(100*time.Millisecond, 5*time.Second))),
		)
		if err == nil {
			return esClient, nil
		}
		fmt.Printf("Failed to connect to Elasticsearch, retrying in 5 seconds... (attempt %d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to Elasticsearch after 5 attempts: %v", err)
}

func SetupElasticsearchIndex(client *elastic.Client, indexName string) error {
	ctx := context.Background()

	// Verifica si el índice existe
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}

	// Si el índice no existe, créalo
	if !exists {
		createIndex, err := client.CreateIndex(indexName).Do(ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			return errors.New("elasticsearch did not acknowledge index creation")
		}
		log.Printf("Created index: %s", indexName)
	} else {
		log.Printf("Index already exists: %s", indexName)
	}

	return nil
}
