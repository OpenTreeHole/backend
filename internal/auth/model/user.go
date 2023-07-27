package model

import (
	"database/sql"
	"time"
)

type User struct {
	// primary key
	ID int `json:"id"`

	// user registration time
	CreatedAt time.Time `json:"created_at"`

	// user last login time
	UpdatedAt time.Time `json:"updated_at"`

	// designed, not using now
	Nickname string `json:"nickname" gorm:"default:user;size:32"`

	// encrypted email, using pbkdf2.Key + sha3.New512 + hex.EncodeToString, 128 length
	Identifier sql.NullString `json:"identifier" gorm:"size:128;uniqueIndex:,length:10"`

	// encrypted password, using pbkdf2.Key + sha256.New + base64.StdEncoding, 78 length
	Password string `json:"password" gorm:"size:128"`

	// check whether user is admin, deprecated after use RBAC
	IsAdmin bool `json:"is_admin" gorm:"not null;default:false;index"`

	// check whether user has completed registration test
	HasCompletedRegistrationTest bool `json:"has_completed_registration_test" gorm:"default:false"`
}
