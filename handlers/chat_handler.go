package handlers

import (
	"chatbot/models"
	"chatbot/services"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	volcService *services.VolcengineService
}

func NewChatHandler(volcService *services.VolcengineService) *ChatHandler {
	return &ChatHandler{volcService: volcService}
}

func RegisterRoutes(r *gin.Engine, volcService *services.VolcengineService) {
	handler := NewChatHandler(volcService)
	r.POST("/api/chat", handler.Chat)
	r.POST("/api/chat/stream", handler.ChatStream)
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 100001,
			"msg":  err.Error(),
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 100001,
			"msg":  "问题不能为空",
		})
		return
	}

	content, reasoningContent, err := h.volcService.Chat(c.Request.Context(), req.Message)
	if err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 100001,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ChatResponse{
		Content:          content,
		ReasoningContent: reasoningContent,
	})
}

func (h *ChatHandler) ChatStream(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 100001,
			"msg":  err.Error(),
		})
		return
	}

	if req.Message == "" {
		c.JSON(http.StatusOK, map[string]interface{}{
			"code": 100001,
			"msg":  "问题不能为空",
		})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	streamChan, errChan := h.volcService.ChatStream(c.Request.Context(), req.Message)

	c.Stream(func(w io.Writer) bool {
		select {
		case content, ok := <-streamChan:
			if !ok {
				c.SSEvent("message", "[DONE]")
				return false
			}
			c.SSEvent("message", content)
			return true
		case err, ok := <-errChan:
			if ok {
				c.SSEvent("error", err.Error())
			}
			return false
		}
	})
}
