package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type ChatGPTAIService struct {
	APIKey string
}

func NewChatGPTAIService() *ChatGPTAIService {
	return &ChatGPTAIService{
		APIKey: os.Getenv("OPENAI_API_KEY"),
	}
}

type chatGPTRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatGPTResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

func (s *ChatGPTAIService) GenerateBlogIdeas(topic string) (string, error) {
	if s.APIKey == "" {
		return "", errors.New("OPENAI_API_KEY not set")
	}

	reqBody := chatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []chatMessage{
			{Role: "system", Content: "You are a helpful assistant that generates blog ideas."},
			{Role: "user", Content: fmt.Sprintf("Generate blog ideas about: %s", topic)},
		},
	}

	bodyBytes, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("OpenAI API error: %s", respBody)
	}

	var chatResp chatGPTResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}
	if len(chatResp.Choices) == 0 {
		return "", errors.New("no ideas generated")
	}
	return chatResp.Choices[0].Message.Content, nil
}
