package user_handlers

import (
	"encoding/json"
	"fmt"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services"
	"subscriber-topic-stars/src/utils/redis"
	"time"
)

type Request struct {
	Token string `json:"token"`
}

type UserHandler interface {
	UserProfileRPCHandler() func([]byte) ([]byte, error)
}

type userHandler struct {
	services services.ServiceCenter
}

func NewUserHandler(services services.ServiceCenter) UserHandler {
	return userHandler{services: services}
}

func (h userHandler) UserProfileRPCHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req Request
		if err := json.Unmarshal(requestBody, &req); err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			})
		}

		cacheKey := fmt.Sprintf("user:profile:%s", req.Token)

		// Cek di Redis
		cachedData, err := redis.GetKey(cacheKey)
		if err == nil && cachedData != "" {
			fmt.Println("Cache hit for user profile:", req.Token)
			return []byte(cachedData), nil
		}

		// Input untuk service
		msg := map[string]interface{}{
			"token": req.Token,
		}

		result, err := h.services.User.GetUser(msg)
		if err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		_ = redis.SetKey(cacheKey, resultJSON, 5*time.Minute)

		return resultJSON, nil
	}
}
