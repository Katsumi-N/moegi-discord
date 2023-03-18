package chatgpt

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"moegi-discord/config"
	"net/http"
)

type SendingMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model       string           `json:"model"`
	Messages    []SendingMessage `json:"messages"`
	Temperature float64          `json:"temperature"`
}

type ChatGPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

const (
	CHAT_COMPLETIONS_ENDPOINT = "https://api.openai.com/v1/chat/completions"
)

func RequestToOpenAI(sending_messages []SendingMessage) (resp ChatGPTResponse, err error) {
	client := &http.Client{}
	data := ChatGPTRequest{
		Model:       "gpt-3.5-turbo",
		Messages:    sending_messages,
		Temperature: 0.7,
	}

	json_data, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
		return ChatGPTResponse{}, err
	}
	req, err := http.NewRequest("POST", CHAT_COMPLETIONS_ENDPOINT, bytes.NewBuffer(json_data))
	if err != nil {
		log.Fatal(err)
		return ChatGPTResponse{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.Config.OpenAISecretKey)

	raw_resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return ChatGPTResponse{}, err
	}
	defer raw_resp.Body.Close()
	body, err := ioutil.ReadAll(raw_resp.Body)
	if err != nil {
		log.Fatal(err)
		return ChatGPTResponse{}, err
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		log.Fatal(err)
		return ChatGPTResponse{}, err
	}

	return resp, err
}
