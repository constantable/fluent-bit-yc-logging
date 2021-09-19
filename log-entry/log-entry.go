package log_entry

import (
	"fmt"
	"github.com/fluent/fluent-bit-go/output"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
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
	msg := make(map[string]interface{})
	msg = parseRecord(record)

	payload, err := structpb.NewStruct(msg)
	if err != nil {
		log.Fatal(err)
	}
	return &logging.IncomingLogEntry{
		Timestamp:   timestamppb.New(timestampTime),
		Message:     fmt.Sprintf("%v", msg["log"]),
		JsonPayload: payload,
	}
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
