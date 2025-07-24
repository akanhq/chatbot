package models

type ChatRequest struct {
	Message string `json:"message"`
	Model   string `json:"model"`
}

type ChatResponse struct {
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
