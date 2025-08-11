package thread_services

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/thread_dtos"
	"subscriber-topic-stars/src/entities/comments"
	"subscriber-topic-stars/src/entities/thread_likes"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/repositories/thread_repositories"

	"github.com/gin-gonic/gin"
)

type ThreadServiceInterface interface {
	CreateThread(token map[string]interface{}, title string, description string) (interface{}, error)
	GetAllThreads(token map[string]interface{}) (map[string]interface{}, error)
	GetThreadDetail(token map[string]interface{}, threadID string) (map[string]interface{}, error)
	LikeThreadService(token map[string]interface{}, threadID string) (*thread_likes.ThreadLike, error)
}

type ThreadService struct {
	threadRepo thread_repositories.ThreadRepositoryInterface
}

func NewThreadService(repo thread_repositories.ThreadRepositoryInterface) *ThreadService {
	return &ThreadService{threadRepo: repo}
}
func (s *ThreadService) CreateThread(token map[string]interface{}, title string, description string) (interface{}, error) {

	userID, err := helpers.ConvertTokenToUserId(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	thread, err := s.threadRepo.Create(title, description, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to create thread: %w", err)
	}

	return gin.H{
		"id":          thread.ID,
		"uuid":        thread.UUID,
		"title":       thread.Title,
		"description": thread.Description,
		"user_id":     thread.UserID,
	}, nil
}

func (s *ThreadService) GetAllThreads(token map[string]interface{}) (map[string]interface{}, error) {
	_, err := helpers.ConvertTokenToUserId(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	threads, err := s.threadRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve threads: %w", err)
	}

	var response []thread_dtos.ThreadResponseDto
	for _, thread := range threads {
		response = append(response, thread_dtos.ThreadResponseDto{
			UUID:        thread.UUID,
			ID:          thread.ID,
			Title:       thread.Title,
			Description: thread.Description,
			CreatedBy:   thread.User.Name,
			CreatedAt:   thread.CreatedAt,
		})
	}

	// Bungkus response dalam map
	result := map[string]interface{}{
		"list_forums": response,
		"total_forum": len(response),
	}

	return result, nil
}

func (s *ThreadService) GetThreadDetail(token map[string]interface{}, threadID string) (map[string]interface{}, error) {
	_, err := helpers.ConvertTokenToUserId(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	thread, err := s.threadRepo.FindThreadDetailByUUID(threadID)
	if err != nil {
		return nil, fmt.Errorf("thread not found: %w", err)
	}

	var comments []map[string]interface{}
	for _, c := range thread.Comments {
		if c.ParentID == nil {
			comments = append(comments, buildCommentTree(c))
		}
	}

	likes, err := s.threadRepo.CountLikes(thread.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count likes: %w", err)
	}

	result := map[string]interface{}{
		"title":         thread.Title,
		"description":   thread.Description,
		"created_by":    thread.User.Name,
		"created_at":    thread.CreatedAt.Format("2006-01-02 15:04:05"),
		"total_comment": len(comments),
		"total_like":    likes,
		"comment_list":  comments,
	}

	return result, nil
}

func (s *ThreadService) LikeThreadService(token map[string]interface{}, threadID string) (*thread_likes.ThreadLike, error) {
	userID, err := helpers.ConvertTokenToUserId(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	thread, err := s.threadRepo.FindThreadDetailByUUID(threadID)
	if err != nil {
		return nil, fmt.Errorf("thread not found: %w", err)
	}

	like, err := s.threadRepo.AddLike(uint64(thread.ID), uint64(userID))
	return like, err
}

func buildCommentTree(comment comments.Comment) map[string]interface{} {
	var replies []map[string]interface{}
	for _, r := range comment.Replies {
		replies = append(replies, buildCommentTree(r))
	}

	return map[string]interface{}{
		"comment":     comment.Comment,
		"comment_by":  comment.User.Name,
		"created_at":  comment.CreatedAt.Format("2006-01-02 15:04:05"),
		"total_reply": len(replies),
		"reply_list":  replies,
	}
}
