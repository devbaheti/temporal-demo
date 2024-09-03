package dsl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"temporal-demo/models"
	"temporal-demo/utils"

	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"
)

type SampleActivities struct {
}

func (a *SampleActivities) HttpBlock(ctx context.Context, url []string) (string, error) {
	fmt.Println("HttpBlock")

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", url[0], nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", "XXXXXXXXXXXX")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	zap.L().Info("processing response", zap.Any("responsePayload", resp.Body), zap.Any("statusCode", resp.StatusCode))

	// Check if the response status code indicates an error
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		zap.L().Error("processing response :: error", zap.Any("responsePayload", resp.Body), zap.Any("statusCode", string(bodyBytes)))
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the response
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}
	fmt.Println(response)

	finalResp := ""

	for k, v := range response {
		finalResp = finalResp + k + ":" + utils.ConvertInterfaceToString(v)
	}

	return finalResp, nil
}

func (a *SampleActivities) ProcessChat(ctx context.Context, text []string) (string, error) {

	prompt := fmt.Sprintf("Text: %s "+" Generate context of the text. Strictly do not include more than 100 words. "+" Strictly include only context of the text", text[0])

	messages := []models.Message{
		{
			Role:    "system",
			Content: "You are a helpful assistant which helps to generate context of a text provided by the user",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	requestPayload := &models.ChatCompletionRequest{}
	requestPayload.Messages = messages
	// Set the default model if none is provided
	if requestPayload.Model == "" {
		requestPayload.Model = "gpt-3.5-turbo-0613"
	}

	// Set the default temperature if not provided
	if requestPayload.Temperature == 0 {
		requestPayload.Temperature = 0.7
	}

	zap.L().Info("ProcessChat :: processing completion request", zap.Any("requestPayload", requestPayload))

	// Convert the payload to JSON
	jsonPayload, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request payload: %v", err)
	}

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/chat/completions", "https://api.openai.com/v1"), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "XXXXXXXXXXXX"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	zap.L().Info("ProcessChat :: processing completion response", zap.Any("responsePayload", resp.Body), zap.Any("statusCode", resp.StatusCode))

	// Check if the response status code indicates an error
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		zap.L().Error("ProcessChat :: processing completion response :: error", zap.Any("responsePayload", resp.Body), zap.Any("statusCode", string(bodyBytes)))
		return "", fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the response
	var chatCompletionResponse models.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatCompletionResponse); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, chatCompletionResponse.Choices[0].Message.Content)
	return chatCompletionResponse.Choices[0].Message.Content, nil
}

func (a *SampleActivities) SampleActivity1(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity2(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity3(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity4(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity5(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}
