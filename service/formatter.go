package golog

import (
	"encoding/json"
	"fmt"
	"time"
)

func JsonFormatHandler(level Level, message *Message) ([]byte, error) {
	addSystemInfo(level, message)
	if bytes, err := json.Marshal(message); err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func TextFormatHandler(level Level, message *Message) ([]byte, error) {
	type MessageText struct {
		prefixes interface{}
		tags     interface{}
		message  interface{}
		fields   interface{}
	}

	addSystemInfo(level, message)
	return []byte(fmt.Sprintf("%+v", MessageText{prefixes: message.Prefixes, tags: message.Tags, message: message.Message, fields: message.Fields})), nil
}

func addSystemInfo(level Level, message *Message) {
	// special prefixes keys
	prefixes := make(map[string]interface{}, len(message.Prefixes))
	for key, value := range message.Prefixes {
		switch value {
		case LEVEL:
			value = level.String()
		case TIME:
			value = time.Now().Format("2006-01-02 15:04:05")
		}
		prefixes[key] = value
	}
	message.Prefixes = prefixes
}
