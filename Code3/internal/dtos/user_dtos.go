package dtos

// for response
type UserDTO struct {
	ID        int64    `json:"id"`
	FullName  string   `json:"fullname"`
	Email     string   `json:"email"`
	Slug      string   `json:"slug"`
	Role      []string `json:"role"`
	CreatedAt string   `json:"created_at"`
}

// for payload request
type CreateUserRequest struct {
	FullName   string   `json:"fullname" validate:"required"`
	Email      string   `json:"email" validate:"required,email"`
	Phone      string   `json:"phone" validate:"required"`
	DOB        string   `json:"dob" validate:"required"`
	Gender     string   `json:"gender" validate:"required"`
	Address    string   `json:"address" validate:"required"`
	Occupation string   `json:"occupation" validate:"required"`
	Role       []string `json:"role" validate:"max=5"`
}
