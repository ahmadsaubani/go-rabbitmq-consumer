package comment_handlers

import (
	"fmt"
	"subscriber-topic-stars/src/dtos/comment_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services/comment_services"
	"subscriber-topic-stars/src/utils/redis"

	"encoding/json"
)

func CreateCommentRPCHandler(commentServices comment_services.CommentServiceInterface) func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req comment_dtos.CreateCommentRequest
		if err := json.Unmarshal(requestBody, &req); err != nil {

			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		token := map[string]interface{}{
			"token": req.Token,
		}

		result, err := commentServices.CreateComment(token, req)

		key := fmt.Sprintf("thread:detail:%s", req.ThreadID)
		redis.DelKey(key)

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
