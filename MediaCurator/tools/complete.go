package tools

import (
	"encoding/json"
	"fmt"

	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Decision struct {
	Decision      string `json:"decision"`
	Title         string `json:"title"`
	URL           string `json:"url"`
	Justification string `json:"justification"`
}

type Complete struct {
	tools.Base
}

func NewComplete(input model.ToolInput) tools.Tool {
	name := "complete"
	args := []string{"decision1", "decision2", "..."}
	brief := "complete: ends the conversation"
	explanation := `args:
- List of decisions, where each decision is a stringified JSON object like so: "{ \"decision\": \"<decision>\", \"title\": \"<title>\", \"url\": \"<url>\", \"justification\": \"<justification>\" }"
	- title: The anchor tag's title/inner text
	- url: The URL from the anchor tag you are making a decision about
	- decision: Your decision. Must be one of the following:
		- IGNORE: Choose this option if you do not think your client will be interested in reading this URL today.
		- NOTIFY: Choose this option if you would like to forward this URL to your client
	- justification: A short explanation for your decision
Example decision: "{ \"decision\": \"IGNORE\", \"title\": \"My Contoso Blog\", \"url\": \"https://blog.contoso.com/article/1337\", \"justification\": \"The article's content is about starting a generic company, which is not related to the client's interests.\" }"
`
	return Complete{
		Base: tools.Base{Input: input, Name: name, Args: args, BriefText: brief, ExplanationText: explanation},
	}
}

func (task Complete) Match() bool {
	return true
}

func (task Complete) Invoke() string {
	var decision Decision
	for i, arg := range task.Input.GetArgs() {
		err := json.Unmarshal([]byte(arg), &decision)
		if err != nil {
			return fmt.Sprintf("error unmarshaling decision %d: %v", i, err)
		}
	}

	return "completed"
}
