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
	"go.uber.org/zap"
)

// Structure represents EZMQX topic discovery.
type EZMQXTopicDiscovery struct {
	ezmqxCtx *EZMQXContext
}

// Get EZMQX topic discovery instance.
func GetEZMQXTopicDiscovery() (*EZMQXTopicDiscovery, EZMQXErrorCode) {
	context := getContextInstance()
	if !context.isCtxInitialized() {
		return nil, EZMQX_NOT_INITIALIZED
	}
	var instance *EZMQXTopicDiscovery
	instance = &EZMQXTopicDiscovery{}
	instance.ezmqxCtx = context
	return instance, EZMQX_OK
}

// Query the given topic to TNS [Topic name server] server.
func (instance *EZMQXTopicDiscovery) Query(topic string) (*EZMQXTopic, EZMQXErrorCode) {
	topics, result := instance.queryInternal(topic, false)
	if result != EZMQX_OK {
		return nil, result
	}
	return topics.Front().Value.(*EZMQXTopic), result
}

// Query the given topic to TNS [Topic name server] server.
// It will send query request with hierarchical option.
//
// For example: If topic name is /Topic then in success case TNS will
// return /Topic/A, /Topic/A/B etc.
func (instance *EZMQXTopicDiscovery) HierarchicalQuery(topic string) (*list.List, EZMQXErrorCode) {
	return instance.queryInternal(topic, true)
}

func (instance *EZMQXTopicDiscovery) queryInternal(topic string, isHierarchical bool) (*list.List, EZMQXErrorCode) {
	if instance.ezmqxCtx.isCtxTerminated() {
		return nil, EZMQX_TERMINATED
	}
	if !instance.ezmqxCtx.isCtxTnsEnabled() {
		return nil, EZMQX_TNS_NOT_AVAILABLE
	}
	result := validateTopic(topic)
	if false == result {
		return nil, EZMQX_INVALID_TOPIC
	}
	return instance.verifyTopic(topic, isHierarchical)
}

func (instance *EZMQXTopicDiscovery) parseTNSResponse(data []byte) (*list.List, EZMQXErrorCode) {
	ezmqxTopicList := list.New()
	topics := make(map[string][]interface{})
	err := json.Unmarshal([]byte(data), &topics)
	if err != nil {
		Logger.Error("parseTNSResponse: Unmarshal failed")
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
		isSecured, exists := stringMap[PAYLOAD_SECURED].(bool)
		if !exists {
			Logger.Error("No secured key exists in json response")
			return nil, EZMQX_REST_ERROR
		}
		ezmqXEndPoint := GetEZMQXEndPoint(endPoint)
		ezmqxTopic := GetEZMQXTopic(name, dataModel, isSecured, ezmqXEndPoint)
		ezmqxTopicList.PushBack(ezmqxTopic)
	}
	return ezmqxTopicList, EZMQX_OK
}

func (instance *EZMQXTopicDiscovery) verifyTopic(topic string, isHierarchical bool) (*list.List, EZMQXErrorCode) {
	tnsURL := instance.ezmqxCtx.ctxGetTnsAddr() + PREFIX + TOPIC
	Logger.Debug("[Topic discovery]", zap.String("Rest URL:", tnsURL))

	var hierarchical string
	if true == isHierarchical {
		hierarchical = QUERY_TRUE
	} else {
		hierarchical = QUERY_FALSE
	}
	query := QUERY_NAME + topic + QUERY_HIERARCHICAL + hierarchical
	Logger.Debug("[Topic discovery]", zap.String("query:", query))

	client := GetRestFactory()
	response, err := client.Get(tnsURL + QUESTION_MARK + query)
	if err != EZMQX_OK {
		Logger.Error("[Topic discovery]: request failed")
		return nil, EZMQX_REST_ERROR
	}
	if response.GetStatusCode() != HTTP_OK {
		Logger.Error("[Topic discovery]: Response code is not HTTP_OK")
		return nil, EZMQX_REST_ERROR
	}
	data := response.GetResponse()
	Logger.Debug("[Topic discovery]: ", zap.String("response:", string(data)))
	return instance.parseTNSResponse(data)
}
