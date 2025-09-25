package domain

import (
	"context"
)

type StorageService interface {
	UploadFiles(ctx context.Context, req []UploadFileRequest) (*UploadFileResponse, error)
	DeleteFiles(ctx context.Context, fileIDs []string) error
}

type UploadFileRequest struct {
	FileName    string
	ContentType string
	Data        []byte
}

type UploadFileResponse struct {
	Infos []UploadedFileInfo
}

type UploadedFileInfo struct {
	ID  string
	URL string
}
