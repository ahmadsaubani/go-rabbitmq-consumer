package comment_dtos

type CreateCommentRequest struct {
	Token    string  `json:"token"`
	ThreadID string  `json:"thread_id"`
	ParentID *string `json:"parent_id"`
	UserID   uint64  `json:"user_id"`
	Comment  string  `json:"comment"`
}

type CommentResponse struct {
	Comment    string            `json:"comment"`
	CommentBy  string            `json:"comment_by"`
	CreatedAt  string            `json:"created_at"`
	ReplyList  []CommentResponse `json:"reply_list,omitempty"`
	TotalReply int               `json:"total_reply"`
}
