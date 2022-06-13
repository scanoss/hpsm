BINARY_NAME=hpsm
LIB_NAME=libhpsm.so

build_lib:
	go build -o ${LIB_NAME}  -buildmode=c-shared lib/libhpsm.go

cli:
	go build -o ${BINARY_NAME} main.go

clean: 
	rm hpsm
	rm libhpsm.so

install:
	cp ${LIB_NAME} /usr/lib
	cp ${BINARY_NAME} /usr/bin
