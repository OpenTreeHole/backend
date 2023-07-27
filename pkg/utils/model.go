package utils

type MessageResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Message return MessageResponse
func Message(msg string) *MessageResponse {
	return &MessageResponse{
		Message: msg,
	}
}
