package user_handlers

import (
	"encoding/json"
	"fmt"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services/user_services"
	"subscriber-topic-stars/src/utils/redis"
	"time"
)

type Request struct {
	Token string `json:"token"`
}

// []byte digunakan utk menangani komunikasi RPC via RabbitMQ, yang bekerja dengan payload berbentuk byte array
func UserProfileRPCHandler(userService user_services.UserServiceInterface) func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		// Struktur permintaan dari publisher

		var req Request
		if err := json.Unmarshal(requestBody, &req); err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		cacheKey := fmt.Sprintf("user:profile:%s", req.Token)

		// Cek di Redis
		cachedData, err := redis.GetKey(cacheKey)
		if err == nil && cachedData != "" {
			fmt.Println("Cache hit for user profile:", req.Token)
			return []byte(cachedData), nil
		}

		// karena service meminta map[string]interface{} sebagai input
		msg := map[string]interface{}{
			"token": req.Token,
		}

		result, err := userService.GetUser(msg)
		if err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			}
			return json.Marshal(resp)
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		// Simpan ke Redis 5 menit aja broww
		_ = redis.SetKey(cacheKey, resultJSON, 5*time.Minute)

		return resultJSON, nil
	}
}
