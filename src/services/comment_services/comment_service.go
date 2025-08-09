package comment_services

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/comment_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/repositories/comment_repositories"
)

type CommentServiceInterface interface {
	CreateComment(token map[string]interface{}, req comment_dtos.CreateCommentRequest) (interface{}, error)
}

type CommentService struct {
	commentRepo comment_repositories.CommentRepositoryInterface
}

func NewCommentService(commentRepo comment_repositories.CommentRepositoryInterface) *CommentService {
	return &CommentService{commentRepo: commentRepo}
}

func (s *CommentService) CreateComment(token map[string]interface{}, req comment_dtos.CreateCommentRequest) (interface{}, error) {
	userID, err := helpers.ConvertTokenToUserId(token)
	if err != nil {
		return nil, err
	}

	req.UserID = userID

	comment, err := s.commentRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return &comment, nil
}
