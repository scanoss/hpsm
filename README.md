# HPSM 

This project defines the functionality to do High Precision Snippet Matching on any file.
The main principle is based on semi brute force search.

* Each line of the file is normalized and then hashed using CRC8.
* The file with the given MD5 is also hashed.
* The longest sequence of CRCs is calculated doing greedy advance on local and remote file.

The functionality is available by:

* **API**: an endpoint receives a json structure defining a set of <md5><[hashes]> to be processed. The API could be deployed on the sources server or can download the sources from several servers.
* **libhpsm**: A shared library that provides local processing (by downloading from ossk.org) or remote processing (calling the above mentioned API). It also provides functionallity to hash the content of a file. ("hpsm=01234586787887....")
* **CLI**: A simple shell console can be used to:
  * Get the *hpsm=* string to be appended to wfp file
  * Get the *.wfp* file containing the hpsm= string
  * Compare two files in a graphical view. The first file must be local but the second file can be local or a remote MD5
* **go module** A go package (Coming soon).

## Remote Files
By default, files from *"https://osskb.org/api/file_contents/"* are retrieved. If other sources server is used, the environmental variable **SCANOSS_FILE_CONTENTS_URL** must be used. Eg:
  
 ``export  SCANOSS_FILE_CONTENTS_URL=https://osskb.org/api/file_contents/``

You can also set the environmental variable `SCANOSS_API_KEY` that will be passed as the `X-Session` header.

 ``export  SCANOSS_API_KEY=[YOUR_KEY]``


## Building
A makefile is provided to automate the building process

``make build_lib`` 

Creates *libhpsm.so* and *libhpsm.h*. *libhpsm.so* must be placed on **/usr/lib** to make the feature available. *libhpsm.h* should be used by the client application, eg: **inc/** folder


``make build_cli`` 

Creates *hpsm* cli application

  
``make build_amd`` 

Builds and packages *hpsm* lib and cli application. Creates temporary directory **scripts** with **hpsm** (binary), **libhpsm.so** (lib), and **libhpsm.h** (header file).

``make install`` 

Copy both files *libhpsm.so* and *hpsm* to their destination to make the feature work (**/usr/lib** and **/usr/bin/**)

  ``make clean`` 

Removes compiled binaries *libhpsm.so* and *hpsm*





&copy; SCANOSS 2018-2022
