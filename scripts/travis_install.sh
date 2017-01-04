#!/usr/bin/env bash

# install glide and enable vendor feature, NOTE: this only works on travis

# switch folder
# get the script path http://stackoverflow.com/questions/4774054/reliable-way-for-a-bash-script-to-get-the-full-path-to-itself
pushd `dirname $0` > /dev/null
SCRIPTPATH=`pwd -P`
popd > /dev/null
# get current working directory
ORIGINAL_WD=${PWD}
# switch to script directory
cd ${SCRIPTPATH}

# download and extract
wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
tar -zxvf glide-v0.12.3-linux-amd64.tar.gz
# add glide to path
export PATH=$PATH:${SCRIPTPATH}/linux-amd64
# show it is working
glide -v

# install dependencies
cd ..
glide install

# go back to the old working directory
cd ${ORIGINAL_WD}
