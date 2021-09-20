package log_entry

import (
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

	if entry.Message != "<nil>" {
		t.Error(`Message should be <nil>`)
	}
	if reflect.TypeOf(entry.Timestamp) != reflect.TypeOf((*timestamppb.Timestamp)(nil)) {
		t.Error(`Entry timestamp shoud be type of timestamppb.Timestamp`)
	}
}
func TestNewLogEntryTimestamp(t *testing.T) {
	var ts output.FLBTime
	var record map[interface{}]interface{}
	ts.Time = time.Now()
	entry := NewLogEntry(ts, record)

	if entry.Message != "<nil>" {
		t.Error(`Message should be <nil>`)
	}
	if reflect.TypeOf(entry.Timestamp) != reflect.TypeOf((*timestamppb.Timestamp)(nil)) {
		t.Error(`Entry timestamp shoud be type of timestamppb.Timestamp`)
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
