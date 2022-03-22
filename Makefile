s?=s
compile:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o  app -work $(shell pwd)/

run:
	go build -o app
	GO_MICRO_ENV=local ./app --config $(shell pwd)/dev.toml
nsq:
	docker-compose up
send:
	curl -d '{"val":"hello"}' 'http://127.0.0.1:4151/pub?topic=xxx'
csend:
	while true;do curl -d '{"val":"hello"}' 'http://127.0.0.1:4151/pub?topic=xxx';done;