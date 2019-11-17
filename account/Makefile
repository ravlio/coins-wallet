gen:
	protoc  -I=$(GOPATH)/src:pb --gogofast_out=plugins=grpc:pb pb/*.proto
	easyjson --all pb/messages.pb.go service.go

tests:
	go test -v ./...

all:
	gen tests
