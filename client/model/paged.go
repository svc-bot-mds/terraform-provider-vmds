package model

type Paged[T any] struct {
	Embedded map[string][]T `json:"_embedded"`
	Page     PageInfo       `json:"page"`
}

type PageInfo struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

type PageQuery struct {
	Index int `schema:"page"`
	Size  int `schema:"size"`
}

func (p *Paged[T]) Get() *[]T {
	empty := make([]T, 0)
	var items = &empty
	for _, v := range p.Embedded {
		items = &v
	}
	return items
}

func (p *Paged[T]) GetPage() *PageInfo {
	return &p.Page
}
