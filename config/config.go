package config

import (
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
)

type PluginConfig struct {
	ServiceAccountId  string
	ServiceAccountKey ServiceAccountKey
	LogGroupId        string
}
type ServiceAccountKey struct {
	Id         string
	PublicKey  string
	PrivateKey string
}

func NewConfig(ctx unsafe.Pointer) (cfg PluginConfig, err error) {
	cfg.ServiceAccountId = output.FLBPluginConfigKey(ctx, "ServiceAccountId")
	cfg.LogGroupId = output.FLBPluginConfigKey(ctx, "LogGroupId")
	cfg.ServiceAccountKey = ServiceAccountKey{
		Id:         output.FLBPluginConfigKey(ctx, "KeyId"),
		PublicKey:  output.FLBPluginConfigKey(ctx, "PublicKey"),
		PrivateKey: output.FLBPluginConfigKey(ctx, "PrivateKey"),
	}
	return
}
