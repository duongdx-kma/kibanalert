package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"kibanalert/alerts"
	"net/http"
	"os"
)

type TeamsMessage struct {
	Type        string        `json:"type"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	ContentType string      `json:"contentType"`
	Content     AdaptiveCard `json:"content"`
}

type AdaptiveCard struct {
	Schema    string      `json:"$schema"`
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	MSTeams   MSTeams     `json:"msteams"`
	Body      []TextBlock `json:"body"`
}

type MSTeams struct {
	Width string `json:"width"`
}

type TextBlock struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Weight string `json:"weight,omitempty"`
	Size  string `json:"size,omitempty"`
	Wrap  bool   `json:"wrap"`
}

// SendToTeamsAlert gửi cảnh báo đến Microsoft Teams Webhook
func SendToTeamsAlert(source alerts.Source) error {
	fmt.Println("DEBUG: Starting SendToTeamsAlert function")
	fmt.Printf("DEBUG: Received alert: %+v\n", source)

	webhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if webhookURL == "" {
		fmt.Println("ERROR: Teams webhook URL is not set")
		return errors.New("Teams webhook URL is not set")
	}

	message := TeamsMessage{
		Type: "message",
		Attachments: []Attachment{
			{
				ContentType: "application/vnd.microsoft.card.adaptive",
				Content: AdaptiveCard{
					Schema:  "http://adaptivecards.io/schemas/adaptive-card.json",
					Type:    "AdaptiveCard",
					Version: "1.4",
					MSTeams: MSTeams{Width: "Full"},
					Body: []TextBlock{
						{Type: "TextBlock", Text: fmt.Sprintf("ㆍ Service Name: %s", source.ServiceName), Wrap: true},
						{Type: "TextBlock", Text: fmt.Sprintf("ㆍ Reason: %s", source.Reason), Wrap: true},
						{Type: "TextBlock", Text: fmt.Sprintf("ㆍ Date: %s", source.Date), Wrap: true},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("ERROR: Failed to marshal JSON:", err)
		return err
	}
	fmt.Println("DEBUG: JSON payload prepared:", string(jsonData))

	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("ERROR: Failed to create HTTP request:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: Failed to send request to Teams webhook:", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Println("DEBUG: Response status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ERROR: Teams webhook error: received status code %d\n", resp.StatusCode)
		return fmt.Errorf("Teams webhook error: received status code %d", resp.StatusCode)
	}

	fmt.Println("DEBUG: Message sent successfully to Teams")
	return nil
}
