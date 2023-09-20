#!/bin/bash

cd ~/cs425/mp1/tools
./kill_server.sh
./local_build.sh
cd ../server
nohup ./server > /dev/null 2>&1 &
cd ..
echo Server is on.