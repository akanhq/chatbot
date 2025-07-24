package services

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"chatbot/config"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type VolcengineService struct {
	client *arkruntime.Client
}

func NewVolcengineService(cfg *config.Config) *VolcengineService {
	client := arkruntime.NewClientWithApiKey(
		cfg.VolcengineAPIKey,
		arkruntime.WithBaseUrl(cfg.VolcengineBaseURL),
		arkruntime.WithRegion(cfg.VolcengineRegion),
		arkruntime.WithTimeout(30*time.Minute),
	)
	return &VolcengineService{client: client}
}

func (s *VolcengineService) Chat(ctx context.Context, message string) (string, string, error) {
	req := model.CreateChatCompletionRequest{
		Model: "deepseek-r1-250528",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String(message),
				},
			},
		},
	}

	resp, err := s.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", "", err
	}

	var content, reasoningContent string
	if resp.Choices[0].Message.ReasoningContent != nil { //推理
		content = *resp.Choices[0].Message.ReasoningContent
	}
	reasoningContent = *resp.Choices[0].Message.Content.StringValue //结果
	return content, reasoningContent, nil
}

func (s *VolcengineService) ChatStream(ctx context.Context, message string) (<-chan string, <-chan error) {
	streamChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(streamChan)
		defer close(errChan)

		req := model.CreateChatCompletionRequest{
			Model: "deepseek-r1-250528",
			Messages: []*model.ChatCompletionMessage{
				{
					Role: model.ChatMessageRoleUser,
					Content: &model.ChatCompletionMessageContent{
						StringValue: volcengine.String(message),
					},
				},
			},
		}

		stream, err := s.client.CreateChatCompletionStream(ctx, req)
		if err != nil {
			errChan <- err
			return
		}
		defer stream.Close()

		buffer := ""
		sentenceEnd := regexp.MustCompile(`[。！？\n]`) // 句子结束标志

		for {
			var recv model.ChatCompletionStreamResponse
			recv, err = stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errChan <- err
				return
			}
			if len(recv.Choices) > 0 {
				var content string
				if recv.Choices[0].Delta.ReasoningContent != nil && len(*recv.Choices[0].Delta.ReasoningContent) > 0 {
					//streamChan <- *recv.Choices[0].Delta.ReasoningContent
					//content = *recv.Choices[0].Delta.ReasoningContent
				} else {
					content = recv.Choices[0].Delta.Content
				}

				if content != "" {

					// 清理可能包含的 SSE 格式关键字
					content = strings.ReplaceAll(content, "event: message", "")
					content = strings.ReplaceAll(content, "data: ", "")
					content = strings.TrimSpace(content)
					if content == "" {
						continue
					}

					buffer += content
					fmt.Println("收到片段:", content) // 调试：打印原始片段
					// 如果缓冲区包含句子结束标志，发送数据
					if sentenceEnd.MatchString(buffer) || len(buffer) > 100 { // 或者缓冲区超过100字符
						fmt.Println("发送完整句子:", buffer)
						streamChan <- buffer
						buffer = ""
					}
				}
			}
		}
	}()

	return streamChan, errChan
}
