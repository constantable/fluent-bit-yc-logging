package log_entry

import (
	"encoding/json"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestEmptyNewLogEntry(t *testing.T) {
	var ts interface{}
	var record map[interface{}]interface{}
	entry := NewLogEntry(ts, record)

	if entry.Message != "" {
		t.Error(`Message should be empty`)
	}
	if reflect.TypeOf(entry.Timestamp) != reflect.TypeOf((*timestamppb.Timestamp)(nil)) {
		t.Error(`Entry timestamp should be type of timestamppb.Timestamp`)
	}
}
func TestNewLogEntryTimestamp(t *testing.T) {
	var ts output.FLBTime
	var record map[interface{}]interface{}
	ts.Time = time.Now()
	entry := NewLogEntry(ts, record)

	if entry.Message != "" {
		t.Error(`Message should be empty`)
	}
	if reflect.TypeOf(entry.Timestamp) != reflect.TypeOf((*timestamppb.Timestamp)(nil)) {
		t.Error(`Entry timestamp should be type of timestamppb.Timestamp`)
	}
	if time.Unix(int64(entry.Timestamp.Seconds), 0).After(time.Now()) {
		t.Error(`Wrong timestamp`)
	}
}

func TestNewLogEntryMessage(t *testing.T) {
	var ts output.FLBTime
	ts.Time = time.Now()

	record := make(map[interface{}]interface{})
	nested := make(map[interface{}]interface{})
	nested["nested"] = "5"
	record["example"] = "string"
	record["log"] = nested

	entry := NewLogEntry(ts, record)

	if entry.Message != "map[nested:5]" {
		t.Error(`Incorrect message`)
	}
}

var GetLevelTestData = []struct {
	in  []byte
	out logging.LogLevel_Level
}{
	{[]byte(`{"@fields": {"level": "100"}}`), logging.LogLevel_DEBUG},
	{[]byte(`{"@fields": {"level": 200}}`), logging.LogLevel_INFO},
	{[]byte(`{"@fields": {"level": 250}}`), logging.LogLevel_INFO},
	{[]byte(`{"@fields": {"level": "300"}}`), logging.LogLevel_WARN},
	{[]byte(`{"@fields": {"level": "400"}}`), logging.LogLevel_ERROR},
	{[]byte(`{"@fields": {"level": "500"}}`), logging.LogLevel_FATAL},
	{[]byte(`{"@fields": {"level": 550}}`), logging.LogLevel_FATAL},
	{[]byte(`{"@fields": {"level": "600"}}`), logging.LogLevel_FATAL},
	{[]byte(`{"@fields": {"level": "700"}}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"@fields": {"level": 0}}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"@fields": {"field": "200"}}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"fields": {"level": "200"}}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"log": {"level": "200"}}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"log": "something"}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"log": "APP.CRIT Message"}`), logging.LogLevel_FATAL},
	{[]byte(`{"log": {"level": "APP.CRIT Message"}}`), logging.LogLevel_FATAL},
	{[]byte(`{"log": "crit info"}`), logging.LogLevel_LEVEL_UNSPECIFIED},
	{[]byte(`{"log": "crit INFO"}`), logging.LogLevel_INFO},
}

func TestGetLevel(t *testing.T) {
	for _, tt := range GetLevelTestData {
		var m map[string]interface{}
		e := json.Unmarshal(tt.in, &m)
		if e != nil {
			log.Fatal(e)
		}
		level := getLevel(m)
		if level != tt.out {
			t.Errorf("got %q, want %q", level, tt.out)
		}
	}
}
