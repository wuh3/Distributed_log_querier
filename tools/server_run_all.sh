#!/bin/bash

for i in {1..9}
do
    ssh haozhew3@fa23-cs425-740$i.cs.illinois.edu "cd ~/cs425/mp1/tools/; ./server_run_local.sh"
    echo VM$i server started...
done
ssh haozhew3@fa23-cs425-7410.cs.illinois.edu "cd ~/cs425/mp1/tools/; ./server_run_local.sh"
echo VM10 server started...