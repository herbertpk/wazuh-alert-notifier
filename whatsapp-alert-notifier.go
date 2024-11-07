package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Alert represents the structure of the alert from Wazuh
type Alert struct {
	Rule struct {
		Description string      `json:"description"`
		Level       int         `json:"level"`
		ID          interface{} `json:"id"` // Accepts both int and string
	} `json:"rule"`
	Agent struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
		ID   string `json:"id"`
	} `json:"agent"`
	Timestamp string `json:"timestamp"`
}

// MessagePayload represents the payload structure for the Green API
type MessagePayload struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

// sendToGreenAPI sends the alert message to the Green API
func sendToGreenAPI(hookURL, chatId, message string) error {
	// Create payload
	payload := MessagePayload{
		ChatID:  chatId,
		Message: message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Set up HTTP request
	req, err := http.NewRequest("POST", hookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 response
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 response from API: %d, response body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Alert sent successfully with response code %d\n", resp.StatusCode)
	return nil
}

func main() {
	// Check if all required arguments are provided
	if len(os.Args) < 4 {
		log.Fatal("Usage: ./custom-greenapi <alert_file> <chatId> <hook_url>")
	}

	// Read arguments
	alertFile := os.Args[1]
	chatId := os.Args[2]
	hookURL := os.Args[3]

	// Read alert data from file
	alertData, err := ioutil.ReadFile(alertFile)
	if err != nil {
		log.Fatalf("failed to read alert file: %v", err)
	}

	// Parse alert JSON
	var alert Alert
	if err := json.Unmarshal(alertData, &alert); err != nil {
		log.Fatalf("failed to parse alert JSON: %v", err)
	}

	// Handle possible types for Rule ID (int or string)
	var ruleID string
	switch id := alert.Rule.ID.(type) {
	case float64:
		ruleID = fmt.Sprintf("%.0f", id) // Convert float64 (JSON numbers) to string
	case string:
		ruleID = id
	default:
		log.Fatalf("unexpected type for rule ID: %T", id)
	}

	// Create message based on alert content
	alertMessage := fmt.Sprintf("Wazuh Alert:\nDescription: %s\nLevel: %d\nRule ID: %s\nAgent: %s (IP: %s)\nTime: %s",
		alert.Rule.Description, alert.Rule.Level, ruleID, alert.Agent.Name, alert.Agent.IP, alert.Timestamp)

	// Send alert message to Green API
	if err := sendToGreenAPI(hookURL, chatId, alertMessage); err != nil {
		log.Fatalf("failed to send alert: %v", err)
	}
}
