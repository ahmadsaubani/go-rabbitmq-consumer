package services

import (
	"subscriber-topic-stars/src/repositories"
	"subscriber-topic-stars/src/services/auth_services"
	"subscriber-topic-stars/src/services/comment_services"
	"subscriber-topic-stars/src/services/thread_services"
	"subscriber-topic-stars/src/services/user_services"
)

type ServiceCenter struct {
	Auth    auth_services.AuthServiceInterface
	Comment comment_services.CommentServiceInterface
	Thread  thread_services.ThreadServiceInterface
	User    user_services.UserServiceInterface
}

func InitServices(repo repositories.RepositoryCenter) ServiceCenter {
	return ServiceCenter{
		Auth:    auth_services.NewAuthService(repo.Auth),
		Comment: comment_services.NewCommentService(repo.Comment),
		Thread:  thread_services.NewThreadService(repo.Thread),
		User:    user_services.NewUserService(repo.User),
	}
}
