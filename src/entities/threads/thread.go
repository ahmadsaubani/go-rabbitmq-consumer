package threads

import (
	"subscriber-topic-stars/src/entities/comments"
	"subscriber-topic-stars/src/entities/thread_likes"
	"subscriber-topic-stars/src/entities/users"
	"time"
)

type Thread struct {
	UUID        string                    `gorm:"type:uuid;uniqueIndex" json:"uuid"`
	ID          uint64                    `gorm:"primaryKey" db:"id,primary,serial" json:"id"`
	UserID      uint64                    `gorm:"not null;index" db:"user_id" json:"user_id"`
	User        *users.User               `gorm:"foreignKey:UserID" db:"-" json:"user"`
	Title       string                    `gorm:"size:255;not null" json:"title"`
	Description string                    `gorm:"type:text;not null" json:"description"`
	Comments    []comments.Comment        `gorm:"foreignKey:ThreadID" json:"comments"`
	Likes       []thread_likes.ThreadLike `gorm:"foreignKey:ThreadID" json:"-"`

	CreatedAt time.Time `gorm:"autoCreateTime" db:"created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" db:"updated_at" json:"updated_at"`
}
