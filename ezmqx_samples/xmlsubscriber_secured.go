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
const TNS_CONFIG_FILE_PATH = "tnsConf.json"
const LOCAL_HOST = "localhost"

func printError() {
	fmt.Printf("\nRe-run the application as shown in below examples: \n")
	fmt.Printf("\n  (1) For running in standalone mode: ")
	fmt.Printf("\n     ./xmlsubscriber_secured -ip 192.168.1.1 -port 5562 -t /topic -secured 1\n")
	fmt.Printf("\n  (2) For running in standalone mode [[With TNS]: ")
	fmt.Printf("\n     ./xmlsubscriber_secured -t /topic -tns 192.168.10.1 -secured 1\n")
	fmt.Printf("\n  (3) For running in docker mode: ")
	fmt.Printf("\n     ./xmlsubscriber_secured -t /topic -secured 1\n")
	fmt.Printf("\n Note:")
	fmt.Printf("\n (1) While testing standalone mode without TNS, Make sure to give same topic on both publisher and subscriber")
	fmt.Printf("\n (2) While testing TNS/docker mode  + secured mode, Make sure to give same topic on both publisher and subscriber")
	fmt.Printf("\n (3) docker mode will work only when sample is running in docker container\n")
	os.Exit(-1)
}

func getTNSAddress(tnsAddress string) string {
	return "http://" + tnsAddress + ":80/tns-server"
}

func main() {
	var ip string
	var port int
	var topic string
	var subscriber *ezmqx.EZMQXXMLSubscriber
	var result ezmqx.EZMQXErrorCode
	var isStandAlone bool = false
	var configInstance *ezmqx.EZMQXConfig = nil
	var isSubscribed bool = false
	var tnsAddr string = ""
	var isSecured int = 0

	// get ip and port from command line arguments
	if len(os.Args) != 5 && len(os.Args) != 7 && len(os.Args) != 9 {
		printError()
	}

	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-ip") {
			ip = os.Args[n+1]
			fmt.Println("Given Ip: ", ip)
			n = n + 1
			isStandAlone = true
		} else if 0 == strings.Compare(os.Args[n], "-port") {
			port, _ = strconv.Atoi(os.Args[n+1])
			fmt.Println("Given Port: ", port)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Println("Topic is: ", topic)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-tns") {
			tnsAddr = os.Args[n+1]
			fmt.Println("TNS Address is : ", tnsAddr)
			n = n + 1
		} else if 0 == strings.Compare(os.Args[n], "-secured") {
			isSecured, _ = strconv.Atoi(os.Args[n+1])
			fmt.Println("Is secured: ", isSecured)
			n = n + 1
			isStandAlone = true
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
		if 0 == len(tnsAddr) {
			result = configInstance.StartStandAloneMode(LOCAL_HOST, false, "")
		} else {
			tnsAddress := getTNSAddress(tnsAddr)
			fmt.Println("Complete TNS address is: " + tnsAddress)
			result = configInstance.StartStandAloneMode("", true, tnsAddress)
		}
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start Stand alone mode: failed")
			os.Exit(-1)
		}
		fmt.Println("Stand alone mode started")
	} else {
		result := configInstance.StartDockerMode(TNS_CONFIG_FILE_PATH)
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Start docker mode: failed")
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

	// this key only used on sample
	// TODO change to real key.
	const serverPublicKey = "tXJx&1^QE2g7WCXbF.$$TVP.wCtxwNhR8?iLi&S<"
	const clientPublicKey = "-QW?Ved(f:<::3d5tJ$[4Er&]6#9yr=vha/caBc("
	const clientSecretKey = "ZB1@RS6Kv^zucova$kH(!o>tZCQ.<!Q)6-0aWFmW"

	if isStandAlone {
		if 0 == isSecured {
			if 0 == len(tnsAddr) {
				endPoint := ezmqx.GetEZMQXEndPoint1(ip, port)
				ezmqxTopic := ezmqx.GetEZMQXTopic(topic, idList.Front().Value.(string), false, endPoint)
				subscriber, result = ezmqx.GetXMLStandAloneSubscriber(*ezmqxTopic, amlSubCB, amlErrorCB)
			} else {
				subscriber, result = ezmqx.GetXMLSubscriber(topic, true, amlSubCB, amlErrorCB)
			}
		} else {
			if 0 == len(tnsAddr) {
				endPoint := ezmqx.GetEZMQXEndPoint1(ip, port)
				ezmqxTopic := ezmqx.GetEZMQXTopic(topic, idList.Front().Value.(string), true, endPoint)
				subscriber, result = ezmqx.GetSecuredXMLSubscriber(*ezmqxTopic, serverPublicKey, clientPublicKey, clientSecretKey, amlSubCB, amlErrorCB)
			} else {
				topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()
				ezmqxTopic, errorCode := topicDiscovery.Query(topic)
				fmt.Println("Topic discovery query respone: ", errorCode)
				if errorCode != ezmqx.EZMQX_OK {
					os.Exit(-1)
				}
				subscriber, result = ezmqx.GetSecuredXMLSubscriber(*ezmqxTopic, serverPublicKey, clientPublicKey, clientSecretKey, amlSubCB, amlErrorCB)
			}
		}
	} else {
		if 0 == isSecured {
			subscriber, result = ezmqx.GetXMLSubscriber(topic, true, amlSubCB, amlErrorCB)
		} else {
			topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()
			ezmqxTopic, errorCode := topicDiscovery.Query(topic)
			fmt.Println("Topic discovery query respone: ", errorCode)
			if errorCode != ezmqx.EZMQX_OK {
				os.Exit(-1)
			}
			subscriber, result = ezmqx.GetSecuredXMLSubscriber(*ezmqxTopic, serverPublicKey, clientPublicKey, clientSecretKey, amlSubCB, amlErrorCB)
		}
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
