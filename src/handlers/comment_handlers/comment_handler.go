package comment_handlers

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/comment_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services"
	"subscriber-topic-stars/src/utils/redis"

	"encoding/json"
)

type CommentHandler interface {
	CreateCommentRPCHandler() func([]byte) ([]byte, error)
}

type commentHandler struct {
	services services.ServiceCenter
}

func NewCommentHandler(services services.ServiceCenter) CommentHandler {
	return commentHandler{services: services}
}

func (h commentHandler) CreateCommentRPCHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req comment_dtos.CreateCommentRequest
		if err := json.Unmarshal(requestBody, &req); err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			})
		}

		token := map[string]interface{}{
			"token": req.Token,
		}

		result, err := h.services.Comment.CreateComment(token, req)

		key := fmt.Sprintf("thread:detail:%s", req.ThreadID)
		redis.DelKey(key)

		if err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return json.Marshal(result)
	}
}
