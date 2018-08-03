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
	"go/aml"
	"go/ezmq"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const AML_FILE_PATH = "sample_data_model.aml"
const TNS_CONFIG_FILE_PATH = "tnsConf.json"
const LOCAL_HOST = "localhost"

func getAMLObject() *aml.AMLObject {
	// create "Model" data
	model, _ := aml.CreateAMLData()
	model.SetValueStr("ctname", "Model_107.113.97.248")
	model.SetValueStr("con", "SR-P7-970")

	// create "Sample" data
	axis, _ := aml.CreateAMLData()
	axis.SetValueStr("x", "20")
	axis.SetValueStr("y", "110")
	axis.SetValueStr("z", "80")

	info, _ := aml.CreateAMLData()
	info.SetValueStr("id", "f437da3b")
	info.SetValueAMLData("axis", axis)

	sample, _ := aml.CreateAMLData()
	sample.SetValueAMLData("info", info)
	appendix := [3]string{"935", "52303", "1442"}
	sample.SetValueStrArr("appendix", appendix[:])

	// set data to object
	amlObj, _ := aml.CreateAMLObject("Robot0001", time.Now().Format("20060102150405"))
	amlObj.AddData("Model", model)
	amlObj.AddData("Sample", sample)
	return amlObj
}

func publishData(publisher *ezmqx.EZMQXAMLPublisher, amlObject *aml.AMLObject, numberOfEvents int) {
	fmt.Printf("\n--------- Will Publish events at interval of 1 seconds ---------\n")
	for i := 0; i < numberOfEvents; i++ {
		result := publisher.Publish(amlObject)
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Error while publishing")
		}
		fmt.Println("Published event result:", result)
		time.Sleep(1000 * time.Millisecond)
	}
}

func printPubError() {
	fmt.Printf("\nRe-run the application as shown in below example: \n")
	fmt.Printf("\n  (1) For running in standalone mode: ")
	fmt.Printf("\n      ./publisher -t /topic -port 5562\n")
	fmt.Printf("\n  (2)  For running in docker mode: ")
	fmt.Printf("\n      ./publisher -t /topic \n")
	os.Exit(-1)
}

func main() {
	var port int
	var topic string
	var configInstance *ezmqx.EZMQXConfig = nil
	var publisher *ezmqx.EZMQXAMLPublisher = nil
	var isStandAlone bool

	// get port/topic from command line arguments
	if len(os.Args) != 3 && len(os.Args) != 5 {
		printPubError()
	}

	for n := 1; n < len(os.Args); n++ {
		if 0 == strings.Compare(os.Args[n], "-port") {
			port, _ = strconv.Atoi(os.Args[n+1])
			fmt.Println("Given Port: ", port)
			n = n + 1
			isStandAlone = true
		} else if 0 == strings.Compare(os.Args[n], "-t") {
			topic = os.Args[n+1]
			fmt.Println("Topic is : ", topic)
			n = n + 1
		} else {
			printPubError()
		}
	}

	//Handler for ctrl+c
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-osSignal
		fmt.Println(sig)
		if nil != publisher {
			publisher.Terminate()
		}
		if nil != configInstance {
			configInstance.Reset()
		}
		os.Exit(0)
	}()

	//get singleton instance
	configInstance = ezmqx.GetConfigInstance()

	//Initialize the EZMQX SDK
	if true == isStandAlone {
		result := configInstance.StartStandAloneMode(LOCAL_HOST, false, "")
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
		fmt.Println("Add aml model: failed")
		os.Exit(-1)
	}

	publisher, errorCode = ezmqx.GetAMLPublisher(topic, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), port)
	if errorCode != ezmq.EZMQ_OK {
		fmt.Println("Get publiser failed")
		os.Exit(-1)
	}

	//create AML object
	amlObject := getAMLObject()

	// This delay is added to prevent ZeroMQ first packet drop during
	// initial connection of publisher and subscriber.
	time.Sleep(1000 * time.Millisecond)

	if isStandAlone {
		publishData(publisher, amlObject, 15)
	} else {
		publishData(publisher, amlObject, 100000)
	}
	result := publisher.Terminate()
	if result != ezmqx.EZMQX_OK {
		fmt.Printf("Error while terminating publisher")
		os.Exit(-1)
	}
	result = configInstance.Reset()
	if result != ezmqx.EZMQX_OK {
		fmt.Printf("Error while resetting config")
		os.Exit(-1)
	}
	fmt.Println("Reset done: ", result)
}
