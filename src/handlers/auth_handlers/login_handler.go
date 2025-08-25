package auth_handlers

import (
	"encoding/json"
	"subscriber-topic-stars/src/dtos/auth_dtos"
	"subscriber-topic-stars/src/helpers"
	"subscriber-topic-stars/src/services"
)

type AuthHandler interface {
	LoginRPCHandler() func([]byte) ([]byte, error)
	RegisterRPCHandler() func([]byte) ([]byte, error)
}

type authHandler struct {
	services services.ServiceCenter
}

func NewAuthHandler(services services.ServiceCenter) AuthHandler {
	return authHandler{services: services}
}

func (h authHandler) LoginRPCHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req auth_dtos.RequestLoginDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			})
		}

		result, err := h.services.Auth.Login(req.Email, req.Password)
		if err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return json.Marshal(result)
	}
}

func (h authHandler) RegisterRPCHandler() func([]byte) ([]byte, error) {
	return func(requestBody []byte) ([]byte, error) {
		var req auth_dtos.RequestRegisterDto
		if err := json.Unmarshal(requestBody, &req); err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: "Invalid request format",
			})
		}

		result, err := h.services.Auth.Register(req)
		if err != nil {
			return json.Marshal(helpers.RPCResponse{
				Success: false,
				Message: err.Error(),
			})
		}
		return json.Marshal(result)
	}
}
