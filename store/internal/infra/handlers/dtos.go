package handlers

import (
	"io"
	"mime/multipart"

	"ichibuy/store/internal/services"
)

type ImageDTO struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type PriceDTO struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"` // cents
	Currency string `json:"currency"`
}

type NewPriceDTO struct {
	Amount   int    `json:"amount"` // cents
	Currency string `json:"currency"`
}

func convertHandlerPriceDTOsToService(handlerDTOs []NewPriceDTO) []services.NewPriceDTO {
	serviceDTOs := make([]services.NewPriceDTO, len(handlerDTOs))
	for i, dto := range handlerDTOs {
		serviceDTOs[i] = services.NewPriceDTO{
			Amount:   dto.Amount,
			Currency: dto.Currency,
		}
	}
	return serviceDTOs
}

func convertMultipartFilesToDTOs(fileHeaders []*multipart.FileHeader) ([]services.FileDTO, error) {
	fileDTOs := make([]services.FileDTO, 0, len(fileHeaders))

	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Read file content
		fileData, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}

		fileDTOs = append(fileDTOs, services.FileDTO{
			FileName:    fileHeader.Filename,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Data:        fileData,
		})
	}

	return fileDTOs, nil
}
