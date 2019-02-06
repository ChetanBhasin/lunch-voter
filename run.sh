#!/bin/bash

git reset --hard origin/master
mkdir -p .results
git pull
git checkout -b user-local
./voter
git add .
git commit -m "Another vote"
git checkout master
git merge user-local
git push