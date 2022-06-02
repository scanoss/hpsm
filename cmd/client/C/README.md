# Interact with C 
HPSM can be accessed directly from C-based programs using a shared library.
libhpsm.go defines the functionality to be shared. Client.c is just a demo client that demostrate the usage.
## Building the shared library
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go

## Building the client program
gcc -v client.c -o client ./libhpsm.so