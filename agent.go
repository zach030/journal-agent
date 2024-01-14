package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	log "github.com/sirupsen/logrus"
)

const prompt = `你是一位博学多识，了解科技与心理的专家，我会给你我最近一段时间记录的思考和笔记内容（用json数组组成），我希望你可以基于这些内容，帮我做几件事：
				1. 整理出所有涉及情绪的内容，并进行专业的心理分析，给我一些心理辅助和调节的建议;
				2. 整理出所有笔记、想法类的内容，进行总结，提取出核心观点并润色，尝试寻找出内容的关联性，给出一些创新独特的思考点;
				3. 整理出所有我的复盘、总结、反思类内容，提取总结我的反思内容，给出可行的执行计划和方案;
				4. 基于所有我给定的内容和你的心情，发挥你的想象力，任意地表达你的想法。`

type FlashNote struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Mtime   time.Time `json:"mtime"`
}

func LoadFlashNote(dir string) []FlashNote {
	notes := make([]FlashNote, 0)
	now := time.Now()
	earlestTime := now.AddDate(0, 0, -8)
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Errorf("read dir error: %v", err)
		return nil
	}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Errorf("get file info error: %v filename: %s", err, entry.Name())
			continue
		}
		if info.ModTime().Before(earlestTime) {
			log.Warnf("skip stale file: %s", entry.Name())
			continue
		}
		log.Infof("start load note: %s", entry.Name())
		content, err := os.ReadFile(dir + "/" + entry.Name())
		if err != nil {
			log.Errorf("read file error: %v filename: %s", err, entry.Name())
			continue
		}
		notes = append(notes, FlashNote{
			Title:   entry.Name(),
			Content: string(content),
			Mtime:   info.ModTime(),
		})
	}
	return notes
}

type JournalAgent struct {
	client *openai.Client
	notes  []FlashNote
}

func NewJournalAgent(apiKey, apiBase, noteDir string) *JournalAgent {
	cfg := openai.DefaultConfig(apiKey)
	if apiBase != "" {
		cfg.BaseURL = apiBase
	}
	client := openai.NewClientWithConfig(cfg)
	return &JournalAgent{client: client, notes: LoadFlashNote(noteDir)}
}

func (a *JournalAgent) NotesReview() (string, error) {
	str, _ := json.Marshal(a.notes)
	msg := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: string(str),
		},
	}
	nums := NumTokensFromMessages(msg, openai.GPT4TurboPreview)
	log.Info("estimated use tokens: ", nums)
	resp, err := a.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			MaxTokens: 4095,
			Model:     openai.GPT4TurboPreview,
			Messages:  msg,
		},
	)
	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (a *JournalAgent) ParseReview(review string) (*AgentReview, error) {
	tools := []openai.Tool{
		{
			Type:     openai.ToolTypeFunction,
			Function: JournalFunction(),
		},
	}
	request := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: fmt.Sprintf("这是你针对我的笔记给出的总结内容: %s, 帮我按照给定的函数格式进行解析", review),
			},
		},
		Tools: tools,
	}
	resp, err := a.client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return nil, err
	}
	if len(resp.Choices) == 0 || len(resp.Choices[0].Message.ToolCalls) == 0 {
		return nil, errors.New("empty response")
	}
	call := resp.Choices[0].Message.ToolCalls[0].Function
	if call.Name != "parse_agent_review" {
		log.Errorf("unexpected call name: %s", call.Name)
		return nil, errors.New("unexpected call name")
	}
	ret := &AgentReview{}
	if err = json.Unmarshal([]byte(call.Arguments), ret); err != nil {
		log.Errorf("json unmarshal error: %v", err)
		return nil, err
	}
	return ret, nil
}

type AgentReview struct {
	Mood     string `json:"mood"`
	Note     string `json:"note"`
	Plan     string `json:"plan"`
	Creative string `json:"creative"`
}

func JournalFunction() openai.FunctionDefinition {
	return openai.FunctionDefinition{
		Name: "parse_agent_review",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"mood": {
					Type:        jsonschema.String,
					Description: "The mood aspect of the review",
				},
				"note": {
					Type:        jsonschema.String,
					Description: "The note aspect of the review",
				},
				"plan": {
					Type:        jsonschema.String,
					Description: "The note aspect of the review",
				},
				"creative": {
					Type:        jsonschema.String,
					Description: "The creative aspect of the review",
				},
			},
			Required: []string{"mood", "note", "plan", "creative"},
		},
	}
}

func NumTokensFromMessages(messages []openai.ChatCompletionMessage, model string) (numTokens int) {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		err = fmt.Errorf("encoding for model: %v", err)
		log.Println(err)
		return
	}

	var tokensPerMessage, tokensPerName int
	switch model {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
		tokensPerName = 1
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
		tokensPerName = -1   // if there's a name, the role is omitted
	default:
		if strings.Contains(model, "gpt-3.5-turbo") {
			log.Println("warning: gpt-3.5-turbo may update over time. Returning num tokens assuming gpt-3.5-turbo-0613.")
			return NumTokensFromMessages(messages, "gpt-3.5-turbo-0613")
		} else if strings.Contains(model, "gpt-4") {
			log.Println("warning: gpt-4 may update over time. Returning num tokens assuming gpt-4-0613.")
			return NumTokensFromMessages(messages, "gpt-4-0613")
		} else {
			err = fmt.Errorf("num_tokens_from_messages() is not implemented for model %s. See https://github.com/openai/openai-python/blob/main/chatml.md for information on how messages are converted to tokens.", model)
			log.Println(err)
			return
		}
	}

	for _, message := range messages {
		numTokens += tokensPerMessage
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokensPerName
		}
	}
	numTokens += 3 // every reply is primed with <|start|>assistant<|message|>
	return numTokens
}
