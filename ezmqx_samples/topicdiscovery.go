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
	"strings"
	"time"
)

const TNS_CONFIG_FILE_PATH = "tnsConf.json"
const LOCAL_HOST = "localhost"

func printTopicError() {
	fmt.Printf("\nRe-run the application as shown in below example: \n")
	fmt.Printf("\n  (1) For running in standalone mode: ")
	fmt.Printf("\n      ./topicdiscovery -t /topic -tns 192.168.1.1\n")
	fmt.Printf("\n  (2)  For running in docker mode: ")
	fmt.Printf("\n      ./topicdiscovery -t /topic \n")
	os.Exit(-1)
}

func printTopicList(topicList list.List) {
	if 0 == topicList.Len() {
		fmt.Println("Topic list is empty.....")
	}
	for ezmqxTopic := topicList.Front(); ezmqxTopic != nil; ezmqxTopic = ezmqxTopic.Next() {
		topic := ezmqxTopic.Value.(*ezmqx.EZMQXTopic)
		fmt.Println("=================================================")
		fmt.Println("Topic: ", topic.GetName())
		fmt.Println("Endpoint: ", topic.GetEndPoint().ToString())
		fmt.Println("Data Model: ", topic.GetDataModel())
		fmt.Println("=================================================")
	}
}

func main() {
	var tnsAddress string
	var topic string
	var configInstance *ezmqx.EZMQXConfig = nil
	var isStandAlone bool

	// get port/topic from command line arguments
	if len(os.Args) != 3 && len(os.Args) != 5 {
		printTopicError()
	}
	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-tns") {
			tnsAddress = os.Args[n+1]
			fmt.Println("Given TNS address: ", tnsAddress)
			n = n + 1
			isStandAlone = true
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Println("Topic is : ", topic)
			n = n + 1
		} else {
			printTopicError()
		}
	}
	//get singleton instance
	configInstance = ezmqx.GetConfigInstance()

	//Initialize the EZMQX SDK
	if true == isStandAlone {
		result := configInstance.StartStandAloneMode(LOCAL_HOST, true, tnsAddress)
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start stand alone mode failed")
			os.Exit(-1)
		}
		fmt.Println("Stand alone mode started")
	} else {
		result := configInstance.StartDockerMode(TNS_CONFIG_FILE_PATH)
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start docker mode failed")
			os.Exit(-1)
		}
		fmt.Println("Docker mode started")
	}
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()
	topicList, errorCode := topicDiscovery.HierarchicalQuery(topic)
	fmt.Println("Topic discovery query respone: ", errorCode)
	if errorCode == ezmqx.EZMQX_OK {
		printTopicList(*topicList)
	} else {
		os.Exit(-1)
	}

	// Wait for 5 minutes before exit [For docker mode].
	if false == isStandAlone {
		fmt.Println("Waiting for 5 minutes before program exit for docker mode... [press ctrl+c to exit]")
		time.Sleep(5 * time.Minute)
	}
	result := configInstance.Reset()
	fmt.Println("Reset Config done: ", result)
}
