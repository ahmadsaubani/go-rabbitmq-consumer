package users

import "time"

type User struct {
	UUID      string    `gorm:"uniqueIndex" db:"uuid" json:"uuid"`
	ID        uint64    `gorm:"primaryKey;autoIncrement" db:"id,primary,serial" json:"id"`
	Email     string    `gorm:"size:255;unique;not null" db:"email" json:"email" binding:"required,email"`
	Name      string    `gorm:"size:255;not null" db:"name" json:"name" binding:"required,min=3,max=255"`
	Username  string    `gorm:"size:255;unique;not null" db:"username" json:"username" binding:"required,min=3,max=255"`
	Password  string    `gorm:"size:255;not null" db:"password" json:"password" binding:"required,min=6"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
