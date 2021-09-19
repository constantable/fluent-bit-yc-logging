build:
	go build -buildmode=c-shared -o fluent-bit-yc-logging.so -v .

clean:
	rm -rf *.so *.h *~

docker-build:
	docker build -f Dockerfile -t fbyl .

docker-run:
	docker run -p 2020:2020 -v test:/test fbyl
