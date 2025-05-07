package dtos

// for response data
type EntityDTO struct {
	ID       int64  `json:"id"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Slug     string `json:"slug"`
	Token    string `json:"token"`
}

type RoleDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type PermissionDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

// for payload request
type CreateEntityRequest struct {
	FullName string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type EntityLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type RolePermissionAssign struct {
	RoleName    string   `json:"roleName" validate:"required"`
	Permissions []string `json:"permissionName" validate:"required"`
}

type RoleUserAssign struct {
	RoleName string `json:"roleName" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
}

type PermissionUserAssign struct {
	Permissions []string `json:"permissionName" validate:"required"`
	Slug        string   `json:"slug" validate:"required"`
}

type DeleteUserPermission struct {
	Permission string `json:"permissionName" validate:"required"`
	Slug       string `json:"slug" validate:"required"`
}

type DeleteRolePermission struct {
	Permission string `json:"permissionName" validate:"required"`
	RoleName   string `json:"roleName" validate:"required"`
}
