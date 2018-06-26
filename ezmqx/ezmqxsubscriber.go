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
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go/aml"
	"go/ezmq"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type EZMQXSubCB func(topic string, ezmqMsg ezmq.EZMQMessage)

type EZMQXSubscriber struct {
	ezmqSubscriber *ezmq.EZMQSubscriber
	context        *EZMQXContext
	storedTopics   *list.List
	amlRepDic      map[string]*aml.Representation
	terminated     atomic.Value
	internalCB     EZMQXSubCB
	mutex          *sync.Mutex
}

func getEZMQXSubscriber() *EZMQXSubscriber {
	var instance *EZMQXSubscriber
	instance = &EZMQXSubscriber{}
	instance.context = getContextInstance()
	instance.storedTopics = list.New()
	instance.amlRepDic = make(map[string]*aml.Representation)
	instance.terminated.Store(false)
	instance.ezmqSubscriber = nil
	instance.mutex = &sync.Mutex{}
	return instance
}

func (instance *EZMQXSubscriber) initialize(topic string, isHierarchical bool) EZMQXErrorCode {
	context := instance.context
	if false == context.isCtxInitialized() {
		Logger.Error("Context is not initialized")
		return EZMQX_NOT_INITIALIZED
	}
	result := validateTopic(topic)
	if false == result {
		Logger.Error("Topic validation failed")
		return EZMQX_INVALID_TOPIC
	}
	if !context.isCtxTnsEnabled() {
		Logger.Error("TNS is not enabled")
		return EZMQX_TNS_NOT_AVAILABLE
	}
	verified, errorCode := instance.verifyTopics(topic, isHierarchical)
	if errorCode != EZMQX_OK {
		Logger.Error("Verify topics failed")
		return errorCode
	}
	return instance.storeTopics(*verified)
}

func (instance *EZMQXSubscriber) parseTNSResponse(data []byte) (*list.List, EZMQXErrorCode) {
	ezmqxTopicList := list.New()
	topics := make(map[string][]interface{})
	err := json.Unmarshal([]byte(data), &topics)
	if err != nil {
		return nil, EZMQX_REST_ERROR
	}
	topicList, exists := topics[PAYLOAD_TOPICS]
	if !exists {
		Logger.Error("No topics key exists in json response")
		return nil, EZMQX_REST_ERROR
	}
	for _, item := range topicList {
		stringMap := item.(map[string]interface{})
		dataModel, exists := stringMap[PAYLOAD_DATAMODEL].(string)
		if !exists {
			Logger.Error("No data model key exists in json response")
			return nil, EZMQX_REST_ERROR
		}
		endPoint, exists := stringMap[PAYLOAD_ENDPOINT].(string)
		if !exists {
			Logger.Error("No end point key exists in json response")
			return nil, EZMQX_REST_ERROR
		}
		name, exists := stringMap[PAYLOAD_NAME].(string)
		if !exists {
			Logger.Error("No name exists in json response")
			return nil, EZMQX_REST_ERROR
		}
		ezmqXEndPoint := GetEZMQXEndPoint(endPoint)
		ezmqxTopic := GetEZMQXTopic(name, dataModel, ezmqXEndPoint)
		topicValue := *ezmqxTopic
		ezmqxTopicList.PushBack(topicValue)
	}
	return ezmqxTopicList, EZMQX_OK
}

