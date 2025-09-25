package services_test

import (
	"context"
	"net/http"
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

	ctx := context.WithValue(context.Background(), sharedCtx.APITokenKey, "eyJhbGciOiJSUzI1NiIsImtpZCI6IjNGYkJHanFib2FnPSIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Inh4bmFjaG8xOTk3eHhAZ21haWwuY29tIiwiZXhwIjoxNzU4ODQxNjMyLCJpYXQiOjE3NTg3NTUyMzIsImlzcyI6ImljaGlidXktYXV0aCIsInVzZXJfaWQiOiJmOTQyYTU2NC04YjA1LTQwNjgtOGM2MC01NzZiM2YxYWFmYTkifQ.S8jK_J58vLWoVbzticIiovU2YmgPY_jEN-kygurcI7voMWL-wHSQc_NR9AK5cQngNHlNw7QmnRzDYdHmt3wCfmrPzJgxHEF7Oa1_mU5NE1NJ3FNKL52VFBHmhQNXzmeaNzjn7IjJniV7C2Wfwz1nZAxdJ0ZzEX-Xt3mI0Rvh9kOgPSO_uKqLDeu38GtLFZ0fEAi00VMDN6Wz1BCxwSE6C1g22rMHZwfWJvGlZwEkqMGYH9TA_6l7UwmxSCSyX12ds9Lu57HEQGHrxupPpa9-QvUxEpR49yPW7vgs5qRvYytePtLyiWJDlCquAZg2tKLRLmqBwvmXzXsFqoe20TZcWA")

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
