# HPSM 

This project defines the functionality to do High Precision Snippet Matching between two text based files.
The main principle is based on semi brute force search.

* Each line of both files are normalized and then hashed using CRC8.
* The longest sequence of matching CRCs is calculated doing greedy advance on both files.

The HPSM functionality is available by:
* **Shared library** that provides interface for using the functionallity from a different programing language. There are two flavors:
  *  _Local processing_ Receives Local file hashes and downloads the remote OSS remote file. Then, HPSM is carried out locally considering both set of hashes.
  *  _Remote processing_ Receives Local file hashes and the remote OSS key and send them to a remote gRPC HPSM calculation service. 
  
  Both flavours provide functionallity to hash the content of a file. ("hpsm=01234586787887....")
  
* **CLI**: A simple shell console can be used to:
  * Get the *hpsm=* string to be appended to wfp file
  * Get the *.wfp* file containing the hpsm= string
  * Compare two files in a graphical view. The first file must be local but the second file can be local or a remote MD5
* **go module** A go package (Coming soon). 

## Building
 A makefile is provided to automate the building process

  ``make local_proc`` 
  
  Creates a shared lirbrary (*libhpsm.so* or *libhpsm.dll*) and header file (*libhpsm.h*). This library is intended to be used for remote download of OSS files and local HPSM calculation. The created library should be placed on a SO library folder (eg: *libhpsm.so* must be placed on **/usr/lib** to make the feature available on linux) *libhpsm.h* should be used by the client application, eg: **inc/** folder. 
  By default _osskb.org_ url is used to download the OSS files. It can be modified by using SRC_URL enviromental variable.
  
---
``make remote_proc`` 
  Creates a shared lirbrary (*libhpsm.so* or *libhpsm.dll*) and header file (*libhpsm.h*). This library is intended to be used for remote processing of HPSM calculation by using a gRPC service. The created library should be placed on a SO library folder (eg: *libhpsm.so* must be placed on **/usr/lib** to make the feature available on linux) *libhpsm.h* should be used by the client application, eg: **inc/** folder. 
  By default _osskb.org_ url is used to do HPSM calculation. It can be modified by using HPSM_URL enviromental variable.

  The HPSM service is also created, ready to be deployed. By running the service without arguments you will receive default service:
  Port: 51015
	Workers: 2
	Threshold: 5 lines

  A new configuration can be used by setting the configuration in a JSON file and sent to the service: 
  Eg: conf.json:

     {
        "port":50050,
        "workers": 5,
         "threshold": 4
     }

To run the service use:
  
  `` hpsm-service conf.json``

---
  ``make cli`` 
  
  creates *hpsm* cli application that let you do file side-by-side comparisson between two local files or between a local file and a remote file given its MD5 key
  
---
  ``make install`` 
  
  copy both files *libhpsm.so* and *hpsm* to their destination to make the feature work (**/usr/lib** and **/usr/bin/**)
  
   ``make clean`` 
  
  removes compiled binaries *libhpsm.so* and *hpsm*
  

  
 
&copy; SCANOSS 2018-2023
