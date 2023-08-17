# HPSM 

This project defines the functionality to do High Precision Snippet Matching between two text based files.
The main principle is based on semi brute force search.

* Each line of the both files are normalized and then hashed using CRC8.
* The longest sequence of matching CRCs is calculated doing greedy advance on both files.

The functionality is available by:
* **gRPC Service**: an gRPC resource running on server sources receives an array of hashes and calculates the hashes of a given MD5 key. Then, the calculation is carried out and returns the list of matching ranges
* **libhpsm**: A shared library that provides local hashing and calls remote processing on the gRPC service. It also provides functionallity to hash the content of a file. ("hpsm=01234586787887....")
* **CLI**: A simple shell console can be used to:
  * Get the *hpsm=* string to be appended to wfp file
  * Get the *.wfp* file containing the hpsm= string
  * Compare two files in a graphical view. The first file must be local but the second file can be local or a remote MD5
* **go module** A go package (Coming soon). 

  ## Building
 A makefile is provided to automate the building process

  ``make build_lib`` 
  
  creates *libhpsm.so* and *libhpsm.h*. *libhpsm.so* must be placed on **/usr/lib** to make the feature available. *libhpsm.h* should be used by the client application, eg: **inc/** folder
  
  ``make server`` 
  
  generate hpsm-service ready to be run on a server that hosts mz files

  ``make cli`` 
  
  creates *hpsm* cli application
  
  
  ``make install`` 
  
  copy both files *libhpsm.so* and *hpsm* to their destination to make the feature work (**/usr/lib** and **/usr/bin/**)
  
   ``make clean`` 
  
  removes compiled binaries *libhpsm.so* and *hpsm*
  

  
 
&copy; SCANOSS 2018-2023
