package models

type Sort struct {
	Field string `json:"field" query:"field"`
	Order string `json:"order" query:"order"`
}

type Pagination struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

type Metadata struct {
	PageSize  int `json:"page_size"`
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

// Fungsi untuk menerapkan pagination
func (p *Pagination) ApplyPagination() {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize < 1 {
		p.PageSize = 10
	}
}

func (p *Sort) ApplySorting() {
	if p.Field == "" {
		p.Field = "created_at"
	}
	if p.Order != "asc" {
		p.Order = "desc"
	}
}
