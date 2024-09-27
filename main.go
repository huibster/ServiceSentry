package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Response struct {
	OK          bool    `json:"ok"`
	Status      int     `json:"status"`
	Name        string  `json:"name"`
	ClusterName string  `json:"cluster_name"`
	Version     Version `json:"version"`
}

type Version struct {
	Number         string `json:"number"`
	BuildHash      string `json:"build_hash"`
	BuildTimestamp string `json:"build_timestamp"`
	BuildSnapshot  bool   `json:"build_snapshot"`
	LuceneVersion  string `json:"lucene_version"`
}

func getJsonData(url string) (*Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResponse Response
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &apiResponse, nil
}

func writeToLog(errMessage string) {
	dir := "/var/log/servicesentry"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
		fmt.Printf("Directory %s created successfully\n", dir)
	} else {
		fmt.Printf("Directory %s already exists\n", dir)
	}

	file, err := os.OpenFile("/var/log/servicesentry/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.Printf("ERROR: %s", errMessage)
}

func main() {
	apiUrl := "http://localhost:4200/"
	serviceCheck := false
	jsonData, err := getJsonData(apiUrl)
	if err == nil {
		if jsonData.Status == 200 {
			serviceCheck = true
		}
	}
	if serviceCheck {
		fmt.Println("Service ok")
	} else {
		fmt.Println("Service down")
		writeToLog("CrateDB Node down")
		cmd := exec.Command("systemctl", "restart", "crate")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			writeToLog("Failed to restart CrateDB Node")
			// log.Fatalf("Failed to restart CrateDB: %v", err)
		}
	}
	os.Exit(0)
}

// sudo systemctl restart crate
// cd /var/www && forever stop server.js && forever start server.js prod
// cmd := exec.Command("/bin/sh", "-c", "cd /var/www && forever stop server.js && forever start server.js prod")
// cmd := exec.Command("sudo", "systemctl", "restart", "crate")

/*
	message := "Planetnine: CrateDB Node down"
	toPhone := "+48668592828"   // The recipient's phone number
	fromPhone := "+13158884479" // Your Twilio phone number
*/
// Function to send an SMS using Twilio
/*
func sendSMS(message string, to string, from string) error {
	// Get Twilio account credentials from environment variables
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(from)
	params.SetBody(message)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}
*/
