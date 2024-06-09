package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/mr-karan/expenseai/pkg/models"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type Manager struct {
	log    *slog.Logger
	client *openai.Client
	model  string
}

func New(token, baseURL, model string, log *slog.Logger) (*Manager, error) {
	// Initialize the OpenAI client.
	cfg := openai.DefaultConfig(token)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	cfg.HTTPClient.Timeout = 10 * time.Second
	client := openai.NewClientWithConfig(cfg)

	return &Manager{
		client: client,
		model:  model,
		log:    log,
	}, nil

}

// Parse the message and extract the expenses.
func (m *Manager) Parse(msg string) (models.Transactions, error) {
	if msg == "" {
		return models.Transactions{}, errors.New("empty message")
	}

	m.log.Debug("Parsing expenses", "message", msg)
	dialogue := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: fmt.Sprintf("Categorize my expenses. Make a best case effort to deduce the nature of expenses from the given description even if a little ambiguous or vague. Ignore any messages not related to user's spends and return a valid error. Today's date is %s", time.Now().Format("2006-01-02"))},
		// {Role: openai.ChatMessageRoleSystem, Content: fmt.Sprintf("The user will input their spends on various things, you need to split the expenses in the correct categories. The amount and a short description must be enough for you to parse the expenses correctly. If the message from user does not contain any details about their spends, return a valid error. Assume transaction date to be today's date if unspecified. Today's date is %s", time.Now().Format("2006-01-02"))},
		{Role: openai.ChatMessageRoleUser, Content: msg},
	}

	fnCategorizeExpenses := openai.FunctionDefinition{
		Name:        "categorize_expense",
		Description: "Categorize expenses from the given input.",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"transactions": {
					Type:        jsonschema.Array,
					Description: "List of items purchased",
					Items: &jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"transaction_date": {
								Type:        jsonschema.String,
								Description: "Date of transaction in ISO 8601 format (e.g., 2021-09-01) if specified else today's date.",
							},
							"amount": {
								Type:        jsonschema.Number,
								Description: "Amount of the item",
							},
							"category": {
								Type:        jsonschema.String,
								Description: "One word category of the expense (e.g., food, travel, entertainment)",
							},
							"description": {
								Type:        jsonschema.String,
								Description: "Concise and short description of the item",
							},
						},
						Required: []string{"transaction_date", "amount", "category", "description"},
					},
				},
			},
			Required: []string{"transactions"},
		},
	}

	t := openai.Tool{
		Type:     openai.ToolTypeFunction,
		Function: &fnCategorizeExpenses,
	}

	resp, err := m.client.CreateChatCompletion(context.TODO(),
		openai.ChatCompletionRequest{
			Model:    m.model,
			Messages: dialogue,
			Tools:    []openai.Tool{t},
		},
	)

	if err != nil || len(resp.Choices) != 1 {
		m.log.Error("Completion error", "error", err, "choices", len(resp.Choices))
		return models.Transactions{}, fmt.Errorf("error completing the request")
	}

	var transactions models.Transactions

	for _, choice := range resp.Choices {
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Function.Name == fnCategorizeExpenses.Name {
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &transactions); err != nil {
					return models.Transactions{}, fmt.Errorf("error unmarshalling response: %s", err)
				}
				return transactions, nil
			}
		}
	}

	if len(transactions.Transactions) == 0 {
		for _, choice := range resp.Choices {
			if choice.FinishReason == "stop" {
				return models.Transactions{}, &NoValidTransactionError{Message: choice.Message.Content}
			}
		}
	}

	return models.Transactions{}, fmt.Errorf("no valid transactions found in response")
}

type NoValidTransactionError struct {
	Message string
}

func (e *NoValidTransactionError) Error() string {
	return e.Message
}
