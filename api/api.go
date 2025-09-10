package api

import "github.com/jexroid/gopi/pkg/models"

// Userdetail represents user information
// @Description User details structure
type Userdetail struct {
	// @Description User's first name
	// @Example John
	Firstname string

	// @Description User's last name
	// @Example Doe
	Lastname string

	// @Description User's phone number
	// @Example 1234567890
	Phone int
}

type OTPRequest struct {
	Phone int `json:"phone"`
}

type OTPVerifyRequest struct {
	Phone int    `json:"phone"`
	Code  string `json:"code" validate:"required,numeric,len=6"`
}

type LoginRequest struct {
	Phone    int    `json:"phone"`
	Password string `json:"password"`
}

// ValidateRequest represents the token validation request
// @Description Token validation request payload
type ValidateRequest struct {
	Token string `json:"token"`
}

type Argonparams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type RegisterResponse struct {
	Ok        bool
	UserExist bool
	Message   string
	Token     string `json:"Token,omitempty"`
}

type LoginResponse struct {
	Ok      bool
	Valid   bool
	Token   string
	Message string
}

// ValidateResponse represents the token validation response
// @Description Token validation response with user details
type ValidateResponse struct {
	Ok      bool
	Payload string
	Exp     bool
	User    Userdetail
}

type OTPVerifyResponse struct {
	Ok    bool   `json:"ok"`
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

type OTPResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

// UsersResponse represents the response structure for paginated users
// @Description Paginated users response with metadata
type UsersResponse struct {
	Users      []models.User  `json:"users"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination metadata
// @Description Pagination information
type PaginationInfo struct {
	Page       int  `json:"page" example:"1"`
	Limit      int  `json:"limit" example:"10"`
	Total      int  `json:"total" example:"100"`
	TotalPages int  `json:"total_pages" example:"10"`
	HasNext    bool `json:"has_next" example:"true"`
	HasPrev    bool `json:"has_prev" example:"false"`
}
