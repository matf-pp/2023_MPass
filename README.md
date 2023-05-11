# 2023_MPass

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/13da5432f0824aa984b9550909697435)](https://app.codacy.com/gh/matf-pp/2023_MPass/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)

CLI based, simple password manager desktop application for Linux. Provides functionalities such as storing, accessing, adding and modifying sensitive login information. Login entries are stored in an encrypted database. User can create one or multiple vaults with unique master passwords, and to gain access to an existing vault a correct master password needs to be entered. Encryption was implemented using AES-256 algorithm with GCM, and encryption keys are created using PBKDF2 hashing algorithm.

Project was written in Go language.

## Prerequisites
[Go language](https://go.dev/dl/)
```
$ sudo snap install go --classic
```
  
XSel or XClip\
[Installation guide](https://ostechnix.com/access-clipboard-contents-using-xclip-and-xsel-in-linux/)
## Installation
1. Building your own binaries: 
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
    Shell script run.sh creates the binaries ./main and ./tui_main that are in the 2023_MPass/main, 2023_MPass/tui directories. To use the CLI version of the program, in terminal run:
    ```
    $ cd main
    $ ./main 
    ```
    To use the TUI version:
    ```
    $ cd tui
    $ ./tui_main
    ```
The binaries are created for Linux based operating systems. 

## Troubleshooting 

While trying to compile and run the program, you might run into errors or missing dependencies.

Common errors:

1. Sqlite-gorm error during build 

    ```
    #0 11.77 # gorm.io/driver/sqlite
    #0 11.77 /go/pkg/mod/gorm.io/driver/sqlite@v1.5.0/error_translator.go:9:35: undefined: sqlite3.ErrNoExtended
    #0 11.77 /go/pkg/mod/gorm.io/driver/sqlite@v1.5.0/error_translator.go:14:36: undefined: sqlite3.Error
    ```
    To fix this, in terminal run:
    
    ```
    $ cd 2023_MPass/main
    $ go get gorm.io/driver/sqlite@v1.4.4
    ```
2. `Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work` runtime error*
    ```
    $ sudo apt update
    $ sudo apt install build-essential
    $ cd 2023_MPass/main
    $ CGO_ENABLED=1 go build main.go
    ```
    *with cgo enabled compile time can take even up to a few minutes
    
## Authors
Project created by
Ana Mihajlović ([cholesski](https://github.com/cholesski)) and 
Katarina Grbović ([gkatarina](https://github.com/gkatarina))
