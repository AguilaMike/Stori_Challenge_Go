package files

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/AguilaMike/Stori_Challenge_Go/internal/common/nats"
	"github.com/google/uuid"
)

type FileUploadService struct {
	natsClient *nats.NatsClient
	uploadDir  string
}

func NewFileUploadService(natsClient *nats.NatsClient, uploadDir string) (*FileUploadService, error) {
	// Asegurarse de que el directorio de carga exista
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return &FileUploadService{
		natsClient: natsClient,
		uploadDir:  uploadDir,
	}, nil
}

func (s *FileUploadService) UploadTransactionFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, userID uuid.UUID) error {
	// Generar un nombre de archivo único
	filename := fmt.Sprintf("%s_%s", userID, filepath.Base(header.Filename))
	filePath := filepath.Join(s.uploadDir, filename)

	// Abrir el archivo con las banderas apropiadas para permitir la sobrescritura
	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error opening file for writing: %w", err)
	}
	defer dst.Close()

	// Copiar el contenido del archivo subido al nuevo archivo
	if _, err = io.Copy(dst, file); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	// Publicar un mensaje en NATS para procesar el archivo
	err = s.natsClient.Publish("transaction.file.uploaded", map[string]interface{}{
		"file_path": filePath,
		"user_id":   userID,
	})
	if err != nil {
		return fmt.Errorf("error publishing NATS message: %w", err)
	}

	return nil
}
