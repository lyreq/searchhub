package models

import (
	"path/filepath"

	"gorm.io/gorm"
)

type StatusType string

const (
	StatusActive  StatusType = "1" // active
	StatusPassive StatusType = "2" // passive
)

type User struct {
	gorm.Model

	Name    string `gorm:"type:varchar(100);"`
	Surname string `gorm:"type:varchar(100);"`
	Avatar  string `gorm:"type:varchar(255);null"`

	Phone    string `gorm:"uniqueIndex;type:varchar(11);not null;"`
	Email    string `gorm:"uniqueIndex;type:varchar(100);not null;"`
	Username string `gorm:"uniqueIndex;type:varchar(50);not null;"`
	Password string `gorm:"type:varchar(255);not null;"`

	Status StatusType `gorm:"type:enum_status_type;default:1;"`
}

func (u *User) AvatarPath() string {
	if u.Avatar == "" {
		// default avatar real path
		return "assets/images/avatars/default.jpg"
	}
	return filepath.Join("uploads", "avatars", u.Avatar)
}
