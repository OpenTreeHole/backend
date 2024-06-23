package model

import (
	"database/sql"
	"time"
)

type User struct {
	// primary key
	ID int `json:"id"`

	// user registration time
	JoinedTime time.Time `gorm:"autoCreateTime"`

	// user last login time
	LastLogin time.Time `gorm:"autoUpdateTime"`

	// user nickname; designed
	Nickname string `gorm:"default:user;size:32"`

	// encrypted email, using pbkdf2.Key + sha3.New512 + hex.EncodeToString, 128 length
	Identifier sql.NullString `gorm:"size:128;uniqueIndex:,length:10"`

	// encrypted password, using pbkdf2.Key + sha256.New + base64.StdEncoding, 78 length
	Password string `gorm:"size:128"`

	// user jwt secret
	UserJwtSecret string

	// check whether user is active or use has been deleted
	IsActive bool `gorm:"default:true"`

	// check whether user is admin, deprecated after use RBAC
	IsAdmin bool `gorm:"not null;default:false;index"`

	// check whether user has completed registration test
	HasCompletedRegistrationTest bool `gorm:"not null;default:false"`
}

type DeleteIdentifier struct {
	// primary key, reference to user.id
	UserID int `gorm:"primaryKey"`

	// encrypted email, using pbkdf2.Key + sha3.New512 + hex.EncodeToString, 128 length, unique
	Identifier string `gorm:"size:128;uniqueIndex:,length:10"`
}
