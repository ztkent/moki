package aiclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

// Ensure we have the right env variables set for the given source
func MustLoadAPIKey(openai bool, anyscale bool) error {
	// Load the .env file if we don't have the env var set
	loadEnvVar := func(varName string) error {
		if os.Getenv(varName) == "" {
			err := godotenv.Load()
			if err != nil || os.Getenv(varName) == "" {
				return fmt.Errorf("Failed to load %s", varName)
			}
		}
		return nil
	}

	if openai {
		if err := loadEnvVar("OPENAI_API_KEY"); err != nil {
			return err
		}
	}
	if anyscale {
		if err := loadEnvVar("ANYSCALE_ENDPOINT_TOKEN"); err != nil {
			return err
		}
	}
	return nil
}

func MustConnectOpenAI(model OpenAIModel, temperature float32) *Client {
	oaiClient := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	MustCheckConnection(oaiClient)
	return &Client{Client: oaiClient, Model: model.String(), Temperature: temperature}
}

func MustConnectAnyscale(model AnyscaleModel, temperature float32) *Client {
	config := openai.DefaultConfig(os.Getenv("ANYSCALE_ENDPOINT_TOKEN"))
	config.BaseURL = "https://api.endpoints.anyscale.com/v1"
	asClient := openai.NewClientWithConfig(config)
	MustCheckConnection(asClient)
	return &Client{Client: asClient, Model: model.String(), Temperature: temperature}
}

func MustCheckConnection(client *openai.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := client.ListModels(ctx)
	if err != nil {
		panic(err)
	}
}
