package thread_likes

import "time"

type ThreadLike struct {
	UUID      string    `gorm:"type:uuid;uniqueIndex" json:"uuid"`
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ThreadID  uint64    `gorm:"not null;index" json:"thread_id"`
	UserID    uint64    `gorm:"not null;index;uniqueIndex:idx_user_thread" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
