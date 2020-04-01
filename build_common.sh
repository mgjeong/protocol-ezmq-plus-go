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
EZMQX_WITH_SECURITY=true

IS_SECURED="secure"

PROTOCOL_EZMQ_GO_VERSION=v1.0_rc1
DATAMODEL_AML_GO_VERSION=v1.0_rc1

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
        git clone git@github.com:edgexfoundry-holding/protocol-ezmq-go.git
    fi
    
    # Build ezmq-go library
    cd $DEP_ROOT/protocol-ezmq-go
    git checkout ${PROTOCOL_EZMQ_GO_VERSION}

    echo -e "${GREEN}Building protocol-ezmq-go library and its dependencies${NO_COLOUR}"
    ./build_auto.sh --target_arch=${TARGET_ARCH} --with_dependencies=${EZMQX_WITH_DEP} --build_mode=${EZMQX_BUILD_MODE} --with_security=${EZMQX_WITH_SECURITY}
    echo -e "${GREEN}Install ezmq-go done${NO_COLOUR}"
    
    cd $DEP_ROOT
    # Check/clone library
    if [ -d "datamodel-aml-go" ] ; then
        echo "datamodel-aml-go folder exist"
    else
        git clone git@github.com:edgexfoundry-holding/datamodel-aml-go.git
    fi
    
    # Build aml-go library
    cd $DEP_ROOT/datamodel-aml-go
    git checkout ${DATAMODEL_AML_GO_VERSION}

    echo -e "${GREEN}Building datamodel-aml-go library and its dependencies${NO_COLOUR}"
    ./build_common.sh --target_arch=${TARGET_ARCH} --build_mode=${EZMQX_BUILD_MODE}
    echo -e "${GREEN}Install aml-go done${NO_COLOUR}"
}

build_x86_and_64() {
    cd $PROJECT_ROOT/src/go/
    #build ezmqx SDK
    cd ./ezmqx  
    go build -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" 
    go install
    #build ezmqx_samples
    cd ../ezmqx_samples
    if [ ${EZMQX_WITH_SECURITY} = true ]; then
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" topicdiscovery.go
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" amlsubscriber_secured.go
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" xmlsubscriber_secured.go      
    else
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" topicdiscovery.go
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" publisher.go 
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" amlsubscriber.go
        go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" xmlsubscriber.go 
    fi
}

build_armhf_native() {
    cd $PROJECT_ROOT/src/go/
    #build ezmqx SDK
    cd ./ezmqx
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" 
    CGO_ENABLED=1 GOOS=linux GOARCH=arm go install
    #build ezmqx_samples
    cd ../ezmqx_samples
    if [ ${EZMQX_WITH_SECURITY} = true ]; then
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" topicdiscovery.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" publisher_secured.go 
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" amlsubscriber_secured.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" xmlsubscriber_secured.go      
    else
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" topicdiscovery.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" publisher.go 
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" amlsubscriber.go
        CGO_ENABLED=1 GOOS=linux GOARCH=arm go build -a -tags="${EZMQX_BUILD_MODE} ${IS_SECURED}" xmlsubscriber.go 
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
    echo "  --with_dependencies=[true|false](default: false)            :  Build ezmq-plus along with dependencies [ezmq and aml]"
    echo "  --with_security=[true|false](default: true)                 :  Build ezmq library with or without Security feature"
    echo "  -c                                                          :  Clean ezmq-plus Repository and its dependencies"
    echo "  -h / --help                                                 :  Display help and exit [Be careful it will also remove GOPATH:src, pkg and bin]"
    echo -e "${GREEN}Note: ${NO_COLOUR}"
    echo "  - While building newly for any architecture use -with_dependencies=true option."
}

build_ezmqx() {
    echo -e "${GREEN}Target Arch is: $EZMQX_TARGET_ARCH${NO_COLOUR}" 
    echo -e "${GREEN}Build mode is: $EZMQX_BUILD_MODE${NO_COLOUR}"
    echo -e "${GREEN}Is security enabled: $EZMQX_WITH_SECURITY${NO_COLOUR}"
    echo -e "${GREEN}Build with depedencies: ${EZMQX_WITH_DEP}${NO_COLOUR}"
    
    if [ ${EZMQX_WITH_SECURITY} = false ]; then
        IS_SECURED="unsecure"
    fi
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
                shift 1;
                ;;
            --target_arch=*)
                EZMQX_TARGET_ARCH="${1#*=}";
                shift 1
                ;;
            --build_mode=*)
                EZMQX_BUILD_MODE="${1#*=}";
                shift 1;
                ;;
            --with_security=*)
                EZMQX_WITH_SECURITY="${1#*=}";
                if [ ${EZMQX_WITH_SECURITY} != true ] && [ ${EZMQX_WITH_SECURITY} != false ]; then
                    echo -e "${RED}Unknown option for --with_security${NO_COLOUR}"
                    shift 1; exit 0
                fi              
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
echo -e "Building ezMQ-plus-go library("${EZMQX_TARGET_ARCH}").."
build_ezmqx
echo -e "Done building ezMQ-plus-go library("${EZMQX_TARGET_ARCH}")"