func (instance *EZMQXSubscriber) verifyTopics(topic string, isHierarchical bool) (*list.List, EZMQXErrorCode) {
	tnsURL := HTTP_PREFIX + instance.context.ctxGetTnsAddr() + COLON + TNS_KNOWN_PORT + PREFIX + TOPIC
	Logger.Debug("[TNS get topic]", zap.String("Rest URL:", tnsURL))
	var hierarchical string
	if true == isHierarchical {
		hierarchical = QUERY_TRUE
	} else {
		hierarchical = QUERY_FALSE
	}
	query := QUERY_NAME + topic + QUERY_HIERARCHICAL + hierarchical
	Logger.Debug("[TNS get topic]", zap.String("query:", query))

	timeout := time.Duration(CONNECTION_TIMEOUT * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Get(tnsURL + QUESTION_MARK + query)
	if err != nil {
		Logger.Debug("[TNS get topic] request failed")
		return nil, EZMQX_REST_ERROR
	}
	if response.StatusCode != HTTP_OK {
		Logger.Debug("[TNS get topic] Response code is not HTTP_OK")
		return nil, EZMQX_REST_ERROR
	}
	data, error := ioutil.ReadAll(response.Body)
	if error != nil {
		Logger.Error("[TNS get topic]: read body failed")
		return nil, EZMQX_REST_ERROR
	}
	Logger.Debug("[TNS get topic]", zap.String("response:", string(data)))
	return instance.parseTNSResponse(data)
}

func (instance *EZMQXSubscriber) createSubscriber(endPoint *EZMQXEndpoint) EZMQXErrorCode {
	instance.ezmqSubscriber = ezmq.GetEZMQSubscriber(endPoint.GetAddr(), endPoint.GetPort(), func(ezmqMsg ezmq.EZMQMessage) {},
		func(topic string, ezmqMsg ezmq.EZMQMessage) {
			contentType := ezmqMsg.GetContentType()
			fmt.Printf("\nTopic: %s", topic)
			if contentType == ezmq.EZMQ_CONTENT_TYPE_BYTEDATA {
				byteData := ezmqMsg.(ezmq.EZMQByteData)
				instance.internalCB(topic, byteData)
			} else {
				Logger.Debug("[Content type is not byte data")
			}
		})
	if nil == instance.ezmqSubscriber {
		Logger.Error("Ezmq subscriber is null")
		return EZMQX_UNKNOWN_STATE
	}
	result := instance.ezmqSubscriber.Start()
	if result != ezmq.EZMQ_OK {
		Logger.Error("Start ezmq subscriber failed", zap.Int("Error code:", int(result)))
		return EZMQX_UNKNOWN_STATE
	}
	Logger.Debug("Started ezmq subscriber", zap.Int("Error code:", int(result)))
	return EZMQX_OK
}

func (instance *EZMQXSubscriber) subscribe(topic EZMQXTopic) EZMQXErrorCode {
	endPoint := topic.GetEndPoint()
	if nil == instance.ezmqSubscriber {
		result := instance.createSubscriber(endPoint)
		if result != EZMQX_OK {
			Logger.Error("Create subscriber failed", zap.Int("Error code:", int(result)))
			return result
		}
	}
	ezmqSubscriber := instance.ezmqSubscriber
	errorCode := ezmqSubscriber.SubscribeWithIPPort(endPoint.GetAddr(), endPoint.GetPort(), topic.GetName())
	if errorCode != ezmq.EZMQ_OK {
		Logger.Error("Subscribe with IP port failed")
		return EZMQX_SESSION_UNAVAILABLE
	}
	Logger.Debug("Subscribed for topic", zap.String("Topic: ", topic.GetName()))
	return EZMQX_OK
}

func (instance *EZMQXSubscriber) storeTopics(topics list.List) EZMQXErrorCode {
	context := instance.context
	if false == context.isCtxInitialized() {
		return EZMQX_NOT_INITIALIZED
	}
	var result EZMQXErrorCode
	for topic := topics.Front(); topic != nil; topic = topic.Next() {
		ezmqxTopic := topic.Value.(EZMQXTopic)
		//validate topic
		isValid := validateTopic(ezmqxTopic.GetName())
		if !isValid {
			Logger.Error("Invalid topic")
			return EZMQX_INVALID_TOPIC
		}
		instance.amlRepDic[ezmqxTopic.GetName()], result = context.getAmlRep(ezmqxTopic.GetDataModel())
		if result != EZMQX_OK {
			Logger.Error("getAmlRep failed", zap.Int("Error code:", int(result)))
			return result
		}
		result = instance.subscribe(ezmqxTopic)
		if result != EZMQX_OK {
			Logger.Error("subscribe failed", zap.Int("Error code:", int(result)))
			return result
		}
		instance.storedTopics.PushBack(ezmqxTopic)
	}
	return EZMQX_OK
}

func (instance *EZMQXSubscriber) terminate() EZMQXErrorCode {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if true == instance.terminated.Load() {
		return EZMQX_TERMINATED
	}
	ezmqSubscriber := instance.ezmqSubscriber
	if ezmqSubscriber != nil {
		result := ezmqSubscriber.Stop()
		if result != ezmq.EZMQ_OK {
			Logger.Error("EZMQ subscriber stop: failed")
			return EZMQX_UNKNOWN_STATE
		}
	}
	instance.terminated.Store(true)
	return EZMQX_OK
}

func (instance *EZMQXSubscriber) isTerminated() bool {
	return instance.terminated.Load().(bool)
}

func (instance *EZMQXSubscriber) getTopics() *list.List {
	return instance.storedTopics
}
