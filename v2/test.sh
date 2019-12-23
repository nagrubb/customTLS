#!/usr/bin/env bash
go build -o client client.go
go build -o server server.go
./server &
server_pid=$!
sleep 1
./client
kill -9 $server_pid
