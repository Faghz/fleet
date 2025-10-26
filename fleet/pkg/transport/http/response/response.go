package response

import (
	"encoding/csv"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ErrorCode int

type BaseResponse struct {
	Status  int         `json:"status" example:"200"`
	Message string      `json:"message,omitempty" example:"A Poem Here"`
	Data    interface{} `json:"data,omitempty"`
}

// Metadata for pagination responses
type Metadata struct {
	Page        int `json:"page" example:"1"`
	Limit       int `json:"limit" example:"50"`
	TotalPages  int `json:"total_pages" example:"10"`
	TotalData   int `json:"total_data" example:"500"`
	CurrentPage int `json:"current_page" example:"1"`
}

// PaginatedResponse is a generic wrapper for paginated responses
type PaginatedResponse struct {
	Items    interface{} `json:"items"`
	Metadata Metadata    `json:"metadata"`
}

// CreateMetadata creates pagination metadata
func CreateMetadata(page, limit, totalData int) Metadata {
	totalPages := 0
	if totalData > 0 {
		totalPages = (totalData + limit - 1) / limit // Ceiling division
	}

	return Metadata{
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		TotalData:   totalData,
		CurrentPage: page,
	}
}

// CreatePaginatedResponse creates a paginated response with metadata
func CreatePaginatedResponse(items interface{}, page, limit, totalData int) PaginatedResponse {
	return PaginatedResponse{
		Items:    items,
		Metadata: CreateMetadata(page, limit, totalData),
	}
}

// Organization responses
type Organization struct {
	ID        string    `json:"id" example:"org_001"`
	Name      string    `json:"name" example:"PT Sample Company"`
	Currency  string    `json:"currency" example:"IDR"`
	CreatedAt time.Time `json:"created_at"`
}

// Account responses
type Account struct {
	ID        string     `json:"id" example:"acc_1110"`
	Code      string     `json:"code" example:"1110"`
	Name      string     `json:"name" example:"Cash"`
	Type      string     `json:"type" example:"asset"`
	ParentID  *string    `json:"parent_id" example:"acc_1000"`
	IsActive  bool       `json:"is_active" example:"true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type AccountBalance struct {
	Account     Account `json:"account"`
	TotalDebit  float64 `json:"total_debit" example:"150000"`
	TotalCredit float64 `json:"total_credit" example:"50000"`
	Balance     float64 `json:"balance" example:"100000"`
}

// Standard list wrapper with metadata for accounts
type AccountList struct {
	Items    []Account `json:"items"`
	Metadata Metadata  `json:"metadata"`
}

// Journal responses
type JournalCategory struct {
	ID          string     `json:"id" example:"cat_sales"`
	Name        string     `json:"name" example:"Sales"`
	Description *string    `json:"description" example:"Sales and revenue transactions"`
	Color       string     `json:"color" example:"#10B981"`
	IsActive    bool       `json:"is_active" example:"true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type JournalCategoryList struct {
	Categories []JournalCategory `json:"categories"`
	Total      int               `json:"total"`
}

type User struct {
	ID   string `json:"id" example:"user_001"`
	Name string `json:"name" example:"John Doe"`
}

type JournalLine struct {
	ID          string  `json:"id" example:"jl_001"`
	AccountID   string  `json:"account_id" example:"acc_1110"`
	Account     Account `json:"account"`
	Debit       float64 `json:"debit" example:"100000"`
	Credit      float64 `json:"credit" example:"0"`
	Description *string `json:"description" example:"Initial payment"`
}

type JournalEntry struct {
	ID            string           `json:"id" example:"je_001"`
	Date          string           `json:"date" example:"2025-01-15"`
	Description   string           `json:"description" example:"Service revenue"`
	Reference     *string          `json:"reference" example:"INV-001"`
	AttachmentURL *string          `json:"attachment_url" example:"https://example.com/receipt.pdf"`
	Category      *JournalCategory `json:"category,omitempty"`
	Status        string           `json:"status" example:"posted"`
	Lines         []JournalLine    `json:"lines"`
	CreatedBy     *User            `json:"created_by,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty"`
}

// Standard list wrapper with metadata for journal entries
type JournalEntryList struct {
	Items    []JournalEntry `json:"items"`
	Metadata Metadata       `json:"metadata"`
}

// Report responses
type BalanceSheetItem struct {
	Account  Account            `json:"account"`
	Balance  float64            `json:"balance" example:"100000"`
	Children []BalanceSheetItem `json:"children,omitempty"`
}

type BalanceSheet struct {
	AsOfDate         string             `json:"as_of_date" example:"2025-01-31"`
	Assets           []BalanceSheetItem `json:"assets"`
	Liabilities      []BalanceSheetItem `json:"liabilities"`
	Equity           []BalanceSheetItem `json:"equity"`
	TotalAssets      float64            `json:"total_assets" example:"1000000"`
	TotalLiabilities float64            `json:"total_liabilities" example:"300000"`
	TotalEquity      float64            `json:"total_equity" example:"700000"`
}

type IncomeStatementItem struct {
	Account  Account               `json:"account"`
	Amount   float64               `json:"amount" example:"500000"`
	Children []IncomeStatementItem `json:"children,omitempty"`
}

type IncomeStatement struct {
	StartDate     string                `json:"start_date" example:"2025-01-01"`
	EndDate       string                `json:"end_date" example:"2025-01-31"`
	Revenue       []IncomeStatementItem `json:"revenue"`
	Expenses      []IncomeStatementItem `json:"expenses"`
	TotalRevenue  float64               `json:"total_revenue" example:"1500000"`
	TotalExpenses float64               `json:"total_expenses" example:"1000000"`
	NetIncome     float64               `json:"net_income" example:"500000"`
}

type NetIncomeResponse struct {
	StartDate string  `json:"start_date" example:"2025-01-01"`
	EndDate   string  `json:"end_date" example:"2025-01-31"`
	NetIncome float64 `json:"net_income" example:"500000"`
}

// Laporan Laba Rugi (Income Statement) Response with Indonesian terminology
type LabaRugi struct {
	StartDate       string  `json:"start_date" example:"2025-01-01"`    // Tanggal Mulai
	EndDate         string  `json:"end_date" example:"2025-01-31"`      // Tanggal Akhir
	TotalPendapatan float64 `json:"total_pendapatan" example:"1500000"` // Total Pendapatan (Revenue)
	TotalBeban      float64 `json:"total_beban" example:"1000000"`      // Total Beban (Expenses)
	LabaRugiBersih  float64 `json:"laba_rugi_bersih" example:"500000"`  // Laba/Rugi Bersih (Net Income)
}

// Balance Sheet (Neraca) Response with English terminology - Summary only
type Neraca struct {
	AsOfDate         string  `json:"as_of_date" example:"2025-01-31"`    // As of date
	TotalAssets      float64 `json:"total_assets" example:"1000000"`     // Total Assets
	TotalLiabilities float64 `json:"total_liabilities" example:"300000"` // Total Liabilities
	TotalEquity      float64 `json:"total_equity" example:"700000"`      // Total Equity
	IsBalanced       bool    `json:"is_balanced" example:"true"`         // Is balanced
}

type Failure struct {
	// http status code
	Code int `json:"-"`
	//  service error code
	ErrorCode ErrorCode `json:"errorCode,omitempty" example:"1001001"`
	// error title
	Title string `json:"message" example:"Not Enough Swipe Token"`
	// error detail
	Description string `json:"description" example:"You run out of swipe token to find you next life partner"`
	// error field if exists
	Errors []FailureError `json:"errors,omitempty"`
}

func (e *Failure) Error() string {
	return e.Description
}

type FailureError struct {
	// This is pointer for guiding which field caused the error, or could be a subject like wallet/points
	Pointer string `json:"pointer" example:"#message[0].content"`

	// The message why the field has an error
	Message string `json:"message" example:"Is a required field, *Note that Errors fields only available on bad request"`
}

func GenerateFailure(httpStatus int, title, description string, errors ...FailureError) error {
	failure := Failure{
		Code:        httpStatus,
		Title:       title,
		Description: description,
	}

	if len(errors) > 0 {
		failure.Errors = append(failure.Errors, errors...)
	}

	return fiber.NewError(httpStatus, description)
}

func GenerateBadRequestFromFiberError(err error) error {
	fiberError, ok := err.(*fiber.Error)
	if !ok {
		return GenerateBadRequest("Bad Request", err.Error())
	}
	return GenerateBadRequest("Bad Request", fiberError.Message)
}

func GenerateBadRequest(title, description string, errors ...FailureError) error {
	return &Failure{
		Code:        http.StatusBadRequest,
		Title:       title,
		Description: description,
		Errors:      errors,
	}
}

// FailureResponse sends a Failure struct as a JSON response with the appropriate HTTP status code
func FailureResponse(c *fiber.Ctx, failure error) error {
	failureDetail, ok := failure.(*Failure)
	if !ok {
		failureDetail = GenerateFailure(http.StatusInternalServerError, "Internal Server Error", failure.Error()).(*Failure)
	}

	return c.Status(failureDetail.Code).JSON(failureDetail)
}

func ErrorBadRequest(message string) error {
	return &Failure{
		Code:        http.StatusBadRequest,
		Title:       "Bad Request",
		Description: message,
	}
}

// ResponseJson will return a standardized JSON response
func ResponseJson(ctx *fiber.Ctx, status int, message string, data ...interface{}) error {
	response := BaseResponse{
		Status:  status,
		Message: message,
	}

	if len(data) > 0 {
		response.Data = data[0]
	}

	return ctx.Status(status).JSON(response)
}

func CSVResponse(ctx *fiber.Ctx, filename string, data [][]string) (err error) {
	writer := csv.NewWriter(ctx.Response().BodyWriter())

	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return GenerateBadRequest("Failed to write CSV", err.Error())
		}
	}

	// Flush the data to the response
	writer.Flush()

	if err := writer.Error(); err != nil {
		return GenerateBadRequest("Failed to flush CSV data", err.Error())
	}

	ctx.Set(fiber.HeaderContentType, "text/csv")
	ctx.Set(fiber.HeaderContentDisposition, "attachment; filename="+filename+".csv")

	return
}
