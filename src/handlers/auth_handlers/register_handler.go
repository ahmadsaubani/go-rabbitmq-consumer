package auth_handlers

import (
	"encoding/json"
	"subscriber-topic-stars/src/dtos/auth_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services/auth_services"
)

func RegisterRPCHandler(authService auth_services.AuthServiceInterface) func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req auth_dtos.RequestRegisterDto
		if err := json.Unmarshal(requestBody, &req); err != nil {

			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		result, err := authService.Register(req)

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
