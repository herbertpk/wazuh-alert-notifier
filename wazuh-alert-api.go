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

// MessagePayload represents the payload structure for the Green API
type MessagePayload struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

// Configuration holds Green API configuration details
type Configuration struct {
	APIURL          string
	IDInstance      string
	APITokenInstance string
	ChatID          string
}

// Load configuration from environment variables
func loadConfig() Configuration {
	return Configuration{
		APIURL:          os.Getenv("GREEN_API_URL"),
		IDInstance:      os.Getenv("GREEN_API_INSTANCE_ID"),
		APITokenInstance: os.Getenv("GREEN_API_TOKEN"),
		ChatID:          os.Getenv("GREEN_API_CHAT_ID"),
	}
}

// sendAlertToAPI sends the alert message to the Green API
func sendAlertToAPI(config Configuration, message string) error {
	// Construct the full URL
	url := fmt.Sprintf("%s/waInstance%s/sendMessage/%s", config.APIURL, config.IDInstance, config.APITokenInstance)

	// Create the payload
	payload := MessagePayload{
		ChatID:  config.ChatID,
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

// alertHandler is an HTTP handler that processes incoming alerts and sends them to the Green API
func alertHandler(w http.ResponseWriter, r *http.Request) {
	config := loadConfig()

	// Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the alert data from the request body
	var alert Alert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, "Failed to parse alert data", http.StatusBadRequest)
		log.Printf("Failed to decode alert: %v", err)
		return
	}

	// Create a message based on the alert
	alertMessage := fmt.Sprintf("Wazuh Alert:\nDescription: %s\nLevel: %d\nAgent: %s (IP: %s)\nTime: %s",
		alert.Rule.Description, alert.Rule.Level, alert.Agent.Name, alert.Agent.IP, alert.Timestamp)

	// Send the alert message to the Green API
	if err := sendAlertToAPI(config, alertMessage); err != nil {
		http.Error(w, "Failed to send alert to Green API", http.StatusInternalServerError)
		log.Printf("Failed to send alert: %v", err)
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Alert sent successfully")
	log.Println("Alert processed and sent to Green API")
}

func main() {
	// Set up the HTTP server
	http.HandleFunc("/alert", alertHandler)
	port := "8083"
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

