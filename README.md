# 2023_MPass
CLI based, simple password manager desktop application for Linux. Provides functionalities such as storing, accessing, adding and modifying sensitive login information. Login entries are stored in an encrypted database. User can create one or multiple databases with unique master passwords, and to gain access to an existing database a correct master password needs to be entered. Encryption was implemented using AES-256 algorithm with GCM, and encryption keys are created using PBKDF2 hashing algorithm.

Project was written in Go language.

## Prerequisites
[Go language](https://go.dev/dl/)
```
$ sudo snap install go --classic
```
  
XSel or XClip\
[Installation guide](https://ostechnix.com/access-clipboard-contents-using-xclip-and-xsel-in-linux/)
## Installation
```
$ git clone https://github.com/matf-pp/2023_MPass.git
```
or just download and extract the zip/tar.gz file.
In terminal type the following commands
 
```
$ cd 2023_MPass
$ chmod +x run.sh
$ ./run.sh
```
Shell script run.sh creates the binary ./main that is in the 2023_MPass/main directory. To continue using the newly created program:
```
$ cd main
$ ./main 
```

The binary is created for Linux based operating systems. 

## Authors
Project created by
Ana Mihajlović ([cholesski](https://github.com/cholesski)) and 
Katarina Grbović ([gkatarina](https://github.com/gkatarina))
