###############################################################################
# Copyright 2018 Samsung Electronics All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
###############################################################################

#!/bin/bash
set +e
#Colors
RED="\033[0;31m"
GREEN="\033[0;32m"
BLUE="\033[0;34m"
NO_COLOUR="\033[0m"

PROJECT_ROOT=$(pwd)
export GOPATH=$PWD
DEP_ROOT=$(pwd)/dependencies
EZMQX_TARGET_ARCH="$(uname -m)"
EZMQX_INSTALL_PREREQUISITES=false
EZMQX_WITH_DEP=false
EZMQX_BUILD_MODE="release"

install_dependencies() { 
    TARGET_ARCH=${EZMQX_TARGET_ARCH}
    if [ "armhf" = ${TARGET_ARCH} ]; then
        TARGET_ARCH="armhf-native";
    fi   
    if [ -d "./dependencies" ] ; then
        echo "dependencies folder exist"
    else
        mkdir dependencies
    fi
    
    cd $DEP_ROOT
    # Check/clone ezmq-go library
    if [ -d "protocol-ezmq-go" ] ; then
        echo "protocol-ezmq-go folder exist"
    else
        git clone git@github.sec.samsung.net:RS7-EdgeComputing/protocol-ezmq-go.git
    fi

    # Build ezmq-go library
    cd $DEP_ROOT/protocol-ezmq-go
    echo -e "${GREEN}Building protocol-ezmq-go library and its dependencies${NO_COLOUR}"
    ./build_auto.sh --target_arch=${TARGET_ARCH} --with_dependencies=${EZMQX_WITH_DEP} --build_mode=${EZMQX_BUILD_MODE}
    echo -e "${GREEN}Install ezmq-go done${NO_COLOUR}"
    
    cd $DEP_ROOT
    # Check/clone library
    if [ -d "datamodel-aml-go" ] ; then
        echo "datamodel-aml-go folder exist"
    else
        git clone git@github.sec.samsung.net:RS7-EdgeComputing/datamodel-aml-go.git
    fi

    # Build aml-go library
    cd $DEP_ROOT/datamodel-aml-go
    echo -e "${GREEN}Building datamodel-aml-go library and its dependencies${NO_COLOUR}"
    ./build_common.sh --target_arch=${TARGET_ARCH} --build_mode=${EZMQX_BUILD_MODE}
    echo -e "${GREEN}Install aml-go done${NO_COLOUR}"
}

build_x86_and_64() {
    cd $PROJECT_ROOT/src/go/
    cd ./ezmqx
    if [ "debug" = ${EZMQX_BUILD_MODE} ]; then
        go build -tags=debug
        go install
        #build ezmqx_samples
        cd ../ezmqx_samples
        go build -a -tags=debug topicdiscovery.go
        go build -a -tags=debug publisher.go	
        go build -a -tags=debug amlsubscriber.go		
        go build -a -tags=debug xmlsubscriber.go
    else
        go build
        go install
        #build ezmqx_samples
        cd ../ezmqx_samples
        go build -a topicdiscovery.go
        go build -a publisher.go	
        go build -a amlsubscriber.go		
        go build -a xmlsubscriber.go
    fi
}

build_armhf_native() {
    cd $PROJECT_ROOT/src/go/
    cd ./ezmqx
    if [ "debug" = ${EZMQX_BUILD_MODE} ]; then
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -tags=debug
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go install
        #build ezmqx_samples
        cd ../ezmqx_samples
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags=debug topicdiscovery.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags=debug publisher.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags=debug amlsubscriber.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags=debug xmlsubscriber.go
    else
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go install
        #build ezmqx_samples
        cd ../ezmqx_samples
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a topicdiscovery.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a publisher.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a amlsubscriber.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a xmlsubscriber.go
    fi
}

clean_ezmqx() {
    echo -e "Cleaning ${BLUE}${PROJECT_ROOT}${NO_COLOUR}"
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/src/${NO_COLOUR}"
    rm -rf ./src
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/pkg/${NO_COLOUR}"
    rm -rf ./pkg
    echo -e "Deleting  ${RED}${PROJECT_ROOT}/dependencies/${NO_COLOUR}"
    rm -rf ./dependencies
    echo -e "Finished Cleaning ${BLUE}${AML}${NO_COLOUR}"
}

