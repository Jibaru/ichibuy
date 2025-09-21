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
