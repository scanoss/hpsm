# Interact from your source code (Local processing)
HPSM can be accessed directly from a different programming language (for example: C) by using a shared library.
libhpsm.go wraps the functionality to be shared. You will find a demo source code (client.c) that demostrates the usage.
Local processing does remote downloading of the OSS file and compares against the local hashes. To compare both sets of hashes use the exported function **HashFileContents**
**Local hashes** can be created by using the exporte function **HashFileContents**
## Building the shared library
From the root of the project type: 
```
make local_proc
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

You should also add the parameter ***-hpsm*** in the libs section of your building phase of the Makefile

## Set up
By default, OSS source codes are downloaded from **osskb.org**, but it can be modified by setting up the Enviroment variable **SRC_URL** where **SRC_URL** points to a server running the Scanoss API. (< server url >/api/file_contents/)