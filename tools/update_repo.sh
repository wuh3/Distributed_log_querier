#!/bin/bash

for i in {1..9}
do
    echo VM$i
    ssh haozhew3@fa23-cs425-740$i.cs.illinois.edu "cd ~/cs425/; git checkout develop; git stash; git pull origin develop; exit"
done
ssh haozhew3@fa23-cs425-7410.cs.illinois.edu "cd ~/cs425/; git checkout develop; git stash; git pull origin develop; exit"
echo 'Repo Update complete.'