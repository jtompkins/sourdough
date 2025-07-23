package recipes

import (
	"context"
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

const LLM_SYSTEM_PROMPT = `
	You are a helpful model that specializes in formatting recipe text into a structured JSON output.
	When you are given text input, if it looks like a recipe, you will do the following steps:
		1. Clean up the formatting of individual ingredients, normalizing the measurements to American standards
		2. Simplify individual steps in the instructions where it makes sense, but DO NOT remove or skip steps
		3. If you cannot determine a value for any of fields, output an empty string ("") for the value, DO NOT substitute any other value or skip the field
		4. Return your modified version of the recipe in JSON format, adhering to the following schema:
			{
				"title": "string",
				"prepTime": "string", # in hours and minutes
				"cookTime": "string", # in hours and minutes
				"servings": "number",
				"ingredients": [
					"string"
				],
				"instructions": [
					"string"
				]
			}
`

type LLMService struct {
	client *openai.Client
	model  string
}

func NewLLMService(client *openai.Client, model string) *LLMService {
	return &LLMService{
		client: client,
		model:  model,
	}
}

func (s *LLMService) FormatRecipe(recipeText string) (LLMRecipe, error) {
	var llmRecipe LLMRecipe

	schema, err := jsonschema.GenerateSchemaForType(llmRecipe)
	if err != nil {
		return LLMRecipe{}, err
	}

	req := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: LLM_SYSTEM_PROMPT,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: recipeText,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "recipe",
				Schema: schema,
				Strict: true,
			},
		},
	}

	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		req,
	)

	if err != nil {
		return LLMRecipe{}, err
	}

	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &llmRecipe)

	if err != nil {
		return LLMRecipe{}, err
	}

	return llmRecipe, nil
}