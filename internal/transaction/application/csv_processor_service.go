package application

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"

	csv "github.com/AguilaMike/Stori_Challenge_Go/internal/common/files/csv_file"
)

type CSVProcessorService struct {
	transactionService *TransactionService
	accountID          uuid.UUID
	filepath           string
}

func NewCSVProcessorService(transactionService *TransactionService) *CSVProcessorService {
	return &CSVProcessorService{
		transactionService: transactionService,
	}
}

func (s *CSVProcessorService) ProcessCSVFile(filePath string, accountID uuid.UUID) error {
	s.accountID = accountID
	s.filepath = filePath

	err := csv.ProcessCSVFileInRow(filePath, s.processTransaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *CSVProcessorService) processTransaction(data []string) error {
	amount, err := strconv.ParseFloat(data[1], 64)
	if err != nil {
		return err
	}

	date, err := time.Parse("2006-01-02", data[0])
	if err != nil {
		return err
	}

	_, err = s.transactionService.CreateTransaction(context.Background(), s.accountID, amount, s.filepath, date)
	if err != nil {
		return err
	}
	return nil
}
