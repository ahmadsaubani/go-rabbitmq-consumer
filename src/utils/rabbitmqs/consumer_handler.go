package rabbitmqs

import (
	"log"
	"subscriber-topic-stars/src/providers"
)

const (
	AuthLoginRequest       = "auth.login.request"
	AuthRegisterRequest    = "auth.register.request"
	UserProfileRequest     = "user.profile.request"
	ThreadCreateRequest    = "thread.create.request"
	ThreadGetAllRequest    = "thread.getAll.request"
	ThreadLikeRequest      = "thread.like.request"
	ThreadGetDetailRequest = "thread.getDetail.request"
	CommentCreateRequest   = "comment.create.request"
)

type Consumer struct {
	app providers.AppProvider
}

func NewConsumerService() Consumer {
	return Consumer{
		app: providers.Register(),
	}
}

func (c Consumer) StartConsumers() error {
	err := InitRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ init error: %v", err)
	}
	err = StartRPCConsumer(AuthLoginRequest, "", c.app.Handlers.Auth.LoginRPCHandler())
	//err = StartRPCConsumer("auth.login.request", "", c.app.Handlers.Auth.LoginRPCHandler())
	if err != nil {
		log.Fatalf("Failed to start login request consumer: %v", err)
	}

	if err := StartRPCConsumer(AuthRegisterRequest, "", c.app.Handlers.Auth.RegisterRPCHandler()); err != nil {
		//if err := StartRPCConsumer("auth.register.request", "", c.app.Handlers.Auth.RegisterRPCHandler()); err != nil {
		log.Fatalf("Failed to start register request consumer: %v", err)
	}

	if err := StartRPCConsumer(UserProfileRequest, "", c.app.Handlers.User.UserProfileRPCHandler()); err != nil {
		//if err := StartRPCConsumer("user.profile.request", "", c.app.Handlers.User.UserProfileRPCHandler()); err != nil {
		log.Fatalf("Failed to start user profile consumer: %v", err)
	}

	if err := StartRPCConsumer("thread.create.request", "", c.app.Handlers.Thread.CreateThreadRPCHandler()); err != nil {
		log.Fatalf("Failed to start thread create consumer: %v", err)
	}
	if err := StartRPCConsumer("thread.getAll.request", "", c.app.Handlers.Thread.GetAllThreadHandler()); err != nil {
		log.Fatalf("Failed to start thread get all consumer: %v", err)
	}
	if err := StartRPCConsumer("thread.like.request", "", c.app.Handlers.Thread.LikeThreadHandler()); err != nil {
		log.Fatalf("Failed to start thread like consumer: %v", err)
	}
	if err := StartRPCConsumer("thread.getDetail.request", "", c.app.Handlers.Thread.GetThreadDetailHandler()); err != nil {
		log.Fatalf("Failed to start thread get detail consumer: %v", err)
	}

	if err := StartRPCConsumer("comment.create.request", "", c.app.Handlers.Comment.CreateCommentRPCHandler()); err != nil {
		log.Fatalf("Failed to start comment create consumer: %v", err)
	}
	return nil
}
