package common

type PagedResponse[T any, U any] struct {
	Items    []*T `json:"items"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Extra    U    `json:"extra,omitempty"`
}
