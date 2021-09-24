package log_entry

import (
	"fmt"
	"log"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewLogEntry(ts interface{}, record map[interface{}]interface{}) (entry *logging.IncomingLogEntry) {
	var timestampTime time.Time
	switch tts := ts.(type) {
	case output.FLBTime:
		timestampTime = tts.Time
	case uint64:
		timestampTime = time.Unix(int64(tts), 0)
	default:
		timestampTime = time.Now()
	}

	msg := parseRecord(record)

	payload, err := structpb.NewStruct(msg)
	if err != nil {
		log.Fatal(err)
	}
	var message interface{}
	ok := false
	if message, ok = msg["log"]; !ok {
		if message, ok = msg["@message"]; !ok {
			message = ""
		}
	}

	return &logging.IncomingLogEntry{
		Timestamp:   timestamppb.New(timestampTime),
		Message:     fmt.Sprintf("%v", message),
		Level:       getLevel(msg),
		JsonPayload: payload,
	}
}

func getLevel(msg map[string]interface{}) logging.LogLevel_Level {
	level := logging.LogLevel_LEVEL_UNSPECIFIED

	// if logstash formatted
	if fields, ok := msg["@fields"]; ok {
		if lvl, ok2 := fields.(map[string]interface{})["level"]; ok2 {
			switch fmt.Sprintf("%v", lvl) {
			case "100":
				return logging.LogLevel_DEBUG
			case "200":
				return logging.LogLevel_INFO
			case "250": // NOTICE generally
				return logging.LogLevel_INFO
			case "300":
				return logging.LogLevel_WARN
			case "400":
				return logging.LogLevel_ERROR
			case "500": // CRITICAL
				return logging.LogLevel_FATAL
			case "550": // ALERT
				return logging.LogLevel_FATAL
			case "600": // EMERGENCY
				return logging.LogLevel_FATAL
			default:
				return logging.LogLevel_LEVEL_UNSPECIFIED
			}
		}
	}

	return level
}

func parseRecord(inputRecord map[interface{}]interface{}) map[string]interface{} {
	return parseValue(inputRecord).(map[string]interface{})
}

func parseValue(value interface{}) interface{} {
	switch value := value.(type) {
	case []byte:
		return string(value)
	case map[interface{}]interface{}:
		remapped := make(map[string]interface{})
		for k, v := range value {
			remapped[k.(string)] = parseValue(v)
		}
		return remapped
	case []interface{}:
		remapped := make([]interface{}, len(value))
		for i, v := range value {
			remapped[i] = parseValue(v)
		}
		return remapped
	default:
		return value
	}
}
