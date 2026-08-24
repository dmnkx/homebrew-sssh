package debuglog

import (
	"encoding/json"
	"os"
	"time"
)

const path = "/Users/cookat/workspaces/golang/.cursor/debug-7329b2.log"

func Log(hypothesisID, location, message string, data map[string]any) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	payload := map[string]any{
		"sessionId":    "7329b2",
		"hypothesisId": hypothesisID,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, _ = f.Write(append(b, '\n'))
}
