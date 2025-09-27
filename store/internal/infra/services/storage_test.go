package services_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	swagger "github.com/Jibaru/ichibuy/api-client/go/fstorage"

	"ichibuy/store/internal/domain"
	"ichibuy/store/internal/infra/services"
	sharedCtx "ichibuy/store/internal/shared/context"
)

func Test_fileUploadClient_UploadFile(t *testing.T) {
	sut := services.NewStorageService(swagger.NewAPIClient(&swagger.Configuration{
		BasePath:   "https://ichibuy-fstorage.vercel.app",
		HTTPClient: &http.Client{},
	}))

	ctx := context.WithValue(context.Background(), sharedCtx.APITokenKey, os.Getenv("FSTORAGE_API_TOKEN"))

	resp, err := sut.UploadFiles(ctx, []domain.UploadFileRequest{{
		FileName:    "test.txt",
		ContentType: "text/plain",
		Data:        []byte("Hello, World!"),
	}})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}
