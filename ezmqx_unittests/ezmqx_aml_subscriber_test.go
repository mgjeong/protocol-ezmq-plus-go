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
	"time"
)

var amlEventCount = 0

func amlSubCB(topic string, amlObject aml.AMLObject) {
	fmt.Printf("\naml SubCB")
	amlEventCount++
}
func errorCB(topic string, errorCode ezmqx.EZMQXErrorCode) {
	fmt.Printf("\naml errorCB")
}

func TestGetAMLStandAloneSubscriber(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestGetAMLStandAloneSubscriber1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic("", idList.Front().Value.(string), false, endPoint)
	_, result := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("subscriber is nil")
	}
	configInstance.Reset()
}

func TestGetSecuredAMLSubscriber(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	subscriber, _ := ezmqx.GetSecuredAMLSubscriber(*topic, utils.SERVER_PUBLIC_KEY, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	isSecured, _ := subscriber.IsSecured()
	if !isSecured {
		t.Errorf("subscriber is secured failed")
	}
	subscriber.Terminate()
	configInstance.Reset()
	// invalid topic
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	subscriber, _ = ezmqx.GetSecuredAMLSubscriber(*topic, utils.SERVER_PUBLIC_KEY, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if subscriber != nil {
		t.Errorf("Get subscriber wrong param failed")
	}
	// invalid key
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	subscriber, _ = ezmqx.GetSecuredAMLSubscriber(*topic, "", utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if subscriber != nil {
		t.Errorf("Get subscriber wrong param failed")
	}
}

func TestGetSecuredAMLSubscriberNegative(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic("", idList.Front().Value.(string), true, endPoint)
	//Invalid topic
	subscriber, _ := ezmqx.GetSecuredAMLSubscriber(*topic, utils.SERVER_PUBLIC_KEY, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if nil != subscriber {
		t.Errorf("subscriber is nil")
	}
	//Invalid server key
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	subscriber, _ = ezmqx.GetSecuredAMLSubscriber(*topic, "", utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if nil != subscriber {
		t.Errorf("subscriber is nil")
	}
}

func TestGetSecuredAMLSubscriber1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	topicKeyMap := make(map[ezmqx.EZMQXTopic]string)
	topicKeyMap[*topic] = utils.SERVER_PUBLIC_KEY
	subscriber, _ := ezmqx.GetSecuredAMLSubscriber1(topicKeyMap, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	isSecured, _ := subscriber.IsSecured()
	if !isSecured {
		t.Errorf("subscriber is secured failed")
	}
	subscriber.Terminate()
	configInstance.Reset()
	// invalid topic
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	topicKeyMap = make(map[ezmqx.EZMQXTopic]string)
	topicKeyMap[*topic] = utils.SERVER_PUBLIC_KEY
	subscriber, _ = ezmqx.GetSecuredAMLSubscriber1(topicKeyMap, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if subscriber != nil {
		t.Errorf("subscriber is nil")
	}
	//invalid key
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	topicKeyMap = make(map[ezmqx.EZMQXTopic]string)
	topicKeyMap[*topic] = utils.SERVER_PUBLIC_KEY
	subscriber, _ = ezmqx.GetSecuredAMLSubscriber1(topicKeyMap, "", utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if subscriber != nil {
		t.Errorf("subscriber is nil")
	}
}

func TestGetSecuredAMLSubscriber2(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), true, endPoint)
	topicKeyMap := make(map[ezmqx.EZMQXTopic]string)
	topicKeyMap[*topic] = utils.SERVER_PUBLIC_KEY
	topic2 := ezmqx.GetEZMQXTopic("/topic2", idList.Front().Value.(string), true, endPoint)
	topicKeyMap[*topic2] = utils.SERVER_PUBLIC_KEY2
	subscriber, _ := ezmqx.GetSecuredAMLSubscriber1(topicKeyMap, utils.CLIENT_PUBLIC_KEY, utils.CLIENT_SECRET_KEY, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestAMLSubscriberStandAlone(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.TEST_LOCAL_HOST, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}

	// Routine to publish data on socket
	go utils.Publish()

	// Wait till publisher is stopped
	<-utils.Exit_Chan

	time.Sleep(1000 * time.Millisecond)
	if amlEventCount < 5 {
		t.Errorf("Received less event")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestAMLSubscriberStandAlone1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	topicList := list.New()
	topicList.PushBack(*topic)
	subscriber, result := ezmqx.GetAMLStandAloneSubscriber1(*topicList, amlSubCB, errorCB)
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Get subscriber failed")
	}
	subscriber.Terminate()
	//Invalid topic
	topic = ezmqx.GetEZMQXTopic("", idList.Front().Value.(string), false, endPoint)
	topicList = list.New()
	topicList.PushBack(*topic)
	_, result = ezmqx.GetAMLStandAloneSubscriber1(*topicList, amlSubCB, errorCB)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Get subscriber failed")
	}
	configInstance.Reset()
	//Without config
	topic = ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	topicList = list.New()
	topicList.PushBack(*topic)
	_, result = ezmqx.GetAMLStandAloneSubscriber1(*topicList, amlSubCB, errorCB)
	if result != ezmqx.EZMQX_NOT_INITIALIZED {
		t.Errorf("Get subscriber failed")
	}
}

func TestAMLSubDockerMode(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.CONFIG_URL, []byte(utils.VALID_CONFIG_RESPONSE))
	utils.SetRestResponse(utils.TNS_INFO_URL, []byte(utils.VALID_TNS_INFO_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APPS_URL, []byte(utils.VALID_RUNNING_APPS_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APP_INFO_URL, []byte(utils.RUNNING_APP_INFO_RESPONSE))
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	instance.StartDockerMode(utils.TNS_CONFIG_FILE_PATH)
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	configInstance.AddAmlModel(*amlFilePath)
	utils.SetRestResponse(utils.SUB_TOPIC_H_URL, []byte(utils.SUB_TOPIC_RESPONSE))
	subscriber, _ := ezmqx.GetAMLSubscriber(utils.TOPIC, true, amlSubCB, errorCB)
	if nil == subscriber {
		t.Errorf("subscriber is nil")
	}
	subscriber.Terminate()
	configInstance.Reset()
}

func TestAMLSubDockerMode1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.CONFIG_URL, []byte(utils.VALID_CONFIG_RESPONSE))
	utils.SetRestResponse(utils.TNS_INFO_URL, []byte(utils.VALID_TNS_INFO_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APPS_URL, []byte(utils.VALID_RUNNING_APPS_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APP_INFO_URL, []byte(utils.RUNNING_APP_INFO_RESPONSE))
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	instance.StartDockerMode(utils.TNS_CONFIG_FILE_PATH)
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	configInstance.AddAmlModel(*amlFilePath)
	utils.SetRestResponse(utils.SUB_TOPIC_H_URL, []byte(utils.SUB_TOPIC_RESPONSE))
	_, result := ezmqx.GetAMLSubscriber("", true, amlSubCB, errorCB)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("subscriber is nil")
	}
	configInstance.Reset()
}

func TestSubTerminate(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
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
	result = subscriber.Terminate()
	if result != ezmqx.EZMQX_UNKNOWN_STATE {
		t.Errorf("Termination failed")
	}
	configInstance.Reset()
}

func TestGetTopics(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	endPoint := ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	topic := ezmqx.GetEZMQXTopic(utils.TOPIC, idList.Front().Value.(string), false, endPoint)
	subscriber, _ := ezmqx.GetAMLStandAloneSubscriber(*topic, amlSubCB, errorCB)
	_, error := subscriber.GetTopics()
	if error != ezmqx.EZMQX_OK {
		t.Errorf("GetTopics failed")
	}
	subscriber.Terminate()
	configInstance.Reset()
}
