package thread_like_dtos

type ThreadLikeRequestDto struct {
	Token    string `json:"token"`
	ThreadID string `json:"thread_id" binding:"required"`
}
