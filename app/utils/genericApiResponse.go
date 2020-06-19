package utils

import (
	"fmt"
	"net/http"
	"time"
)

// Message in json object
func Message(status string, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// APIResponse in json format
func APIResponse(w http.ResponseWriter, valIn []byte) {
	const (
		layout = "2006-01-02T15:04:05-0700"
	)
	t := time.Now()
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"status":"OK","lastfetched":%s, "data": %s}`, t.Format(layout), valIn)))
}
