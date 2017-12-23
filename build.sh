#!/bin/bash -e -x
docker build -t meshtest:latest --build-arg USER=$USER --build-arg USER_ID=`id -u $USER` --build-arg USER_GROUP_ID=`id -g $USER` .
docker run -e GITHUB_TOKEN=$GITHUB_TOKEN -e TRAVIS_TAG=$TRAVIS_TAG -v $GOPATH:/project meshtest ; echo "Done" ; docker container prune -f
cd cmd/experimentsrv
docker build -t experimentsrv .
cd ../echosrv
docker build -t echosrv .
cd ../timesrv
docker build -t timesrv .
cd ../restpoc
docker build -t restpoc .
cd ../..
