/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

package main

import (
	"go/ezmqx"

	"container/list"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const AML_FILE_PATH = "sample_data_model.aml"

func printError() {
	fmt.Printf("\nRe-run the application as shown in below examples: \n")
	fmt.Printf("\n  (1) For running in standalone mode: ")
	fmt.Printf("\n     ./xmlsubscriber -ip 192.168.1.1 -port 5562 -t /topic\n")
	fmt.Printf("\n  (2) For running in docker mode: ")
	fmt.Printf("\n     ./xmlsubscriber -t /topic -h true\n")
	fmt.Printf("\n Note: -h stands for hierarchical search for topic from TNS server\n")
	os.Exit(-1)
}

func main() {
	var ip string
	var port int
	var topic string
	var hierarchical bool
	var subscriber *ezmqx.EZMQXXMLSubscriber
	var result ezmqx.EZMQXErrorCode
	var isStandAlone bool = false
	var configInstance *ezmqx.EZMQXConfig = nil
	var isSubscribed bool = false

	// get ip and port from command line arguments
	if len(os.Args) != 5 && len(os.Args) != 7 {
		printError()
	}

	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-ip") {
			ip = os.Args[n+1]
			fmt.Printf("\nGiven Ip: %s", ip)
			n = n + 1
			isStandAlone = true
		} else if 0 == strings.Compare(os.Args[n], "-port") {
			port, _ = strconv.Atoi(os.Args[n+1])
			fmt.Printf("\nGiven Port %d: ", port)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Printf("Topic is : %s", topic)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-h") {
			isHierarchical := os.Args[n+1]
			hierarchical, _ := strconv.ParseBool(isHierarchical)
			fmt.Println("Is hierarchical: ", hierarchical)
			n = n + 1
		} else {
			printError()
		}
	}

	// Handler for ctrl+c
	osSignal := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignal
		fmt.Println(sig)
		if false == isSubscribed {
			os.Exit(-1)
		}
		if nil != subscriber {
			subscriber.Terminate()
		}
		if nil != configInstance {
			configInstance.Reset()
		}
		exit <- true
	}()

	// callbacks
	amlSubCB := func(topic string, data string) {
		fmt.Println("\n[APP Callback] Topic : " + topic + "\n")
		fmt.Println("[APP Callback] Data : " + data + "\n")
	}
	amlErrorCB := func(topic string, errorCode ezmqx.EZMQXErrorCode) {
		fmt.Println("\n[APP Error Callback] ErrorCode : ", errorCode)
	}

	//get singleton instance
	configInstance = ezmqx.GetConfigInstance()

	//Initialize the EZMQX SDK
	if true == isStandAlone {
		result := configInstance.StartStandAloneMode(false, "")
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start Stand alone mode failed")
			os.Exit(-1)
		}
		fmt.Println("Stand alone mode started")

	} else {
		result := configInstance.StartDockerMode()
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start Docker mode failed")
			os.Exit(-1)
		}
		fmt.Println("Docker mode started")
	}
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, errorCode := configInstance.AddAmlModel(*amlFilePath)
	if ezmqx.EZMQX_OK == errorCode {
		for id := idList.Front(); id != nil; id = id.Next() {
			fmt.Println("id: ", id.Value.(string))
		}
	} else {
		fmt.Println("Add AML model failed")
		os.Exit(-1)
	}

	if isStandAlone {
		endPoint := ezmqx.GetEZMQXEndPoint1(ip, port)
		ezmqxTopic := ezmqx.GetEZMQXTopic(topic, idList.Front().Value.(string), endPoint)
		subscriber, result = ezmqx.GetXMLStandAloneSubscriber(*ezmqxTopic, amlSubCB, amlErrorCB)

	} else {
		subscriber, result = ezmqx.GetXMLDockerSubscriber(topic, hierarchical, amlSubCB, amlErrorCB)
	}
	if result != ezmqx.EZMQX_OK {
		fmt.Println("Get XML subscriber failed")
		os.Exit(-1)
	}
	isSubscribed = true
	fmt.Printf("\nSuscribed to publisher.. -- Waiting for Events --\n")
	<-exit
	fmt.Println("exiting")
}
