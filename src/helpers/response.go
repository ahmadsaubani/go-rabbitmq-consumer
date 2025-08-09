package helpers

import "encoding/json"

type RPCResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func Response(success bool, message string, data map[string]interface{}) ([]byte, error) {
	res := RPCResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
	return json.Marshal(res)
}
