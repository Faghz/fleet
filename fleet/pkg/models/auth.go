package models

type AuthClaims struct {
	ID      string `json:"id"`
	Subject string `json:"sub"`
	OrgID   string `json:"org_id"`
}

type InsertAuthParams struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Password  string `json:"password"`
	CreatedBy string `json:"created_by"`
}
