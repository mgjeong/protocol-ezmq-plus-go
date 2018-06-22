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

package ezmqx_unittests

import (
	"container/list"
	"fmt"
	"go/aml"
	"go/ezmqx"
	"go/ezmqx_unittests/utils"
	"testing"
)

func amlSubCB(topic string, amlObject aml.AMLObject)       { fmt.Printf("amlSubCB") }
func errorCB(topic string, errorCode ezmqx.EZMQXErrorCode) { fmt.Printf("errorCB") }

func TestGetAMLStandAloneSubscriber(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestGetAMLStandAloneSubscriber1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), endPoint)
	topicList := list.New()
	topicList.PushBack(*topic)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber1(*topicList, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestSubTerminate(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	isTerminated, _ := subscriber.IsTerminated()
	if true == isTerminated {
		t.Errorf("subscriber terminated")
	}
	result := subscriber.Terminate()
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Termination failed")
	}
	isTerminated, _ = subscriber.IsTerminated()
	if false == isTerminated {
		t.Errorf("Termination failed")
	}
	configInstance.Reset()
}

func TestGetTopics(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	_, error := subscriber.GetTopics()
	if error != ezmqx.EZMQX_OK {
		t.Errorf("GetTopics failed")
	}
	subscriber.Terminate()
	configInstance.Reset()
}
