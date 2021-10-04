package utils

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// JSONFormatter is a logger for use with Logrus
type JSONFormatter struct {
	Program string
	GinEnv  string
}

// UseJSONLogFormat sets up the JSON log formatter
func UseJSONLogFormat() {

	log.SetFormatter(&JSONFormatter{
		Program: "fibonacci",
	})

	log.SetLevel(log.DebugLevel)
}

// Timestamps in microsecond resolution (like time.RFC3339Nano but microseconds)
var timeStampFormat = "2006-01-02T15:04:05.000000Z07:00"

// Format includes the program, environment, and a custom time format: microsecond resolution
func (f *JSONFormatter) Format(entry *log.Entry) ([]byte, error) {
	data := make(log.Fields, len(entry.Data)+3)
	for k, v := range entry.Data {
		data[k] = v
	}
	data["time"] = entry.Time.UTC().Format(timeStampFormat)
	data["msg"] = entry.Message
	data["level"] = strings.ToUpper(entry.Level.String())
	data["program"] = f.Program

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}

// GetDurationInMilliseconds takes a start time and returns a duration in milliseconds
func GetDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}