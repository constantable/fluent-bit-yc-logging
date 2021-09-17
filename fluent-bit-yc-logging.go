package main

import (
	"C"
	"time"

	"log"
	"unsafe"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	_ "github.com/yandex-cloud/go-sdk/iamkey"

	"fluent-bit-yc-logging/config"
	"fluent-bit-yc-logging/connection"
	"fmt"
)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	// Gets called only once when the plugin.so is loaded
	return output.FLBPluginRegister(def, "fluent-bit-yc-logging", "Fluent Bit Yandex Cloud Logging output plugin")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	cfg, err := config.NewConfig(plugin)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		return output.FLB_ERROR
	}

	client, err := connection.New(cfg)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	output.FLBPluginSetContext(plugin, client)
	return output.FLB_OK
}

//export FLBPluginFlushCtx
func FLBPluginFlushCtx(ctx, data unsafe.Pointer, length C.int, tag *C.char) int {

	dec := output.NewDecoder(data, int(length))
	var entries []*logging.IncomingLogEntry
	for {
		ret, ts, record := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		var timestampTime time.Time
		switch tts := ts.(type) {
		case output.FLBTime:
			timestampTime = tts.Time
		case uint64:
			timestampTime = time.Unix(int64(tts), 0)
		default:
			timestampTime = time.Now()
		}
		msg := make(map[string]string)
		for key, value := range record {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			msg[strKey] = strValue
		}
		//out.Write(record, timestamp, C.GoString(tag))
		entries = append(entries, &logging.IncomingLogEntry{
			Timestamp: timestamppb.New(timestampTime),
			Message:   fmt.Sprintf("%#v", msg),
		})
	}

	client := output.FLBPluginGetContext(ctx).(connection.Client)
	err := connection.WriteEntries(client, entries)
	if err != nil {
		log.Fatal(err)
		return output.FLB_ERROR
	}
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {

}
