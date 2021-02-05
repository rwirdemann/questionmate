build:
	go build -o bin/questionmate world/http/main.go

build-linux:
	env GOOS=linux CGO_ENABLED=0 go build ${LDFLAGS} -a -installsuffix cgo -o bin/questionmate world/http/main.go

deploy: build
	ssh 95.217.222.60 "pkill questionmate"
	scp bin/questionmate 95.217.222.60:~/questionmate/bin
	scp -r config/* 95.217.222.60:~/questionmate/config
	ssh 95.217.222.60 "sh -c 'nohup /home/ralf/questionmate/bin/questionmate -directory=/home/ralf/questionmate/config/coma > /dev/null 2>&1 &'"

test:
	env go test -count=1 ./...

test-all:
	env go test -count=1 ./... -tags=integration

clean:
	rm -rf bin

.PHONY: test
