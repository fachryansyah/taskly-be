package dto

type ResponseError struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Message string `json:"message"`
	Target  string `json:"target"`
	Tag     string `json:"tag"`
}

type ResponseWrapper[T any] struct {
	Data       *T                  `json:"data"`
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
	Error      *[]ResponseError    `json:"error,omitempty"`
}
