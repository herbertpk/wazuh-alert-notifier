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
		Description string `json:"description"`
		Level       int    `json:"level"`
	} `json:"rule"`
	Agent struct {
		Name string `json:"name"`
		IP   string `json:"ip"`
	} `json:"agent"`
	Timestamp string `json:"timestamp"`
}

// MessagePayload represents the payload structure for the API
type MessagePayload struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

// sendAlertToAPI sends the alert message to the WhatsApp API
func sendAlertToAPI(apiURL, idInstance, apiTokenInstance, chatID, message string) error {
	// Construct the full URL
	url := fmt.Sprintf("%s/waInstance%s/sendMessage/%s", apiURL, idInstance, apiTokenInstance)

	// Create the payload
	payload := MessagePayload{
		ChatID:  chatID,
		Message: message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	
	// Set up the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-200 response and log response body
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("non-200 response from API: %d, response body: %s", resp.StatusCode, string(body))
	}

	log.Printf("Alert sent successfully with response code %d\n", resp.StatusCode)
	return nil
}

func main() {
	// Check if all required arguments are provided
	if len(os.Args) < 5 {
		log.Fatal("Usage: ./notify-alert <apiUrl> <idInstance> <apiTokenInstance> <chatId>")
	}

	// Read arguments
	apiURL := os.Args[1]
	idInstance := os.Args[2]
	apiTokenInstance := os.Args[3]
	chatID := os.Args[4]

	// Read the alert data from standard input
	var alert Alert
	if err := json.NewDecoder(os.Stdin).Decode(&alert); err != nil {
		log.Fatalf("failed to decode alert: %v", err)
	}

	// Create a message based on the alert
	alertMessage := fmt.Sprintf("Wazuh Alert:\nDescription: %s\nLevel: %d\nAgent: %s (IP: %s)\nTime: %s",
		alert.Rule.Description, alert.Rule.Level, alert.Agent.Name, alert.Agent.IP, alert.Timestamp)

	// Send the alert message to the API
	if err := sendAlertToAPI(apiURL, idInstance, apiTokenInstance, chatID, alertMessage); err != nil {
		log.Fatalf("failed to send alert: %v", err)
	}
}
