package connection

import (
	"context"
	"fluent-bit-yc-logging/config"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/logging/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/iamkey"
)

type Client struct {
	ctx        context.Context
	sdk        *ycsdk.SDK
	LogGroupId string
}

func New(cfg config.PluginConfig) (client Client, err error) {
	client.ctx = context.Background()

	var credentials ycsdk.Credentials
	credentials, err = ycsdk.ServiceAccountKey(&iamkey.Key{
		Id:         cfg.ServiceAccountKey.Id,
		Subject:    &iamkey.Key_ServiceAccountId{ServiceAccountId: cfg.ServiceAccountId},
		PublicKey:  cfg.ServiceAccountKey.PublicKey,
		PrivateKey: cfg.ServiceAccountKey.PrivateKey,
	})

	client.LogGroupId = cfg.LogGroupId

	if err != nil {
		return
	}

	client.sdk, err = ycsdk.Build(client.ctx, ycsdk.Config{
		Credentials: credentials,
	})
	return
}

func WriteEntries(client Client, entries []*logging.IncomingLogEntry) (err error) {
	request := &logging.WriteRequest{
		Destination: &logging.Destination{
			Destination: &logging.Destination_LogGroupId{
				LogGroupId: client.LogGroupId,
			},
		},
		Entries: entries,
	}
	_, err = client.sdk.LogIngestion().LogIngestion().Write(client.ctx, request)
	return
}
