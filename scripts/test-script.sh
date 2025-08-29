#!/bin/bash

cd ..
make build
rm -rf sandbox/
mkdir sandbox
cp ./bin/arm ./sandbox
cd sandbox

./arm config add registry ar https://github.com/PatrickJS/awesome-cursorrules/ --type git
./arm config add sink q --directories .amazonq/rules --include ar/*
