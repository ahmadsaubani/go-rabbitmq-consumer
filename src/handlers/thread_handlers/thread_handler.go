package thread_handlers

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/thread_dtos"
	"subscriber-topic-stars/src/dtos/thread_like_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services"
	"subscriber-topic-stars/src/utils/redis"
	"time"

	"encoding/json"
)

type ThreadHandler interface {
	CreateThreadRPCHandler() func([]byte) ([]byte, error)
	GetAllThreadHandler() func([]byte) ([]byte, error)
	GetThreadDetailHandler() func([]byte) ([]byte, error)
	LikeThreadHandler() func([]byte) ([]byte, error)
}

type threadHandler struct {
	services services.ServiceCenter
}

func NewThreadHandler(services services.ServiceCenter) ThreadHandler {
	return threadHandler{services: services}
}

func (h threadHandler) CreateThreadRPCHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req thread_dtos.ThreadRequestDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			resp := helpers.RPCResponse{Success: false, Message: "Invalid request format"}
			return json.Marshal(resp)
		}

		msg := map[string]interface{}{
			"token": req.Token,
		}

		result, err := h.services.Thread.CreateThread(msg, req.Title, req.Description)
		if err != nil {
			resp := helpers.RPCResponse{Success: false, Message: err.Error()}
			return json.Marshal(resp)
		}

		redis.DelKey(fmt.Sprintf("threads:list"))

		return json.Marshal(result)
	}
}

func (t threadHandler) GetAllThreadHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req thread_dtos.ThreadRequestDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			resp := helpers.RPCResponse{Success: false, Message: "Invalid request format"}
			return json.Marshal(resp)
		}

		msg := map[string]interface{}{
			"token": req.Token,
		}

		// Key Redis untuk semua thread
		cacheKey := "threads:list"

		if cached, err := redis.GetKey(cacheKey); err == nil && cached != "" {
			fmt.Println("Serving all threads from Redis cache")
			return []byte(cached), nil
		}

		result, err := t.services.Thread.GetAllThreads(msg)
		if err != nil {
			fmt.Println("Error retrieving threads:", err)
			resp := helpers.RPCResponse{Success: false, Message: err.Error()}
			return json.Marshal(resp)
		}

		// 5 menit
		if resultJSON, err := json.Marshal(result); err == nil {
			fmt.Printf("added cache for threads list")
			_ = redis.SetKey(cacheKey, resultJSON, 5*time.Minute)
		}

		return json.Marshal(result)
	}
}

func (t threadHandler) GetThreadDetailHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {

		var req thread_dtos.ThreadDetailRequestDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		key := fmt.Sprintf("thread:detail:%s", req.ThreadID)

		cachedData, err := redis.GetKey(key)
		if err == nil && cachedData != "" {
			fmt.Printf("get cache for thread detail : %s\n", req.ThreadID)
			return []byte(cachedData), nil
		}

		msg := map[string]interface{}{
			"token": req.Token,
		}
		result, err := t.services.Thread.GetThreadDetail(msg, req.ThreadID)
		if err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			}
			return json.Marshal(resp)
		}

		if resultJSON, err := json.Marshal(result); err == nil {
			fmt.Printf("set cache for thread detail: %s\n", req.ThreadID)
			_ = redis.SetKey(key, resultJSON, 5*time.Minute)
		}

		return json.Marshal(result)
	}
}

func (t threadHandler) LikeThreadHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req thread_like_dtos.ThreadLikeRequestDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		key := fmt.Sprintf("thread:detail:%s", req.ThreadID)

		redis.DelKey(key)

		token := map[string]interface{}{"token": req.Token}

		result, err := t.services.Thread.LikeThreadService(token, req.ThreadID)
		if err != nil {
			resp := helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			}
			return json.Marshal(resp)
		}

		return json.Marshal(result)
	}
}
