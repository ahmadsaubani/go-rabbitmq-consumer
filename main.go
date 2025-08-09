package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"subscriber-topic-stars/src/configs"
	"subscriber-topic-stars/src/handlers/auth_handlers"
	"subscriber-topic-stars/src/handlers/comment_handlers"
	"subscriber-topic-stars/src/handlers/thread_handlers"
	"subscriber-topic-stars/src/handlers/user_handlers"
	"subscriber-topic-stars/src/repositories/auth_repositories"
	"subscriber-topic-stars/src/repositories/comment_repositories"
	"subscriber-topic-stars/src/repositories/thread_repositories"
	"subscriber-topic-stars/src/repositories/user_repositories"
	"subscriber-topic-stars/src/seeders"
	"subscriber-topic-stars/src/services/auth_services"
	"subscriber-topic-stars/src/services/comment_services"
	"subscriber-topic-stars/src/services/thread_services"
	"subscriber-topic-stars/src/services/user_services"
	"subscriber-topic-stars/src/utils/rabbitmqs"
	"subscriber-topic-stars/src/utils/redis"

	"github.com/joho/godotenv"
)

func main() {
	err := rabbitmqs.InitRabbitMQ()
	if err != nil {
		log.Fatalf("RabbitMQ init error: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redis.InitRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), redisDB)

	db := configs.ConnectDatabase()
	configs.RunMigrations(db.Gorm)
	seeders.Run(db)

	repo := auth_repositories.NewAuthRepository()
	authService := auth_services.NewAuthService(repo)

	// Pasangkan handler Login
	err = rabbitmqs.StartRPCConsumer("auth.login.request", "", auth_handlers.LoginRPCHandler(authService))
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// USER
	userRepo := user_repositories.NewUserRepositoryInterface()
	userService := user_services.NewUserService(userRepo)

	err = rabbitmqs.StartRPCConsumer("user.profile.request", "", user_handlers.UserProfileRPCHandler(userService))
	if err != nil {
		log.Fatalf("Failed to start user profile consumer: %v", err)
	}

	// THREAD
	threadRepo := thread_repositories.NewThreadRepository()
	threadService := thread_services.NewThreadService(threadRepo)

	err = rabbitmqs.StartRPCConsumer("thread.create.request", "", thread_handlers.CreateThreadRPCHandler(threadService))
	if err != nil {
		log.Fatalf("Failed to start thread create consumer: %v", err)
	}

	err = rabbitmqs.StartRPCConsumer("thread.getAll.request", "", thread_handlers.GetAllThreadHandler(threadService))
	if err != nil {
		log.Fatalf("Failed to start thread get all consumer: %v", err)
	}

	err = rabbitmqs.StartRPCConsumer("thread.like.request", "", thread_handlers.LikeThreadHandler(threadService))
	if err != nil {
		log.Fatalf("Failed to start thread like consumer: %v", err)
	}

	err = rabbitmqs.StartRPCConsumer("thread.getDetail.request", "", thread_handlers.GetThreadDetailHandler(threadService))
	if err != nil {
		log.Fatalf("Failed to start thread get detail consumer: %v", err)
	}

	// COMMENT
	commentRepo := comment_repositories.NewCommentRepository()
	commentService := comment_services.NewCommentService(commentRepo)
	err = rabbitmqs.StartRPCConsumer("comment.create.request", "", comment_handlers.CreateCommentRPCHandler(commentService))
	if err != nil {
		log.Fatalf("Failed to start comment create consumer: %v", err)
	}

	// Wait signal (CTRL+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
