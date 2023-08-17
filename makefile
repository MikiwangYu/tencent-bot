all: bin/randbot.exe

bin/randbot.exe: cmd/randbot.go
	go build -o $@ $^
