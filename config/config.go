package config

import (
	"encoding/base64"
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

	publicKey, err := base64Decode(output.FLBPluginConfigKey(ctx, "PublicKey"))
	if err != nil {
		return
	}
	privateKey, err := base64Decode(output.FLBPluginConfigKey(ctx, "PrivateKey"))
	if err != nil {
		return
	}
	cfg.ServiceAccountKey = ServiceAccountKey{
		Id:         output.FLBPluginConfigKey(ctx, "KeyId"),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	return
}

func base64Decode(str string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), err
}
