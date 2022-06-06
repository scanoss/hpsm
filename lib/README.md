# Interact with C 
HPSM can be accessed directly from C-based programs using a shared library.
libhpsm.go defines the functionality to be shared. Client.c is just a demo client that demostrate the usage.
## Building the shared library
```
go build -o libhpsm.so  -buildmode=c-shared libhpsm.go
```
The command will generate two files:
- **libhpsm.so** This is the shared library. Carry on the functionallity related of HPSM. Should be placed in ***/usr/lib*** folder
- **libhpsm.h** The header file that defines the prototypes of functions provided by libhpsm.so library. Should be copied in your ***inc/*** folder.

## Building the demo client program
The demo client is just a basic example of how to instantiate the library in your source code
To create an executable type:
```
gcc -v client.c -o client ./libhpsm.so
```
## Including the library on your project
In order use the library inside your C Project, you should copy the ***libhpsm.h*** to any folder containing header files.

You should also add the parameter ***-hpsm*** in your building phase of the Makefile
