package coder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jwhenry28/LLMAgents/coding-buddy/prompts"
	local "github.com/jwhenry28/LLMAgents/coding-buddy/tools"
	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
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
	tools.Registry["goinit"] = local.NewGoInit
	tools.Registry["goget"] = local.NewGoGet
	tools.Registry["gotidy"] = local.NewGoTidy
	tools.Registry["fetch"] = local.NewFetch
	tools.Registry["report"] = local.NewReport
	tools.Registry["finish"] = local.NewFinish
}

func setupSandbox() {
	if err := os.MkdirAll("sandbox", 0755); err != nil {
		panic("failed to create sandbox directory")
	}

	dir, err := os.ReadDir("sandbox")
	if err != nil {
		panic("failed to read sandbox directory")
	}

	for _, entry := range dir {
		if err := os.RemoveAll(filepath.Join(utils.SANDBOX_DIR, entry.Name())); err != nil {
			panic("failed to remove sandbox entry")
		}
	}
}

func (c *Coder) Code() {
	registerTools()
	setupSandbox()

	// c.generateSuccessCriteria()
	// successCriteria := c.generateSuccessCriteria()
	c.generateCode("nil")
}

func (c *Coder) generateSuccessCriteria() string {
	conversationIsOver := func(c conversation.Conversation) bool {
		messages := c.GetMessages()
		modelResponse := messages[len(messages)-2]
		selectedTool, _ := model.NewTextToolInput(modelResponse.Content)

		return selectedTool.GetName() == "report"
	}

	messages := c.initialSuccessCriteriaMessages()
	messages = c.runConversation(messages, conversationIsOver)
	return messages[len(messages)-1].Content
}

func (c *Coder) generateCode(successCriteria string) {
	conversationIsOver := func(c conversation.Conversation) bool {
		messages := c.GetMessages()
		modelResponse := messages[len(messages)-2]
		selectedTool, _ := model.NewTextToolInput(modelResponse.Content)

		return selectedTool.GetName() == "finish"
	}

	messages := c.initialCodeMessages(successCriteria)
	_ = c.runConversation(messages, conversationIsOver)
}

func (c *Coder) runConversation(initMessages []model.Chat, isOver func(conversation.Conversation) bool) []model.Chat {
	conversation := conversation.NewChatConversation(c.llm, initMessages, isOver)
	conversation.RunConversation()

	return conversation.GetMessages()
}

func (c *Coder) initialSuccessCriteriaMessages() []model.Chat {
	userPrompt := "PROBLEM STATEMENT:\n" + c.getUserPrompt()
	return []model.Chat{
		{Role: "system", Content: fmt.Sprintf(prompts.SYSTEM_PRBLM, tools.GetToolList())},
		{Role: "user", Content: userPrompt},
	}
}

func (c *Coder) initialCodeMessages(successCriteria string) []model.Chat {
	userPrompt := `INSTRUCTION LIST:
1. Import necessary packages: net/url, github.com/gocolly/colly, and any other required packages.

2. Define the Anchor struct in a separate package named model with the following fields:
   - Text string
   - HRef string

3. Define the HackerNewsScraper struct with the following fields:
   - URL *url.URL: to store the URL of the page being scraped.
   - Anchors []model.Anchor: to store the list of anchors found on the page.
   - InnerText string: to store the inner text of the page.
   - collector *colly.Collector: to store the Gocolly object used to scrape the page.

4. Implement a constructor function NewHackerNewsScraper that initializes a HackerNewsScraper object:
   - Parse the Hacker News URL using url.Parse and assign it to the URL field.
   - Initialize the collector field with colly.NewCollector().

5. Implement a method Scrape for the HackerNewsScraper struct:
   - Use the collector.OnHTML method to find all <a> tags on the page.
   - For each <a> tag, extract the href attribute and the inner text.
   - Create a new model.Anchor object with the extracted href and inner text, and append it to the Anchors slice.
   - Use the collector.OnResponse method to capture the entire response body and assign it to the InnerText field.

6. Implement a method GetFeaturedArticles for the HackerNewsScraper struct:
   - Filter the Anchors slice to include only those anchors that are likely to be featured articles. This can be done by checking if the HRef is a full URL (starts with http or https) and the Text is not empty.
   - Return the filtered list of model.Anchor objects.

7. In the main function, create an instance of HackerNewsScraper using the constructor.
   - Call the Scrape method to perform the scraping.
   - Retrieve the featured articles using the GetFeaturedArticles method.
   - Print or process the list of featured articles as needed.`
	return []model.Chat{
		{Role: "system", Content: fmt.Sprintf(prompts.SYSTEM_CODE, tools.GetToolList())},
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
	return `
Scrape all featured news articles from Hacker News (https://news.ycombinator.com).

You should encapsulate your scraper into a HackerNewsScraper struct that exposes the following methods and fields:
type HackerNewsScraper struct {
	URL       *url.URL  // the URL of the page being scraped
	Anchors   []model.Anchor // the list of anchors found on the page; Anchor has fields Text and HRef
	InnerText string // the inner text of the page

	collector *colly.Collector // the Gocolly object used to scrape the page
}
`
	//	return `
	//
	// Write a Golang function "WhoisDomain" that does the following:
	// 1. Takes in a target domain name and a WHOIS server.
	// 2. Queries the server for the target domain's WHOIS information.
	// 3. Parses the response into a Go-friendly data structure.
	// 4. If the queried server is the authorative server, return the parsed response.
	// 5. Otherwise, recursively query the referred WHOIS server in the response until you reach the authoritative server.
	// 6. If your WHOIS query fails, use exponential backoff to wait for a short time before trying again. Retry a maximum of 5 times.
	// `
}
