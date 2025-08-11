package thread_dtos

import (
	"subscriber-topic-stars/src/dtos/comment_dtos"
	"time"
)

type ThreadRequestDto struct {
	Token       string `json:"token" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type ThreadResponseDto struct {
	UUID        string    `json:"uuid"`
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type ThreadDetailResponse struct {
	CommentList  []comment_dtos.CommentResponse `json:"comment_list"`
	CreatedAt    string                         `json:"created_at"`
	CreatedBy    string                         `json:"created_by"`
	Description  string                         `json:"description"`
	Title        string                         `json:"title"`
	TotalComment int                            `json:"total_comment"`
	TotalLikes   int                            `json:"total_likes"`
}

type ThreadDetailRequestDto struct {
	Token    string `json:"token" binding:"required"`
	ThreadID string `json:"thread_id" binding:"required"`
}
