package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

type Anthropic struct {
	ApiKey      string
	Model       string
	Temperature int

	apiUrl string
}

func NewAnthropic(apikey string, model string, temperature int) *Anthropic {
	return &Anthropic{
		ApiKey:      apikey,
		Model:       model,
		Temperature: temperature,
		apiUrl:      "https://api.anthropic.com",
	}
}

func (llm *Anthropic) Type() string {
	return "anthropic"
}

func (llm *Anthropic) CompleteChat(messages []model.Chat) (string, error) {
	return llm.completeChat(messages, DEFAULT_RETRIES)
}

func (llm *Anthropic) completeChat(messages []model.Chat, retries int) (string, error) {
	endpoint := llm.apiUrl + "/v1/messages"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       llm.Model,
		"messages":    messages[1:],
		"temperature": llm.Temperature,
		"system":      llm.getSystemMessage(messages),
		"max_tokens":  512,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", llm.ApiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		if retries <= 1 {
			return "", fmt.Errorf("rate limit exceeded")
		}
		jitter := 1.0
		duration := llm.getRetryDelay(string(body)) + jitter
		time.Sleep(time.Duration(duration) * time.Second)
		return llm.completeChat(messages, retries-1)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	content := result["content"].([]interface{})
	if len(content) == 0 {
		return "", fmt.Errorf("no content returned from API")
	}

	message := content[0].(map[string]interface{})["text"].(string)

	return message, nil
}

func (llm *Anthropic) getSystemMessage(messages []model.Chat) string {
	return messages[0].Content
}

func (llm *Anthropic) getRetryDelay(errorResponse string) float64 {
	re := regexp.MustCompile(`try again in (\d+\.?\d*)s`)
	matches := re.FindStringSubmatch(errorResponse)
	if len(matches) > 1 {
		seconds, _ := strconv.ParseFloat(matches[1], 64)
		return seconds
	}
	return 0.0
}
