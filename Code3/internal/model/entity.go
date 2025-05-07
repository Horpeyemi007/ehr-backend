package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// define a db instance patient struct
type entityStore struct {
	db *sql.DB
}

var (
	ErrEntityDuplicateEmail = errors.New("an entity with the given email already exists")
)

func (e *entityStore) Create(ctx context.Context, entity *Entity) error {
	// check if the entity already exist in the record
	isFound, _ := e.FindByEmail(ctx, entity.Email)
	if isFound != nil {
		return ErrEntityDuplicateEmail
	}

	return withTransaction(ctx, e.db, func(tx *sql.Tx) error {
		// insert to the entity table
		query := `INSERT INTO entity (fullname, email, password) 
		VALUES ($1, $2, $3) RETURNING id`

		ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
		defer cancel()

		var entityId int64
		err := tx.QueryRowContext(ctx, query, entity.FullName, entity.Email, entity.Password.hash).Scan(&entityId)
		if err != nil {
			return err
		}

		// insert to the entity information table
		query = `INSERT INTO entity_information (entity_id, slug)
		VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, entityId, entity.EntityInformation.Slug.Value)
		if err != nil {
			return err
		}
		return nil
	})
}

func (e *entityStore) FindByEmail(ctx context.Context, email string) (*Entity, error) {
	query := `
	SELECT e.id, e.fullname, e.email, e.password, e.user_type, ei.slug
	FROM entity e
	JOIN entity_information ei ON e.id = ei.entity_id
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	entity := &Entity{}
	err := e.db.QueryRowContext(ctx, query, email).Scan(
		&entity.ID, &entity.FullName, &entity.Email,
		&entity.Password.hash, &entity.UserType,
		&entity.EntityInformation.Slug.Value,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return entity, nil
}

func (e *entityStore) FindById(ctx context.Context, slug string) (*Entity, error) {
	query := `SELECT id, email, ei.slug FROM entity e
	JOIN entity_information ei ON e.id = ei.entity_id
	WHERE slug = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	entity := &Entity{}
	err := e.db.QueryRowContext(ctx, query, slug).Scan(
		&entity.ID, &entity.Email, &entity.EntityInformation.Slug.Value,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return entity, nil
}

func (e *entityStore) CreateRole(ctx context.Context, role *Role) error {
	// verify if the role already exist in the record
	isFound, _ := e.findRoleByName(ctx, role.Name)
	if isFound {
		return ErrRoleNameConflict
	}

	query := `INSERT INTO roles (name, description, type) 
	VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := e.db.ExecContext(ctx, query, role.Name, role.Description, role.Type)
	if err != nil {
		return err
	}

	return nil
}

// assignment of permission to role (recommended)
func (e *entityStore) AssignPermissionToRole(ctx context.Context, role string, permissions []string) error {
	// Generate placeholders ($2, $3, ...)
	placeholders := make([]string, len(permissions))
	args := make([]any, 0, len(permissions)+1)
	args = append(args, role)

	for i, perm := range permissions {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, perm)
	}
	query := fmt.Sprintf(`
		INSERT INTO role_permissions (role_id, permission_id)
		SELECT r.id, p.id
		FROM roles r, permissions p
		WHERE r.name = $1
		AND p.name IN (%s)
	`, strings.Join(placeholders, ", "))

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := e.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoChangesMade
	}

	return nil
}

func (e *entityStore) AssignRoleToUser(ctx context.Context, slug string, role string) error {
	const query = `INSERT INTO user_roles (user_id, role_id)
		SELECT u.id AS user_id, r.id AS role_id
		FROM users u, roles r
		WHERE u.slug = $1
		AND r.name = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := e.db.ExecContext(ctx, query, slug, role)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoChangesMade
	}

	return nil
}

// direct assignment of permission to user (not recommended)
func (e *entityStore) AssignPermissionToUser(ctx context.Context, slug string, permissions []string) error {
	// Generate placeholders ($2, $3, ...)
	placeholders := make([]string, len(permissions))
	args := make([]any, 0, len(permissions)+1)
	args = append(args, slug)

	for i, perm := range permissions {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, perm)
	}
	query := fmt.Sprintf(`
		INSERT INTO user_permissions (user_id, permission_id)
    SELECT u.id AS user_id, p.id AS permission_id
    FROM users u, permissions p
    WHERE u.slug = $1
		AND p.name IN (%s)
	`, strings.Join(placeholders, ", "))

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := e.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoChangesMade
	}

	return nil
}

func (e *entityStore) RemovePermissionFromUser(ctx context.Context, slug string, permission string) error {
	query := `
		DELETE FROM user_permissions
		WHERE user_id = (SELECT id FROM users WHERE slug = $1)
		AND permission_id = (SELECT id FROM permissions WHERE name = $2)
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := e.db.ExecContext(ctx, query, slug, permission)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoChangesMade
	}
	return nil
}

func (e *entityStore) RemovePermissionFromRole(ctx context.Context, role string, permission string) error {
	query := `
    DELETE FROM role_permissions
    WHERE role_id = (SELECT id FROM roles WHERE name = $1)
    AND permission_id = (SELECT id FROM permissions WHERE name = $2)
  `
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := e.db.ExecContext(ctx, query, role, permission)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNoChangesMade
	}
	return nil
}

func (e *entityStore) GetAllRoles(ctx context.Context) ([]Role, error) {
	query := `SELECT id, name, description, type FROM roles`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []Role{}
	for rows.Next() {
		var r Role
		err := rows.Scan(
			&r.ID, &r.Name,
			&r.Description, &r.Type,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}

func (e *entityStore) GetAllPermissions(ctx context.Context) ([]Permission, error) {
	query := `SELECT id, name, description, type FROM permissions`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	permissions := []Permission{}
	for rows.Next() {
		var p Permission
		err := rows.Scan(
			&p.ID, &p.Name,
			&p.Description, &p.Type,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// private method to check if the role already exist in the record
func (e *entityStore) findRoleByName(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM roles WHERE name = $1)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var exists bool
	err := e.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, err
	}

	// check if the role exists
	if exists {
		return true, nil
	}
	return false, nil
}
