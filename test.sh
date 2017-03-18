#!/bin/bash
go run src/main/wc_s3.go master localhost:7777 pg > master.log &
go run src/main/wc_s3.go worker localhost:7777 localhost:7778 > worker_1.log &
go run src/main/wc_s3.go worker localhost:7777 localhost:7779 > worker_2.log &
