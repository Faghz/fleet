package request

// Organization requests
type CreateOrganization struct {
	Name     string `json:"name" validate:"required,min=1,max=255" example:"PT Sample Company"`
	Currency string `json:"currency" validate:"required,len=3" example:"IDR"`
}

type UpdateOrganization struct {
	Name     string `json:"name" validate:"required,min=1,max=255" example:"PT Sample Company"`
	Currency string `json:"currency" validate:"required,len=3" example:"IDR"`
}

// Account requests
type CreateAccount struct {
	Code     string  `json:"code" validate:"required,min=1,max=50" example:"1110"`
	Name     string  `json:"name" validate:"required,min=1,max=255" example:"Cash"`
	Type     string  `json:"type" validate:"required,oneof=asset liability equity revenue expense" example:"asset"`
	ParentID *string `json:"parent_id" example:"acc_1000"`
	IsActive bool    `json:"is_active" example:"true"`
}

type UpdateAccount struct {
	Code     string  `json:"code" validate:"required,min=1,max=50" example:"1110"`
	Name     string  `json:"name" validate:"required,min=1,max=255" example:"Cash"`
	Type     string  `json:"type" validate:"required,oneof=asset liability equity revenue expense" example:"asset"`
	ParentID *string `json:"parent_id" example:"acc_1000"`
	IsActive bool    `json:"is_active" example:"true"`
}

// Journal requests
type JournalLine struct {
	AccountID   string  `json:"account_id" validate:"required" example:"acc_1110"`
	Debit       float64 `json:"debit" validate:"min=0" example:"100000"`
	Credit      float64 `json:"credit" validate:"min=0" example:"0"`
	Description *string `json:"description" example:"Initial payment"`
}

type CreateJournalEntry struct {
	Date          string        `json:"date" validate:"required" example:"2025-01-15"`
	Description   string        `json:"description" validate:"required,min=1" example:"Service revenue"`
	Reference     *string       `json:"reference" example:"INV-001"`
	AttachmentURL *string       `json:"attachment_url" example:"https://example.com/receipt.pdf"`
	CategoryID    *string       `json:"category_id" example:"cat_sales"`
	Lines         []JournalLine `json:"lines" validate:"required,min=2,dive"`
}

type UpdateJournalEntry struct {
	Date          string        `json:"date" validate:"required" example:"2025-01-15"`
	Description   string        `json:"description" validate:"required,min=1" example:"Service revenue"`
	Reference     *string       `json:"reference" example:"INV-001"`
	AttachmentURL *string       `json:"attachment_url" example:"https://example.com/receipt.pdf"`
	CategoryID    *string       `json:"category_id" example:"cat_sales"`
	Status        string        `json:"status" validate:"oneof=draft posted cancelled" example:"posted"`
	Lines         []JournalLine `json:"lines" validate:"required,min=2,dive"`
}

// Journal Category requests
type CreateJournalCategory struct {
	Name        string  `json:"name" validate:"required,min=1,max=255" example:"Sales"`
	Description *string `json:"description" example:"Sales and revenue transactions"`
	Color       *string `json:"color" validate:"omitempty,hexcolor" example:"#10B981"`
}

type UpdateJournalCategory struct {
	Name        string  `json:"name" validate:"required,min=1,max=255" example:"Sales"`
	Description *string `json:"description" example:"Sales and revenue transactions"`
	Color       *string `json:"color" validate:"omitempty,hexcolor" example:"#10B981"`
	IsActive    bool    `json:"is_active" example:"true"`
}

type GetJournalCategory struct {
	ID string `param:"id" validate:"required" example:"cat_12345"`
}

type DeleteJournalCategory struct {
	ID string `param:"id" validate:"required" example:"cat_12345"`
}

type SearchJournalCategories struct {
	Query  string `query:"q" example:"sales"`
	Limit  int    `query:"limit" validate:"min=1,max=100" example:"10"`
	Offset int    `query:"offset" validate:"min=0" example:"0"`
}

// Signup with organization
type SignupWithOrganization struct {
	Email       string `json:"email" validate:"required,email" example:"admin@example.com"`
	Password    string `json:"password" validate:"required,min=8" example:"SecurePass123!"`
	Name        string `json:"name" validate:"required,min=1,max=255" example:"John Doe"`
	OrgName     string `json:"org_name" validate:"required,min=1,max=255" example:"PT Sample Company"`
	OrgCurrency string `json:"org_currency" validate:"required,len=3" example:"IDR"`
}
