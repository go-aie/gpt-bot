package gptbot

import (
	"context"

	"github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/chat"
	"github.com/rakyll/openai-go/completion"
)

type ModelType string

const (
	// GPT-4
	GPT4     ModelType = "gpt-4"
	GPT40301 ModelType = "gpt-4-0314"

	// GPT-3.5
	GPT3Dot5Turbo     ModelType = "gpt-3.5-turbo"
	GPT3Dot5Turbo0301 ModelType = "gpt-3.5-turbo-0301"
	TextDavinci003    ModelType = "text-davinci-003"
	TextDavinci002    ModelType = "text-davinci-002"
	// GPT-3
	TextAda001     ModelType = "text-ada-001"
	TextCurie001   ModelType = "text-curie-001"
	TextBabbage001 ModelType = "text-babbage-001"
)

type EngineMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type EngineRequest struct {
	Messages    []*EngineMessage `json:"messages,omitempty"`
	Temperature float64          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
}

type EngineResponse struct {
	Text string `json:"text,omitempty"`
}

type Engine interface {
	Infer(context.Context, *EngineRequest) (*EngineResponse, error)
}

// https://platform.openai.com/docs/models/model-endpoint-compatibility

// OpenAIChatEngine powered by /v1/chat/completions completion api, supported model like:
// `gpt-4`, `gpt-4-0314`, `gpt-3.5-turbo`, `gpt-3.5-turbo-0301` ...
type OpenAIChatEngine struct {
	Client *chat.Client
}

// OpenAICompletionEngine powered by /v1/completions completion api, supported model like:
// `text-davinci-003`, `text-davinci-002`, `text-ada-001`, `text-curie-001`, `text-babbage-001` ...
type OpenAICompletionEngine struct {
	Client *completion.Client
}

func NewOpenAIChatEngine(apiKey string, model ModelType) *OpenAIChatEngine {
	client := chat.NewClient(openai.NewSession(apiKey), string(model))

	return &OpenAIChatEngine{
		Client: client,
	}
}

func NewOpenAICompletionEngine(apiKey string, model ModelType) *OpenAICompletionEngine {
	compClient := completion.NewClient(openai.NewSession(apiKey), string(model))

	return &OpenAICompletionEngine{
		Client: compClient,
	}
}

func (e *OpenAIChatEngine) Infer(ctx context.Context, req *EngineRequest) (resp *EngineResponse, err error) {
	apiResp, err := e.Client.CreateCompletion(ctx, &chat.CreateCompletionParams{
		Messages: []*chat.Message{
			{
				Role:    req.Messages[0].Role,
				Content: req.Messages[0].Content,
			},
		},
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return
	}

	return &EngineResponse{
		Text: apiResp.Choices[0].Message.Content,
	}, nil
}

func (e *OpenAICompletionEngine) Infer(ctx context.Context, req *EngineRequest) (resp *EngineResponse, err error) {
	apiResp, err := e.Client.Create(ctx, &completion.CreateParams{
		Prompt:      []string{req.Messages[0].Content},
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
	})
	if err != nil {
		return
	}

	return &EngineResponse{
		Text: apiResp.Choices[0].Text,
	}, nil
}