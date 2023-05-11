#!/bin/bash

cd main
go get gorm.io/driver/sqlite@v1.4.4
go build main.go
cd ../tui
go build tui_main.go
./tui_main
