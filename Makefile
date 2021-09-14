build:
	go build -buildmode=c-shared -o fluent-bit-yc-logging.so .

clean:
	rm -rf *.so *.h *~

proto:
	protoc  -Icloudapi -Icloudapi/third_party/googleapis --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    yandex/cloud/logging/v1/log_ingestion_service.proto

clone-api-proto:
	git clone https://github.com/yandex-cloud/cloudapi.git 
