package comment_repositories

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/comment_dtos"
	"subscriber-topic-stars/src/entities/comments"
	"subscriber-topic-stars/src/entities/threads"
	"subscriber-topic-stars/src/entities/users"
	"subscriber-topic-stars/src/helpers"
)

type CommentRepositoryInterface interface {
	Create(req comment_dtos.CreateCommentRequest) (*comments.Comment, error)
	FindUserById(userID uint64) (*users.User, error)
	FindParentByUUID(parentID string) (*comments.Comment, error)
	FindThreadByUUID(threadID string) (*threads.Thread, error)
	FindCommentsByThreadID(threadID uint64) ([]comments.Comment, error)
}

type commentRepository struct{}

func NewCommentRepository() commentRepository {
	return commentRepository{}
}

func (r commentRepository) Create(req comment_dtos.CreateCommentRequest) (*comments.Comment, error) {
	fmt.Println("di repo :", req)
	user, err := r.FindUserById(req.UserID)
	if err != nil {
		return nil,
			err
	}

	thread, err := r.FindThreadByUUID(req.ThreadID)
	if err != nil {
		return nil,
			err
	}

	newComment := comments.Comment{
		ThreadID: uint64(thread.ID),
		UserID:   uint64(user.ID),
		Comment:  req.Comment,
	}

	if req.ParentID != nil {

		parent, err := r.FindParentByUUID(*req.ParentID)
		if err != nil {
			return nil,
				err
		}
		if parent != nil {
			newComment.ParentID = &parent.ID
		}
	}

	if err := helpers.InsertModel(&newComment); err != nil {
		return nil, err
	}
	return &newComment, nil
}

func (r commentRepository) FindUserById(userID uint64) (*users.User, error) {
	var user users.User

	err := helpers.FindOneByField(&user, "id", userID)
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)

	}
	return &user, nil
}

func (r commentRepository) FindParentByUUID(parentID string) (*comments.Comment, error) {
	var comment comments.Comment
	if parentID == "" {

		return nil,
			nil
	}

	err := helpers.FindOneByField(&comment, "uuid", parentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}

	fmt.Println(comment)
	return &comment, nil
}

func (r commentRepository) FindThreadByUUID(threadID string) (*threads.Thread, error) {
	var thread threads.Thread

	err := helpers.FindOneByField(&thread, "uuid", threadID)
	if err != nil {
		return nil, fmt.Errorf("thread not found: %w", err)
	}
	return &thread, nil
}

func (r commentRepository) FindCommentsByThreadID(threadID uint64) ([]comments.Comment, error) {
	var varComment []comments.Comment
	err := helpers.GettingAllModels(&varComment, []string{"User", "Thread"}, "thread_id = ?", threadID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve comments: %w", err)
	}
	return varComment, nil
}
