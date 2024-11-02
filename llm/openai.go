package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"hackandpray.com/media-curator/model"
)

type OpenAI struct {
	ApiKey      string
	Model       string
	Temperature int

	apiUrl string
}

func NewOpenAI(apikey string, model string, temperature int) *OpenAI {
	return &OpenAI{
		ApiKey:      apikey,
		Model:       model,
		Temperature: temperature,
		apiUrl:      "https://api.openai.com",
	}
}

// TODO: call request in loop N times in case of failure or ratelimiting.
func (llm *OpenAI) CompleteChat(messages []model.Chat) (string, error) {
	endpoint := llm.apiUrl + "/v1/chat/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":       llm.Model,
		"messages":    messages,
		"temperature": llm.Temperature,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+llm.ApiKey)

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
