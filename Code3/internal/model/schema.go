package model

import "backend/ehr/internal/utils"

type password struct {
	text string
	hash []byte
}

// define the entity struct model
type Entity struct {
	ID                int64
	FullName          string
	Email             string
	Password          password
	UserType          string
	CreatedAt         string
	IsActive          bool
	EntityInformation EntityInformation
}

type EntityInformation struct {
	EntityId string
	Slug     utils.Slug
}

// define the user struct model
type User struct {
	ID         int64
	FullName   string
	Email      string
	Phone      string
	DOB        string
	Gender     string
	Address    string
	Occupation string
	UserType   string
	IdType     string
	IdNumber   string
	Slug       utils.Slug
	CreatedAt  string
	IsActive   bool
}

// define the patient struct model
type Patient struct {
	ID                 int64
	FullName           string
	Email              string
	DOB                string
	Gender             string
	Phone              string
	Slug               utils.Slug
	Address            string
	Occupation         string
	EmergencyName      string
	EmergencyTelephone string
	CreatedAt          string
}

type Role struct {
	ID          int64
	Name        string
	Description string
	Type        string
}

type Permission struct {
	ID          int64
	Name        string
	Description string
	Type        string
}

type RolePermission struct {
	RoleID       int64
	PermissionID int64
}

type UserRole struct {
	UserID int64
	RoleID int64
}

type UserPermission struct {
	UserID       int64
	PermissionID int64
}
