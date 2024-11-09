package coder

import (
	"fmt"

	"github.com/jwhenry28/LLMAgents/coding-buddy/prompts"
	local "github.com/jwhenry28/LLMAgents/coding-buddy/tools"
	"github.com/jwhenry28/LLMAgents/shared/conversation"
	"github.com/jwhenry28/LLMAgents/shared/llm"
	"github.com/jwhenry28/LLMAgents/shared/model"
	"github.com/jwhenry28/LLMAgents/shared/tools"
)

type Coder struct {
	llm llm.LLM
}

func NewCoder(llm llm.LLM) *Coder {
	return &Coder{llm: llm}
}

func registerTools() {
	tools.RegisterTool("help", tools.NewHelp)
	tools.Registry["read"] = local.NewRead
	tools.Registry["write"] = local.NewWrite
	tools.Registry["run"] = local.NewRun
	tools.Registry["finish"] = local.NewFinish
	tools.Registry["goget"] = local.NewGoGet
	tools.Registry["gotidy"] = local.NewGoTidy
	tools.Registry["fetch"] = local.NewFetch
}

func (c *Coder) Code() {
	conversationIsOver := func(c *conversation.Conversation) bool {
		modelResponse := c.Messages[len(c.Messages)-2]
		selectedTool, _ := model.ToolInputFromJSON(modelResponse.Content)

		return selectedTool.Name == "finish"
	}

	registerTools()

	messages := c.initialMessages()
	conversation := conversation.NewConversation(c.llm, messages, conversationIsOver)
	conversation.RunConversation()
}

func (c *Coder) initialMessages() []model.Chat {
	userPrompt := c.getUserPrompt()
	return []model.Chat{
		{Role: "system", Content: fmt.Sprintf(prompts.SYSTEM_CODER, tools.GetToolList())},
		{Role: "user", Content: userPrompt},
	}
}

func (c *Coder) getUserPrompt() string {
	// fmt.Print("What would you like me to help you code? ")
	// scanner := bufio.NewScanner(os.Stdin)
	// prompt := ""
	// if scanner.Scan() {
	// 	prompt = scanner.Text()
	// }

	// return prompt
	// return "I want you to write a program that prints 'Hello, world!'"
	return `I want you to write a Gocolly Collector.OnHTML callback function that extracts data from Hacker News. You
	should scrape the root page of Hacker News (https://news.ycombinator.com) and extract all featured articles. Parse
	the articles into a JSON list, where each object is formatted like so:

	{"title": "Anchor tag inner text", "url": "Anchor tag href"}

	Your callback should ignore any links that are not featured articles, such as /jobs, /newcomments, or /submit.

	You may find Gocolly's 'Getting Started' documentation helpful: https://go-colly.org/docs/introduction/start/

	Use the 'fetch' tool to read Gocolly's documentation or to inspect the Hacker News page. I recommend reviewing 
	both materials before writing any Golang code.

	Please do not submit a program until you have tested it, and it returns the expected output (a list of parsed anchor tags).
	`
}
