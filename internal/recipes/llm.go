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
		4. If the recipe you're given is missing cook time or prep time, output an empty string ("") for the value, DO NOT substitute any other value or skip the field
		5. Return your modified version of the recipe in JSON format, adhering to the following schema:
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

const LLM_IMAGE_SYSTEM_PROMPT = `
	You are a helpful model that specializes in extracting recipe information from images and formatting it into structured JSON output.
	When you are given an image that contains a recipe, you will do the following steps:
		1. Extract all visible text from the image, paying special attention to ingredients lists and cooking instructions
		2. Clean up the formatting of individual ingredients, normalizing the measurements to American standards
		3. Organize the instructions into clear, sequential steps
		4. If you cannot determine a value for any of the fields, output an empty string ("") for the value, DO NOT substitute any other value or skip the field
		5. If the recipe is missing cook time or prep time, output an empty string ("") for the value, DO NOT substitute any other value or skip the field
		6. Return your extracted and formatted version of the recipe in JSON format, adhering to the following schema:
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

func (s *LLMService) FormatRecipeFromImage(base64Image, contentType string) (LLMRecipe, error) {
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
				Content: LLM_IMAGE_SYSTEM_PROMPT,
			},
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL:    "data:" + contentType + ";base64," + base64Image,
							Detail: openai.ImageURLDetailHigh,
						},
					},
				},
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
