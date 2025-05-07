package dtos

// for response data
type UserDTO struct {
	ID        int64  `json:"id"`
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	Slug      string `json:"slug"`
	UserType  string `json:"userType"`
	CreatedAt string `json:"created_at"`
}

// for payload request
type CreateUserRequest struct {
	FullName   string `json:"fullname" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	DOB        string `json:"dob" validate:"required"`
	Gender     string `json:"gender" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	UserType   string `json:"userType" validate:"required"`
	IdType     string `json:"idType"`
	IdNumber   string `json:"idNumber"`
}

type CreateAdminUserRequest struct {
	FullName   string `json:"fullname" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required"`
	DOB        string `json:"dob" validate:"required"`
	Gender     string `json:"gender" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Occupation string `json:"occupation" validate:"required"`
	UserType   string `json:"userType" validate:"required,oneof=admin"`
	IdType     string `json:"idType"`
	IdNumber   string `json:"idNumber"`
}
