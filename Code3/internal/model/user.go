package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// define the db instance user struct
type userStore struct {
	db *sql.DB
}

var (
	ErrUserDuplicateEmail = errors.New("a user with the given email already exists")
)

// By default, admin user is created with full permission
func (u *userStore) CreateAdmin(ctx context.Context, user *User) error {
	// check if the user already exist in the record
	isFound, _ := u.FindByEmail(ctx, user.Email)
	if isFound != nil {
		return ErrUserDuplicateEmail
	}
	return withTransaction(ctx, u.db, func(tx *sql.Tx) error {
		query := `INSERT INTO users (fullname, email, phone, dob, gender, address, occupation, slug, id_type, id_number, user_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING slug, id`

		ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
		defer cancel()
		// insert to the users table
		err := tx.QueryRowContext(
			ctx, query, user.FullName, user.Email, user.Phone, user.DOB, user.Gender, user.Address,
			user.Occupation, user.Slug.Value, user.IdType, user.IdNumber, user.UserType,
		).Scan(&user.Slug.Value, &user.ID)
		if err != nil {
			return err
		}
		// get the admin role id from the roles table
		var rowId int64
		query = `SELECT id FROM roles WHERE name = 'administrator'`
		err = tx.QueryRowContext(ctx, query).Scan(&rowId)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return errors.New("admin role not found")
			default:
				return err
			}
		}
		// insert the user to the user-role table with full admin permission
		query = `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
		_, err = tx.ExecContext(ctx, query, user.ID, rowId)
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *userStore) Create(ctx context.Context, user *User) error {
	// check if the user already exist in the record
	isFound, _ := u.FindByEmail(ctx, user.Email)
	if isFound != nil {
		return ErrUserDuplicateEmail
	}

	query := `INSERT INTO users (fullname, email, phone, dob, gender, address, occupation, slug, user_type, id_type, id_number) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING slug`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := u.db.QueryRowContext(
		ctx,
		query,
		user.FullName, user.Email, user.Phone, user.DOB, user.Gender,
		user.Address, user.Occupation, user.Slug.Value, user.UserType,
		user.IdType, user.IdNumber,
	).Scan(&user.Slug.Value)

	if err != nil {
		return err
	}
	return nil
}

func (u *userStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT	id, fullname, email, slug, user_type FROM users
	WHERE email = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}

	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.FullName, &user.Email,
		&user.Slug.Value, &user.UserType,
	)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userStore) FindById(ctx context.Context, slug string) (*User, error) {
	query := `SELECT id, email, slug FROM users
	WHERE slug = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	user := &User{}
	err := u.db.QueryRowContext(ctx, query, slug).Scan(
		&user.ID, &user.Email, user.Slug.Value,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (u *userStore) VerifyPermission(ctx context.Context, userId int64, permissions []string) (bool, error) {
	// Generate placeholders ($2, $3, ...)
	placeholders := make([]string, len(permissions))
	args := make([]any, 0, len(permissions)+1)
	args = append(args, userId)

	for i, perm := range permissions {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args = append(args, perm)
	}

	query := fmt.Sprintf(`
		SELECT 1 FROM (
			SELECT p.name FROM user_permissions up
			JOIN permissions p ON up.permission_id = p.id
			WHERE up.user_id = $1
			UNION
			SELECT p.name FROM user_roles ur
			JOIN role_permissions rp ON ur.role_id = rp.role_id
			JOIN permissions p ON rp.permission_id = p.id
			WHERE ur.user_id = $1
		) AS perms
		WHERE name IN (%s)
		LIMIT 1
	`, strings.Join(placeholders, ", "))

	var found int
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := u.db.QueryRowContext(ctx, query, args...).Scan(&found)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, ErrNotFound
		default:
			return false, err
		}
	}

	return true, nil
}
