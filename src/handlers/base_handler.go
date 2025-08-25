package handlers

import (
	"subscriber-topic-stars/src/handlers/auth_handlers"
	"subscriber-topic-stars/src/handlers/comment_handlers"
	"subscriber-topic-stars/src/handlers/thread_handlers"
	"subscriber-topic-stars/src/handlers/user_handlers"
	"subscriber-topic-stars/src/services"
)

type HandlerCenter struct {
	Auth    auth_handlers.AuthHandler
	User    user_handlers.UserHandler
	Thread  thread_handlers.ThreadHandler
	Comment comment_handlers.CommentHandler
}

func InitHandlers(service services.ServiceCenter) HandlerCenter {
	return HandlerCenter{
		Auth:    auth_handlers.NewAuthHandler(service),       // Inisialisasi handler Auth
		User:    user_handlers.NewUserHandler(service),       // Inisialisasi handler User
		Thread:  thread_handlers.NewThreadHandler(service),   // Inisialisasi handler Thread
		Comment: comment_handlers.NewCommentHandler(service), // Inisialisasi handler Comment
	}
}
