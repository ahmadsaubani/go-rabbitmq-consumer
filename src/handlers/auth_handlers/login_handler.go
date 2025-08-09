package auth_handlers

import (
	"subscriber-topic-stars/src/dtos/auth_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services/auth_services"

	"encoding/json"
)

func LoginRPCHandler(authService auth_services.AuthServiceInterface) func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req auth_dtos.RequestLoginDto
		if err := json.Unmarshal(requestBody, &req); err != nil {

			resp := helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			}
			return json.Marshal(resp)
		}

		result, err := authService.Login(req.Email, req.Password)

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