usage() {
    echo -e "${BLUE}Usage:${NO_COLOUR} ./build_common.sh <option>"
    echo -e "${GREEN}Options:${NO_COLOUR}"
    echo "  --target_arch=[x86|x86_64|armhf]                            :  Choose Target Architecture"
    echo "  --build_mode=[release|debug](default: release)              :  Build in release or debug mode"
    echo "  --with_dependencies=(default: false)                        :  Build ezmq-plus along with dependencies [ezmq and aml]"
    echo "  -c                                                          :  Clean ezmq-plus Repository and its dependencies"
    echo "  -h / --help                                                 :  Display help and exit [Be careful it will also remove GOPATH:src, pkg and bin]"
}

build_ezmqx() { 
    if [ ${EZMQX_WITH_DEP} = true ]; then
        install_dependencies
        # Copy "ezmq-go" package to GOPATH
		cd $PROJECT_ROOT
        cp -r $DEP_ROOT/protocol-ezmq-go/src/ .
        # Copy "aml-go" package to GOPATH
        cp -r $DEP_ROOT/datamodel-aml-go/src/go/aml/ ./src/go/
        # copy aml libs 
        cd $PROJECT_ROOT/src/go
        mkdir ezmqx_extlibs && cd ezmqx_extlibs
        cp $DEP_ROOT/datamodel-aml-go/dependencies/datamodel-aml-c/dependencies/datamodel-aml-cpp/out/linux/${EZMQX_TARGET_ARCH}/${EZMQX_BUILD_MODE}/libaml.so .
        cp $DEP_ROOT/datamodel-aml-go/dependencies/datamodel-aml-c/out/linux/${EZMQX_TARGET_ARCH}/${EZMQX_BUILD_MODE}/libcaml.so .
    fi
    
    cd $PROJECT_ROOT
    #copy ezmq-plus SDK files
    cp -r ezmqx ./src/go
    #copy ezmq-plus ezmqx_samples
    cp -r ezmqx_samples ./src/go
    # Copy ezmq-plus unit test cases
    cp -r ezmqx_unittests ./src/go
	
    # set flags for aml-go includes and libs 
    export CGO_CFLAGS=-I$PWD/dependencies/datamodel-aml-go/dependencies/datamodel-aml-c/include/
    export CGO_LDFLAGS=-L$PWD/src/go/ezmqx_extlibs
    export CGO_LDFLAGS=$CGO_LDFLAGS" -lcaml -laml"
    
    if [ "x86" = ${EZMQX_TARGET_ARCH} ]; then
         build_x86_and_64;
    elif [ "x86_64" = ${EZMQX_TARGET_ARCH} ]; then
         build_x86_and_64;
    elif [ "armhf" = ${EZMQX_TARGET_ARCH} ]; then
         build_armhf_native;
    else
         echo -e "${RED}Not a supported architecture${NO_COLOUR}"
         usage; exit 1;
    fi
}

process_cmd_args() {
    if [ "$#" -eq 0  ]; then
        echo -e "No argument.."
        usage; exit 1
    fi

    while [ "$#" -gt 0  ]; do
        case "$1" in
            --with_dependencies=*)
                EZMQX_WITH_DEP="${1#*=}";
                if [ ${EZMQX_WITH_DEP} != true ] && [ ${EZMQX_WITH_DEP} != false ]; then
                    echo -e "${RED}Unknown option for --with_dependencies${NO_COLOUR}"
                    exit 1
                fi
                echo -e "${GREEN}Install dependencies [ezmq and aml] before build: ${EZMQX_WITH_DEP}${NO_COLOUR}"
                shift 1;
                ;;
            --target_arch=*)
                EZMQX_TARGET_ARCH="${1#*=}";
                echo -e "${GREEN}Target Arch is: $EZMQX_TARGET_ARCH${NO_COLOUR}"
                shift 1
                ;;
            --build_mode=*)
                EZMQX_BUILD_MODE="${1#*=}";
                echo -e "${GREEN}Build mode is: $EZMQX_BUILD_MODE${NO_COLOUR}"
                shift 1;
                ;;
            -c)
                clean_ezmqx
                shift 1; exit 0
                ;;
            -h)
                usage; exit 0
                ;;
            --help)
                usage; exit 0
                ;;
            -*)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
             *)
                echo -e "${RED}"
                echo "unknown option: $1" >&2;
                echo -e "${NO_COLOUR}"
                usage; exit 1
                ;;
        esac
    done
}

process_cmd_args "$@"
build_ezmqx
echo -e "${GREEN}ezmq-plus-go Build done${NO_COLOUR}"

