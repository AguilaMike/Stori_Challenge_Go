package files

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/google/uuid"
)

type FileUploadService struct {
	natsClient *nats.NatsClient
}

func NewFileUploadService(natsClient *nats.NatsClient) (*FileUploadService, error) {
	return &FileUploadService{
		natsClient: natsClient,
	}, nil
}

func (s *FileUploadService) UploadTransactionFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID) error {
	// Leer el contenido del archivo
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file content: %w", err)
	}

	// Codificar el contenido en base64
	encodedContent := base64.StdEncoding.EncodeToString(content)

	// Publicar un mensaje en NATS con el contenido del archivo
	err = s.natsClient.Publish("transaction.file.uploaded", map[string]interface{}{
		"file_name":    header.Filename,
		"file_content": encodedContent,
		"user_id":      userID,
	})
	if err != nil {
		return fmt.Errorf("error publishing NATS message: %w", err)
	}

	return nil
}
