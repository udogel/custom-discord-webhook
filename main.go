package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	WebhookURL  string `json:"webhook_url"`
	Username    string `json:"username,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	Embed       Embed  `json:"embed"`
	TimeoutTime int    `json:"timeout_time,omitempty"`
}

func main() {
	config, err := readConfig("config.json")
	if err != nil {
		log.Fatal("Error loading configuration file: ", err)
	}

	if err := validateConfig(config); err != nil {
		log.Fatal("Invalid configuration: ", err)
	}

	timeoutTime := config.TimeoutTime
	if timeoutTime == 0 {
		timeoutTime = 10
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutTime) * time.Second,
	}

	if err := PostWebhook(config, client, timeoutTime); err != nil {
		log.Fatal("Error sending webhook: ", err)
	} else {
		log.Print("Webhook sent successfully")
	}
}

func readConfig(path string) (Config, error) {
	var config Config
	configFile, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err = json.Unmarshal(configFile, &config); err != nil {
		return config, fmt.Errorf("failed to unmarshall config file: %w", err)
	}

	return config, nil
}

func validateConfig(config Config) error {
	if config.WebhookURL == "" || config.Embed.Title == "" || config.Embed.Description == "" {
		return fmt.Errorf("all required config fields must be set")
	}
	for _, v := range config.Embed.Fields {
		if v.Name == "" || v.Value == "" {
			return fmt.Errorf("all required config fields must be set")
		}
	}
	return nil
}
