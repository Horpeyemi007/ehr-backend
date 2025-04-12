package dtos

// for response
type PatientDTO struct {
	ID        int64  `json:"id"`
	FullName  string `json:"fullname"`
	Email     string `json:"email"`
	Slug      string `json:"slug"`
	CreatedAt string `json:"created_at"`
}

// for payload request
type CreatePatientRequest struct {
	FullName          string `json:"fullname" validate:"required"`
	DOB               string `json:"dob" validate:"required"`
	Gender            string `json:"gender" validate:"oneof=Male Female"`
	Phone             string `json:"phone" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	Address           string `json:"address" validate:"required"`
	Occupation        string `json:"occupation" validate:"required"`
	GuardianName      string `json:"guardian_name" validate:"required"`
	GuardianTelephone string `json:"guardian_telephone" validate:"required"`
}
