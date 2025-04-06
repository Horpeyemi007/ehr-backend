package dtos

// for response
type PatientDTO struct {
	ID        int64  `json:"id"`
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// for request
type CreatePatientRequest struct {
	FullName string `json:"fullname" validate:"required"`
	DOB      string `json:"dob" validate:"required"`
	Sex      string `json:"sex" validate:"oneof=Male Female"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Street   string `json:"street" validate:"required"`
	City     string `json:"city" validate:"required"`
	State    string `json:"state" validate:"required"`
	Zipcode  string `json:"zipcode" validate:"required"`
}
