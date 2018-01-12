all: talker listener

talker:
	go build ./cmd/talker

listener:
	go build ./cmd/listener
