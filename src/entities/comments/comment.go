package comments

import (
	"subscriber-topic-stars/src/entities/users"
	"time"
)

type Comment struct {
	ID       uint64     `gorm:"primaryKey" json:"id"`
	UUID     string     `gorm:"type:uuid;uniqueIndex" json:"uuid"`
	ThreadID uint64     `gorm:"not null;index" json:"thread_id"`
	ParentID *uint64    `gorm:"index" json:"parent_id"`
	Parent   *Comment   `gorm:"foreignKey:ParentID" json:"-"`
	UserID   uint64     `gorm:"not null;index" json:"user_id"`
	User     users.User `gorm:"foreignKey:UserID" json:"-"`
	Replies  []Comment  `gorm:"foreignKey:ParentID" json:"replies"`

	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
