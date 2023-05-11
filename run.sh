#!/bin/bash

cd main
go build main.go
cd ../tui
go build tui_main.go
./tui_main