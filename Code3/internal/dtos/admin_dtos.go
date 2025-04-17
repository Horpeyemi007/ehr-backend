package dtos

// for response data
type AdminDTO struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Slug     string `json:"slug"`
	Token    string `json:"token"`
}

// for payload request
type CreateAdminRequest struct {
	FullName string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
