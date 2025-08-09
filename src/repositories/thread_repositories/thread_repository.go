package thread_repositories

import (
	"fmt"
	"subscriber-topic-stars/src/entities/thread_likes"
	"subscriber-topic-stars/src/entities/threads"
	"subscriber-topic-stars/src/entities/users"
	"subscriber-topic-stars/src/helpers"

	"gorm.io/gorm"
)

type ThreadRepositoryInterface interface {
	Create(title string, description string, userID uint64) (*threads.Thread, error)
	GetAll() ([]threads.Thread, error)
	FindUserById(userID uint64) (*users.User, error)
	FindThreadDetailByUUID(threadID string) (*threads.Thread, error)
	AddLike(threadID uint64, userID uint64) (*thread_likes.ThreadLike, error)
	CountLikes(threadID uint64) (int64, error)
}

type threadRepository struct{}

func NewThreadRepository() *threadRepository {
	return &threadRepository{}
}

func (r *threadRepository) Create(title string, description string, userID uint64) (*threads.Thread, error) {

	if _, err := r.FindUserById(userID); err != nil {
		return nil, fmt.Errorf("user not found :  %w", err)
	}

	newThread := threads.Thread{
		Title:       title,
		Description: description,
		UserID:      userID,
	}

	if err := helpers.InsertModel(&newThread); err != nil {
		return nil, fmt.Errorf("could not insert user: %w", err)
	}
	return &newThread, nil
}

func (r *threadRepository) FindUserById(userID uint64) (*users.User, error) {
	var user users.User

	err := helpers.FindOneByField(&user, "id", userID)
	if err != nil {
		return nil, fmt.Errorf("email not found: %w", err)

	}
	return &user, nil
}

func (r *threadRepository) GetAll() ([]threads.Thread, error) {

	var threadsList []threads.Thread
	if err := helpers.GettingAllModels(&threadsList, []string{"User"}); err != nil {
		fmt.Println("Error retrieving threads:", err)
		return nil, fmt.Errorf("could not retrieve threads: %w", err)
	}
	return threadsList, nil
}

func (r *threadRepository) FindThreadDetailByUUID(threadID string) (*threads.Thread, error) {
	var thread threads.Thread
	err := helpers.FindOneByFieldWithPreload(
		&thread,
		[]string{
			"User",
			"Comments",                      // komentar utama
			"Comments.User",                 // user yang komentar
			"Comments.Replies",              // nested reply 1 level
			"Comments.Replies.User",         // user yang reply
			"Comments.Replies.Replies",      // nested reply 2 level (optional)
			"Comments.Replies.Replies.User", // user di nested reply 2 level (optional)
		},
		"uuid", threadID,
	)

	return &thread, err
}

func (r *threadRepository) AddLike(threadID uint64, userID uint64) (*thread_likes.ThreadLike, error) {
	like := thread_likes.ThreadLike{
		ThreadID: threadID,
		UserID:   userID,
	}
	if err := helpers.InsertModel(&like); err != nil {
		return nil, fmt.Errorf("could not insert like: %w", err)
	}
	return &like, nil
}

func (r *threadRepository) CountLikes(threadID uint64) (int64, error) {
	return helpers.CountModelWithFilter[thread_likes.ThreadLike](func(db *gorm.DB) *gorm.DB {
		return db.Where("thread_id = ?", threadID)
	})
}
