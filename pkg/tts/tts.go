package tts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TTSRequest struct {
	Text    string `json:"text"`
	Speaker string `json:"speaker"`
}

func TextToSpeech(text, speaker string) ([]byte, error) {
	ttsReq := TTSRequest{
		Text:    text,
		Speaker: speaker,
	}

	jsonData, err := json.Marshal(ttsReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	resp, err := http.Post("http://localhost:5000/tts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return audioData, nil
}
