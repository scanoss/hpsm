VERSION=$(shell ./version.sh)

BINARY_NAME=hpsm
LIB_NAME=libhpsm.so

clean: 
	rm hpsm
	rm libhpsm.so
	
version:  ## Produce Semgrep version text file
	@echo "Writing version file..."
	echo $(VERSION) > version.txt

unit_test:  ## Run all unit tests in the pkg folder
	@echo "Running unit test framework..."
	go test -v ./test/...

build_lib:
	go build -o ${LIB_NAME}  -buildmode=c-shared lib/libhpsm.go

build_cli:
	go build -o ${BINARY_NAME} main.go

install:
	cp ${LIB_NAME} /usr/lib
	cp ${BINARY_NAME} /usr/bin

build_amd: version  ## Build an AMD 64 binary
	@echo "Building AMD binary $(VERSION)..."
	go build -o ${LIB_NAME}  -buildmode=c-shared lib/libhpsm.go

package: package_amd  ## Build & Package an AMD 64 binary

package_amd: version  ## Build & Package an AMD 64 binary
	@echo "Building AMD binary $(VERSION) and placing into scripts..."
	mkdir -p scripts
	go build -o scripts/${LIB_NAME}  -buildmode=c-shared lib/libhpsm.go
	go build -o scripts/${BINARY_NAME} main.go
	bash ./package-scripts.sh linux-amd64 $(VERSION)