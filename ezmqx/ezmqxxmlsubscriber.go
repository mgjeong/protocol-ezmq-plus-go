/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
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

package ezmqx

import (
	"container/list"
	"go.uber.org/zap"
	"go/aml"
	"go/ezmq"
)

// Callback to get all the subscribed events for a specific topic.
type EZMQXXmlSubCB func(topic string, data string)

// Callback to get error for the subscribed topic.
type EZMQXXmlErrorCB func(topic string, errorCode EZMQXErrorCode)

// Structure represents EZMQX XML subscriber.
type EZMQXXMLSubscriber struct {
	subscriber    *EZMQXSubscriber
	subCallback   EZMQXXmlSubCB
	errorCallback EZMQXXmlErrorCB
}

// Get XML subscriber instance for given topic.
// It will work, if EZMQX is configured in docker mode.
func GetXMLDockerSubscriber(topic string, isHierarchical bool, subCallback EZMQXXmlSubCB, errorCallback EZMQXXmlErrorCB) (*EZMQXXMLSubscriber, EZMQXErrorCode) {
	instance := createXmlSubscriber(subCallback, errorCallback)
	result := instance.subscriber.initialize(topic, isHierarchical)
	if result != EZMQX_OK {
		Logger.Error("initialization failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Get XML subscriber instance for given topic.
// It will work, if EZMQX is configured in standalone mode.
func GetXMLStandAloneSubscriber(topic EZMQXTopic, subCallback EZMQXXmlSubCB, errorCallback EZMQXXmlErrorCB) (*EZMQXXMLSubscriber, EZMQXErrorCode) {
	instance := createXmlSubscriber(subCallback, errorCallback)
	ezmqxTopicList := list.New()
	ezmqxTopicList.PushBack(topic)
	result := instance.subscriber.storeTopics(*ezmqxTopicList)
	if result != EZMQX_OK {
		Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Get XML subscriber instance for given topic list.
// It will work, if EZMQX is configured in standalone mode.
func GetXMLStandAloneSubscriber1(topics list.List, subCallback EZMQXXmlSubCB, errorCallback EZMQXXmlErrorCB) (*EZMQXXMLSubscriber, EZMQXErrorCode) {
	instance := createXmlSubscriber(subCallback, errorCallback)
	result := instance.subscriber.storeTopics(topics)
	if result != EZMQX_OK {
		Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Terminate EZMQX XML subscriber.
func (instance *EZMQXXMLSubscriber) Terminate() EZMQXErrorCode {
	subscriber := instance.subscriber
	if nil == subscriber {
		Logger.Error("Subscriber is null")
		return EZMQX_UNKNOWN_STATE
	}
	return subscriber.terminate()
}

// Check whether subscriber is terminated or not.
func (instance *EZMQXXMLSubscriber) IsTerminated() (bool, EZMQXErrorCode) {
	subscriber := instance.subscriber
	if nil == subscriber {
		return false, EZMQX_UNKNOWN_STATE
	}
	return subscriber.isTerminated(), EZMQX_OK
}

// Get list of topics that subscribed by this subscriber.
func (instance *EZMQXXMLSubscriber) GetTopics() (*list.List, EZMQXErrorCode) {
	subscriber := instance.subscriber
	if nil == subscriber {
		return nil, EZMQX_UNKNOWN_STATE
	}
	return subscriber.getTopics(), EZMQX_OK
}

func createXmlSubscriber(subCallback EZMQXXmlSubCB, errorCallback EZMQXXmlErrorCB) *EZMQXXMLSubscriber {
	var instance *EZMQXXMLSubscriber
	instance = &EZMQXXMLSubscriber{}
	instance.subCallback = subCallback
	instance.errorCallback = errorCallback
	instance.subscriber = getEZMQXSubscriber()
	subscriber := instance.subscriber
	subscriber.internalCB = func(topic string, ezmqMsg ezmq.EZMQMessage) {
		representation := subscriber.amlRepDic[topic]
		if 0 == len(topic) || nil == representation {
			instance.errorCallback(topic, EZMQX_UNKNOWN_TOPIC)
			return
		}
		ezmqByteData := ezmqMsg.(ezmq.EZMQByteData)
		amlObject, result := representation.ByteToData(ezmqByteData.ByteData)
		if result != aml.AML_OK {
			instance.errorCallback(topic, EZMQX_BROKEN_PAYLOAD)
			return
		}
		amlString, result := representation.DataToAml(amlObject)
		instance.subCallback(topic, amlString)
	}
	return instance
}
