#!/usr/bin/env bash

# switch folder
# get the script path http://stackoverflow.com/questions/4774054/reliable-way-for-a-bash-script-to-get-the-full-path-to-itself
pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
# get current working directory
ORIGINAL_WD=${PWD}
# switch to script directory
cd ${SCRIPTPATH}
# switch to parent folder
cd ..

# add glide to path
export PATH=$PATH:${SCRIPTPATH}/linux-amd64
# show it is working
glide -v

# run the test
go test -v -cover $(glide novendor)

# go back to the old working directory
cd ${ORIGINAL_WD}
