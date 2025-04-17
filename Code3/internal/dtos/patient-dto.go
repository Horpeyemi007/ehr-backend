package dtos

// for response data
type PatientDTO struct {
	Slug string `json:"slug"`
}

// for response data
type GetAllPatientResponse struct {
	Slug               string `json:"slug"`
	FullName           string `json:"fullname"`
	Email              string `json:"email"`
	Gender             string `json:"gender"`
	Phone              string `json:"phone"`
	DOB                string `json:"dob"`
	Address            string `json:"address"`
	Occupation         string `json:"occupation"`
	EmergencyName      string `json:"emergencyName"`
	EmergencyTelephone string `json:"emergencyPhone"`
}

// for payload request
type CreatePatientRequest struct {
	FullName           string `json:"fullname" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Gender             string `json:"gender" validate:"oneof=male female"`
	Phone              string `json:"phone" validate:"required"`
	DOB                string `json:"dob" validate:"required"`
	Address            string `json:"address" validate:"required"`
	Occupation         string `json:"occupation" validate:"required"`
	EmergencyName      string `json:"emergencyName" validate:"required"`
	EmergencyTelephone string `json:"emergencyPhone" validate:"required"`
}
