package services

type CreateUpdateResponse struct {
	ID string `json:"id"`
}

// Request pagination filters

type Pagination struct {
	Limit  int
	Offset int
}

type Sorting struct {
	Field string
	Order string
}

type ImageDTO struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type NewPriceDTO struct {
	Amount   int    `json:"amount"` // cents
	Currency string `json:"currency"`
}

type PriceDTO struct {
	ID       string `json:"id"`
	Amount   int    `json:"amount"` // cents
	Currency string `json:"currency"`
}

type FileDTO struct {
	FileName    string
	ContentType string
	Data        []byte
}
