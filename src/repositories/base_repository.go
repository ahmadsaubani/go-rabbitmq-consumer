package repositories

import (
	"subscriber-topic-stars/src/repositories/auth_repositories"
	"subscriber-topic-stars/src/repositories/comment_repositories"
	"subscriber-topic-stars/src/repositories/thread_repositories"
	"subscriber-topic-stars/src/repositories/user_repositories"
)

type RepositoryCenter struct {
	Auth    auth_repositories.AuthRepositoryInterface
	Comment comment_repositories.CommentRepositoryInterface
	Thread  thread_repositories.ThreadRepositoryInterface
	User    user_repositories.UserRepositoryInterface
}

func InitRepositories() RepositoryCenter {
	return RepositoryCenter{
		Auth:    auth_repositories.NewAuthRepository(),
		Comment: comment_repositories.NewCommentRepository(),
		Thread:  thread_repositories.NewThreadRepository(),
		User:    user_repositories.NewUserRepository(),
	}
}
