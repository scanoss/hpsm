# HPSM 

This project defines the functionality to do High Precision Snippet Matching on any file.
The main principle is based on semi brute force search.

* Each line of the file is normalized and then hashed using CRC8.
* The file with the given MD5 is also hashed.
* The longest sequence of CRCs is calculated doing greedy advance on local and remote file.

The functionality is available by:

* **API**: an endpoint receives a json structure defining a set of <md5><[hashes]> to be processed. The API could be deployed on the sources server or can download the sources from several servers.
* **libhpsm**: A shared library that provides local processing (by downloading from ossk.org) or remote processing (calling the above mentioned API). It also provides functionallity to hash the content of a file. ("hpsm=01234586787887....")
* **go module** A go package (Coming soon). 

&copy; SCANOSS 2018-2022