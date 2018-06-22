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
type EZMQXAmlSubCB func(topic string, amlObject aml.AMLObject)

// Callback to get error for the subscribed topic.
type EZMQXAmlErrorCB func(topic string, errorCode EZMQXErrorCode)

// Structure represents EZMQX AML subscriber.
type EZMQXAMLSubscriber struct {
	subscriber    *EZMQXSubscriber
	subCallback   EZMQXAmlSubCB
	errorCallback EZMQXAmlErrorCB
}

// Get AML subscriber instance for given topic.
// It will work, if EZMQX is configured in docker mode.
func GetAMLDockerSubscriber(topic string, isHierarchical bool, subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) (*EZMQXAMLSubscriber, EZMQXErrorCode) {
	instance := createAmlSubscriber(subCallback, errorCallback)
	result := instance.subscriber.initialize(topic, isHierarchical)
	if result != EZMQX_OK {
		Logger.Error("initialization failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Get AML subscriber instance for given topic.
// It will work, if EZMQX is configured in standalone mode.
func GetAMLStandAloneSubscriber(topic EZMQXTopic, subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) (*EZMQXAMLSubscriber, EZMQXErrorCode) {
	instance := createAmlSubscriber(subCallback, errorCallback)
	ezmqxTopicList := list.New()
	ezmqxTopicList.PushBack(topic)
	result := instance.subscriber.storeTopics(*ezmqxTopicList)
	if result != EZMQX_OK {
		Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Get AML subscriber instance for given topic list.
// It will work, if EZMQX is configured in standalone mode.
func GetAMLStandAloneSubscriber1(topics list.List, subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) (*EZMQXAMLSubscriber, EZMQXErrorCode) {
	instance := createAmlSubscriber(subCallback, errorCallback)
	result := instance.subscriber.storeTopics(topics)
	if result != EZMQX_OK {
		Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	return instance, result
}

// Terminate EZMQX AML subscriber.
func (instance *EZMQXAMLSubscriber) Terminate() EZMQXErrorCode {
	subscriber := instance.subscriber
	if nil == subscriber {
		Logger.Error("Subscriber is null")
		return EZMQX_UNKNOWN_STATE
	}
	return subscriber.terminate()
}

// Check whether subscriber is terminated or not.
func (instance *EZMQXAMLSubscriber) IsTerminated() (bool, EZMQXErrorCode) {
	subscriber := instance.subscriber
	if nil == subscriber {
		return false, EZMQX_UNKNOWN_STATE
	}
	return subscriber.isTerminated(), EZMQX_OK
}

// Get list of topics that subscribed by this subscriber.
func (instance *EZMQXAMLSubscriber) GetTopics() (*list.List, EZMQXErrorCode) {
	subscriber := instance.subscriber
	if nil == subscriber {
		return nil, EZMQX_UNKNOWN_STATE
	}
	return subscriber.getTopics(), EZMQX_OK
}

func createAmlSubscriber(subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) *EZMQXAMLSubscriber {
	var instance *EZMQXAMLSubscriber
	instance = &EZMQXAMLSubscriber{}
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
		instance.subCallback(topic, *amlObject)
	}
	return instance
}
