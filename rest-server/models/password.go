package models

import (
	"time"

	"gorm.io/gorm"
)

// type Password struct {
// 	ID        uint           `json:"id,omitempty"`
// 	Name      string         `json:"name,omitempty" gorm:"unique"`
// 	Username  string         `json:"username,omitempty"`
// 	Email     string         `json:"email,omitempty"`
// 	Password  string         `json:"password,omitempty"`
// 	CreatedAt time.Time      `json:"created_at,omitempty"`
// 	UpdatedAt time.Time      `json:"updated_at,omitempty"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
// }

// Using pointers instead of non-pointer so that we can omit
// undesired field from the json response
// src: - https://stackoverflow.com/a/20413334
//      - https://go.dev/play/p/JJNa0m2_nw
type Password struct {
	ID        *uint           `json:"id,omitempty"`
	Name      *string         `json:"name,omitempty" gorm:"unique"`
	Username  *string         `json:"username,omitempty"`
	Email     *string         `json:"email,omitempty"`
	Password  *string         `json:"password,omitempty"`
	CreatedAt *time.Time      `json:"created_at,omitempty"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
