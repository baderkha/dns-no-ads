start:
	go run .\cmd\client\main.go
start-unix:
	sudo /usr/local/go/bin/go run ./cmd/client/main.go
pull-block-list:
	CGO_ENABLED=1 go run .\scripts\dns\blocklist\main.go