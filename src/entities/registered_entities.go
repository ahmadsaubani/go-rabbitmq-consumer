package entities

import (
	"subscriber-topic-stars/src/entities/auth"
	"subscriber-topic-stars/src/entities/comments"
	"subscriber-topic-stars/src/entities/thread_likes"
	"subscriber-topic-stars/src/entities/threads"
	"subscriber-topic-stars/src/entities/users"
)

var RegisteredEntities = []any{
	auth.AccessToken{},
	auth.RefreshToken{},
	users.User{},
	threads.Thread{},
	comments.Comment{},
	thread_likes.ThreadLike{},
}
