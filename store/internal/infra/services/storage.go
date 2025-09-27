package services

import (
	"context"
	"os"
	"path/filepath"

	fstorage "github.com/Jibaru/ichibuy/api-client/go/fstorage"

	"ichibuy/store/internal/domain"
	sharedCtx "ichibuy/store/internal/shared/context"
)

type storageService struct {
	client *fstorage.APIClient
}

func NewStorageService(client *fstorage.APIClient) domain.StorageService {
	return &storageService{client: client}
}

func (c *storageService) UploadFiles(ctx context.Context, req []domain.UploadFileRequest) (*domain.UploadFileResponse, error) {
	ctx = sharedCtx.AddToken(ctx, fstorage.ContextAccessToken)

	files := make([]*os.File, len(req))
	for i, r := range req {
		f, err := fileRequestToFile(&r)
		if err != nil {
			return nil, err
		}
		files[i] = f
		defer f.Close()
		defer os.Remove(f.Name())
	}

	resp, _, err := c.client.FilesApi.ApiV1FilesUploadPost(ctx, files, "store")
	if err != nil {
		return nil, err
	}

	infos := make([]domain.UploadedFileInfo, len(resp.Files))
	for i, f := range resp.Files {
		infos[i] = domain.UploadedFileInfo{
			ID:  f.Id,
			URL: f.Url,
		}
	}

	return &domain.UploadFileResponse{Infos: infos}, nil
}

func (c *storageService) DeleteFiles(ctx context.Context, fileIDs []string) error {
	ctx = sharedCtx.AddToken(ctx, fstorage.ContextAPIKey)
	_, err := c.client.FilesApi.ApiV1FilesBatchDelete(ctx, fstorage.DeleteRequest{FileIds: fileIDs})
	if err != nil {
		return err
	}
	return nil
}

func fileRequestToFile(u *domain.UploadFileRequest) (*os.File, error) {
	ext := filepath.Ext(u.FileName)
	f, err := os.CreateTemp("", "upload_*"+ext)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(u.Data)
	if err != nil {
		f.Close()
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		f.Close()
		return nil, err
	}

	return f, nil
}
