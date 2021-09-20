package main

import (
	"C"
	log_entry "fluent-bit-yc-logging/log-entry"
	"log"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	_ "github.com/yandex-cloud/go-sdk/iamkey"

	"fluent-bit-yc-logging/config"
	"fluent-bit-yc-logging/connection"
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
		log.Printf("[ERROR fluent-bit-yc-logging] %v", err)
		return output.FLB_ERROR
	}
	client, err := connection.New(cfg)
	if err != nil {
		log.Printf("[ERROR fluent-bit-yc-logging] %v", err)
		return output.FLB_ERROR
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

		entry := log_entry.NewLogEntry(ts, record)
		entries = append(entries, entry)

		if len(entries) == 100 {
			err := flush(ctx, entries)
			entries = nil
			if err != nil {
				log.Fatal(err)
				return output.FLB_ERROR
			}
		}
	}

	err := flush(ctx, entries)
	if err != nil {
		log.Fatal(err)
		return output.FLB_ERROR
	}
	return output.FLB_OK
}

func flush(ctx unsafe.Pointer, entries []*logging.IncomingLogEntry) (err error) {
	client := output.FLBPluginGetContext(ctx).(connection.Client)
	err = connection.WriteEntries(client, entries)
	return
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {

}
