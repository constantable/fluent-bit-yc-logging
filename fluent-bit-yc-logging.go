package main

import (
	"C"
	"context"
	"time"

	"github.com/yandex-cloud/go-sdk/iamkey"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	_ "github.com/yandex-cloud/go-sdk/iamkey"
)

//export FLBPluginRegister
func FLBPluginRegister(def unsafe.Pointer) int {
	// Gets called only once when the plugin.so is loaded
	return output.FLBPluginRegister(def, "fluent-bit-yc-logging", "Fluent Bit Yandex Cloud Logging output plugin")
}

//export FLBPluginInit
func FLBPluginInit(plugin unsafe.Pointer) int {
	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	//ctx := context.Background()
	//
	//var credentials ycsdk.Credentials
	//credentials, err := ycsdk.ServiceAccountKey( &iamkey.Key	{
	//	Id: "",
	//	Subject: &iamkey.Key_ServiceAccountId{ServiceAccountId: ""},
	//	PublicKey: "",
	//	PrivateKey: "",
	//})
	//if err != nil {
	//	panic(err.Error())
	//}
	//
	//sdk, err := ycsdk.Build(ctx, ycsdk.Config{
	//	Credentials: credentials,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//dec := output.NewDecoder(data, int(length))
	//var entries []*logging.IncomingLogEntry
	//for {
	//	ret, ts, record := output.GetRecord(dec)
	//	if ret != 0 {
	//		break
	//	}
	//
	//	var timestampTime time.Time
	//	switch tts := ts.(type) {
	//	case output.FLBTime:
	//		timestampTime = tts.Time
	//	case uint64:
	//		// From our observation, when ts is of type uint64 it appears to
	//		// be the amount of seconds since unix epoch.
	//		timestampTime = time.Unix(int64(tts), 0)
	//	default:
	//		timestampTime = time.Now()
	//	}
	//
	//	//out.Write(record, timestamp, C.GoString(tag))
	//	//entries = append(entries, &logging.IncomingLogEntry{
	//	//	Timestamp: timestamppb.New(timestampTime),
	//	//	Message: *record,
	//	//
	//	//})
	//}
	//
	//request := &logging.WriteRequest{
	//	Destination: &logging.Destination{
	//		Destination: &logging.Destination_LogGroupId{
	//			LogGroupId: "e23f25rgeoguvimpdlad",
	//		},
	//	},
	//	Entries: entries,
	//}
	//_, err = sdk.LogIngestion().LogIngestion().Write(ctx, request)
	//
	//if err != nil {
	//	log.Fatal(err)
	//	return output.FLB_ERROR
	//}
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	return output.FLB_OK
}

func main() {
	ctx := context.Background()

	var credentials ycsdk.Credentials
	credentials, err := ycsdk.ServiceAccountKey(&iamkey.Key{
		Id:         "",
		Subject:    &iamkey.Key_ServiceAccountId{ServiceAccountId: ""},
		PublicKey:  "",
		PrivateKey: "",
	})
	if err != nil {
		panic(err.Error())
	}

	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: credentials,
	})
	if err != nil {
		log.Fatal(err)
	}
	request := &logging.WriteRequest{
		Destination: &logging.Destination{
			Destination: &logging.Destination_LogGroupId{
				LogGroupId: "e23b9okkahsg3qak8h6m",
			},
		},
		Entries: []*logging.IncomingLogEntry{
			{
				Message:   "Hello world",
				Timestamp: timestamppb.New(time.Now()),
			},
		},
	}
	response, err := sdk.LogIngestion().LogIngestion().Write(ctx, request)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Response errors: %s", response.GetErrors()[0].GetMessage())
}
