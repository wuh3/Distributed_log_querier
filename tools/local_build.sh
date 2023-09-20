#!/bin/bash

cd ~/cs425/mp1/server/

go build -o server server.go

cd ~/cs425/mp1/client/
go build -o client client.go

echo server and client build complete.