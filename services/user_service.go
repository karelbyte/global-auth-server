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
	Password        *string `db:"password,omitempty"`
	RolID           *string `db:"rol_id" json:"rol_id,omitempty"`
	IsStaff         bool    `db:"is_staff" json:"is_staff"`
	IsActive        bool    `db:"is_active" json:"is_active"`
	BossID          *string `db:"boss_id" json:"boss_id,omitempty"`
	CreatedAt       string  `db:"created_at" json:"created_at"`
	UpdatedAt       string  `db:"updated_at" json:"updated_at"`
	Token           *string `db:"token" json:"token,omitempty"`
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
	sqlxdb := sqlx.NewDb(db, "sqlserver")

	var user User
	err := sqlxdb.Get(&user, `
		SELECT 
			CAST(id AS varchar(36)) as id,
			username, 
			code, 
			names, 
			email, 
			password, 
			CAST(rol_id AS varchar(36)) as rol_id, 
			is_staff, 
			is_active, 
			boss_id, 
			created_at, 
			updated_at, 
			token, 
			logins, 
			can_download_xlsx, 
			bank_id, 
			filial_id
		FROM users WHERE email = @p1
	`, email)
	if err != nil {
		return nil, fmt.Errorf("User not found: %w", err)
	}
	return &user, nil
}

func GetRolesByUserID(userID string) ([]Role, error) {
	db := libs.GetDB()
	sqlxdb := sqlx.NewDb(db, "sqlserver")

	var roles []Role
	err := sqlxdb.Select(&roles, `
		SELECT r.code, r.description
		FROM user_roles ur
		INNER JOIN rols r ON ur.rol_id = r.id
		WHERE ur.user_id = @p1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("roles not found: %w", err)
	}
	return roles, nil
}
