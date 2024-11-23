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

type Claude struct {
	ApiKey      string
	Model       string
	Temperature int

	apiUrl string
}

func NewClaude(apikey string, model string, temperature int) *Claude {
	return &Claude{
		ApiKey:      apikey,
		Model:       model,
		Temperature: temperature,
		apiUrl:      "https://api.anthropic.com",
	}
}

func (llm *Claude) CompleteChat(messages []model.Chat) (string, error) {
	return llm.completeChat(messages, DEFAULT_RETRIES)
}

func (llm *Claude) completeChat(messages []model.Chat, retries int) (string, error) {
	endpoint := llm.apiUrl + "/v1/messages"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       llm.Model,
		"messages":    messages[1:],
		"temperature": llm.Temperature,
		"system":      llm.getSystemMessage(messages),
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

	choices := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", fmt.Errorf("no choices returned from API")
	}

	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})

	return message["content"].(string), nil
}

func (llm *Claude) getSystemMessage(messages []model.Chat) string {
	return messages[0].Content
}

func (llm *Claude) getRetryDelay(errorResponse string) float64 {
	re := regexp.MustCompile(`try again in (\d+\.?\d*)s`)
	matches := re.FindStringSubmatch(errorResponse)
	if len(matches) > 1 {
		seconds, _ := strconv.ParseFloat(matches[1], 64)
		return seconds
	}
	return 0.0
}
