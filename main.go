package main

import (
	"fmt"
	"os"
	"os/signal"
	"subscriber-topic-stars/src/configs"
	"subscriber-topic-stars/src/seeders"
	"subscriber-topic-stars/src/utils/rabbitmqs"
	"subscriber-topic-stars/src/utils/redis"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	db := configs.ConnectDatabase()
	configs.RunMigrations(db.Gorm)
	seeders.Run(db)

	redisError := redis.InitRedis()
	if redisError != nil {
		fmt.Println("Error initializing redis: " + redisError.Error())
	}

	consumerError := rabbitmqs.NewConsumerService().StartConsumers()
	if consumerError != nil {
		fmt.Errorf("could not start consumer: %w", consumerError)
	}

	// Wait signal (CTRL+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
