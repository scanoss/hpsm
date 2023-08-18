BINARY_NAME=hpsm
LIB_NAME=libhpsm.so
SRV_NAME=hpsm-service

local_proc:
	go build -o ${LIB_NAME}  -buildmode=c-shared lib/local/libhpsm.go

remote_proc:
	go build -o ${LIB_NAME}  -buildmode=c-shared lib/remote/libhpsm.go
	go build -o ${SRV_NAME} cmd/server/grpc/server.go


cli:
	go build -o ${BINARY_NAME} main.go

server:
	go build -o ${SRV_NAME} cmd/server/grpc/server.go

clean: 
	rm hpsm
	rm libhpsm.so

install:
	cp ${LIB_NAME} /usr/lib
	cp ${BINARY_NAME} /usr/bin
