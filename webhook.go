package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Color       int          `json:"color,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
}

type WebhookRequestBody struct {
	Username  string  `json:"username,omitempty"`
	Embeds    []Embed `json:"embeds"`
	AvatarURL string  `json:"avatar_url,omitempty"`
}

func PostWebhook(config Config, client *http.Client, timeoutTime int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutTime)*time.Second)
	defer cancel()

	if err := sendWebhook(config, client, ctx); err != nil {
		return err
	}

	return nil
}

func sendWebhook(config Config, client *http.Client, ctx context.Context) error {
	requestBody := WebhookRequestBody{
		Username:  config.Username,
		AvatarURL: config.AvatarURL,
		Embeds:    []Embed{config.Embed},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", config.WebhookURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d, response body: %s", resp.StatusCode, string(body))
	}

	return nil
}
