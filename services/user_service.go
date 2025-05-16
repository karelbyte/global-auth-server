package services

import (
	"fmt"
	"global-auth-server/libs"

	"github.com/jmoiron/sqlx"
)

// User represents the structure of the users table
type User struct {
	ID              string  `db:"id" json:"id"`
	Username        string  `db:"username" json:"username"`
	Code            *string `db:"code" json:"code,omitempty"`
	Names           string  `db:"names" json:"names"`
	Email           string  `db:"email" json:"email"`
	Password        *string `db:"password" json:"-"`
	RolID           *string `db:"rol_id" json:"rol_id,omitempty"`
	IsStaff         bool    `db:"is_staff" json:"is_staff"`
	IsActive        bool    `db:"is_active" json:"is_active"`
	BossID          *string `db:"boss_id" json:"boss_id,omitempty"`
	CreatedAt       string  `db:"created_at" json:"-"`
	UpdatedAt       string  `db:"updated_at" json:"-"`
	Token           *string `db:"token" json:"-"`
	Logins          *int    `db:"logins" json:"logins,omitempty"`
	CanDownloadXlsx *bool   `db:"can_download_xlsx" json:"can_download_xlsx,omitempty"`
	BankID          *string `db:"bank_id" json:"bank_id,omitempty"`
	FilialID        *string `db:"filial_id" json:"filial_id,omitempty"`
}

type Role struct {
	Code        string `db:"code" json:"code"`
	Description string `db:"description" json:"description"`
}

// Getuserbyemail looks for a user by email using SQLX and the Singleton connection
func GetUserByEmail(email string) (*User, error) {
	db := libs.GetDB()
	sqlxdb := sqlx.NewDb(db, "postgres")

	var user User
	err := sqlxdb.Get(&user, `
		SELECT 
			id::text,
			username, 
			code, 
			names, 
			email, 
			password, 
			rol_id::text, 
			is_staff, 
			is_active, 
			boss_id::text, 
			created_at, 
			updated_at, 
			token, 
			logins, 
			can_download_xlsx, 
			bank_id::text, 
			filial_id::text
		FROM users 
		WHERE email = $1
	`, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}

func GetRolesByUserID(userID string) ([]Role, error) {
	db := libs.GetDB()
	sqlxdb := sqlx.NewDb(db, "postgres")

	var roles []Role
	err := sqlxdb.Select(&roles, `
		SELECT r.code, r.description
		FROM user_roles ur
		INNER JOIN rols r ON ur.rol_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("roles not found: %w", err)
	}
	return roles, nil
}
