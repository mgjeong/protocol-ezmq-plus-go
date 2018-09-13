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
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go/ezmq"
	"sync/atomic"
)

type EZMQXPublisher struct {
	ezmqPublisher *ezmq.EZMQPublisher
	context       *EZMQXContext
	topic         *EZMQXTopic
	topicHandler  *EZMQXTopicHandler
	localPort     int
	status        uint32
}

func getPublisher() *EZMQXPublisher {
	var instance *EZMQXPublisher
	instance = &EZMQXPublisher{}
	instance.context = getContextInstance()
	instance.status = CREATED
	return instance
}

func (instance *EZMQXPublisher) initialize(optionalPort int) EZMQXErrorCode {
	if !instance.context.isCtxInitialized() {
		return EZMQX_NOT_INITIALIZED
	}
	if instance.context.isCtxStandAlone() {
		instance.localPort = optionalPort
	} else {
		var error EZMQXErrorCode
		instance.localPort, error = instance.context.assignDynamicPort()
		if error != EZMQX_OK {
			return error
		}
	}
	// create ezmq publisher
	instance.ezmqPublisher = ezmq.GetEZMQPublisher(instance.localPort, nil, nil, nil)
	if nil == instance.ezmqPublisher {
		Logger.Error("Could not create ezmq publisher")
		return EZMQX_UNKNOWN_STATE
	}
	// Start ezmq publisher
	if ezmq.EZMQ_OK != instance.ezmqPublisher.Start() {
		Logger.Error("Could not start ezmq publisher")
		return EZMQX_UNKNOWN_STATE
	}
	// Init topic handler
	if instance.context.isCtxTnsEnabled() {
		instance.topicHandler = getTopicHandler()
		instance.topicHandler.initHandler()
		Logger.Debug("Initialized topic handler")
	}
	atomic.StoreUint32(&instance.status, INITIALIZED)
	return EZMQX_OK
}

func (instance *EZMQXPublisher) parseTopicResponse(response RestResponse) EZMQXErrorCode {
	statusCode := response.GetStatusCode()
	Logger.Debug("parseTopicResponse ", zap.Int(" Status code: ", statusCode))
	if statusCode != HTTP_CREATED {
		Logger.Error("parseTopicResponse, status code is not HTTP_CREATED")
		return EZMQX_REST_ERROR
	}
	data := response.GetResponse()
	result := make(map[string]int)
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		Logger.Error("Unmarshal error")
		return EZMQX_REST_ERROR
	}
	interval, exists := result[PAYLOAD_KEEPALIVE_INTERVAL]
	if !exists {
		Logger.Error("No keep alive interval key in json response")
		return EZMQX_REST_ERROR
	}
	if interval < 1 {
		Logger.Error("Invalid keepAlive interval")
		return EZMQX_REST_ERROR
	}
	Logger.Debug("Keep alive interval", zap.Int("Interval: ", interval))
	topicHandler := instance.topicHandler
	// fmt.println is used as logger is not supporting for atomic values
	fmt.Println("[parseTopicResponse] Current Keep Alive interval:", topicHandler.getKeepAliveInterval())
	if topicHandler.getKeepAliveInterval() < 0 {
		topicHandler.updateKeepAliveInterval(int64(interval))
	}
	if !topicHandler.isKeepAliveServiceStarted() {
		result := topicHandler.send(KEEPALIVE, "")
		if result != EZMQX_OK {
			Logger.Error("Topic handler send failed")
			return result
		}
	}
	return EZMQX_OK
}

func (instance *EZMQXPublisher) registerTopic(topic *EZMQXTopic) EZMQXErrorCode {
	isValid := validateTopic(topic.GetName())
	if false == isValid {
		Logger.Error("Topic validation failed")
		return EZMQX_INVALID_TOPIC
	}
	instance.topic = topic
	context := instance.context
	if !context.isCtxTnsEnabled() {
		return EZMQX_OK
	}
	// Send post request to TNS server
	jsonData := map[string]interface{}{PAYLOAD_NAME: topic.GetName(), PAYLOAD_DATAMODEL: topic.GetDataModel(), PAYLOAD_ENDPOINT: topic.GetEndPoint().ToString(), PAYLOAD_SECURED: topic.IsSecured()}
	payload := make(map[string]interface{})
	payload[PAYLOAD_TOPIC] = jsonData
	fmt.Println("TNS register topic payload: \n\n", payload)
	jsonValue, err := json.Marshal(payload)
	if err != nil {
		Logger.Error("TNS register topic: Json marshal failed")
		return EZMQX_REST_ERROR
	}
	client := GetRestFactory()
	topicURL := context.ctxGetTnsAddr() + PREFIX + TOPIC
	Logger.Debug("[TNS register topic] ", zap.String("Rest URL: ", string(topicURL)))
	response, error := client.Post(topicURL, jsonValue)
	if error != EZMQX_OK {
		Logger.Error("TNS register topic: Post request failed")
		return EZMQX_REST_ERROR
	}
	result := instance.parseTopicResponse(*response)
	if result != EZMQX_OK {
		Logger.Error("TNS register topic: Parse response failed")
		return result
	}
	//send a request to topic handler to add topic to topic list
	result = getTopicHandler().send(REGISTER, topic.GetName())
	if result != EZMQX_OK {
		Logger.Error("Topic handler send failed")
		return result
	}
	Logger.Debug("Sent request to topic handler to add topic to list: ", zap.String("Topic: ", topic.GetName()))
	return EZMQX_OK
}

func (instance *EZMQXPublisher) unRegisterTopic(topic *EZMQXTopic) EZMQXErrorCode {
	context := instance.context
	if !context.isCtxTnsEnabled() {
		return EZMQX_OK
	}
	topicURL := context.ctxGetTnsAddr() + PREFIX + TOPIC
	query := QUERY_NAME + topic.GetName()
	Logger.Debug("[TNS unregister topic]", zap.String("Rest URL: ", string(topicURL)))
	Logger.Debug("[TNS unregister topic]", zap.String("Query: ", string(query)))

	client := GetRestFactory()
	response, _ := client.Delete(topicURL+QUESTION_MARK+query, nil)
	Logger.Debug("[TNS unregister topic]", zap.Int("Status: ", response.GetStatusCode()))
	if response.GetStatusCode() != HTTP_OK {
		return EZMQX_REST_ERROR
	}

	//send request to topic handler to remove from topic list
	result := getTopicHandler().send(UNREGISTER, topic.GetName())
	if result != EZMQX_OK {
		Logger.Error("Topic handler send failed")
		return result
	}
	Logger.Debug("Sent request to topic handler to remove topic from list: ", zap.String("Topic: ", topic.GetName()))
	return EZMQX_OK
}

func (instance *EZMQXPublisher) terminate() EZMQXErrorCode {
	if false == atomic.CompareAndSwapUint32(&instance.status, INITIALIZED, TERMINATING) {
		Logger.Error("terminate failed : Not initialized")
		return EZMQX_UNKNOWN_STATE
	}
	context := instance.context
	if !context.isCtxStandAlone() {
		result := instance.context.releaseDynamicPort(instance.localPort)
		if result != EZMQX_OK {
			Logger.Error("Release dynamic port: failed")
		} else {
			Logger.Debug("Released local port")
		}
	}
	if context.isCtxTnsEnabled() {
		result := instance.unRegisterTopic(instance.topic)
		if result != EZMQX_OK {
			Logger.Error("Unregister topic: failed")
		} else {
			Logger.Debug("Unregistered topic on TNS")
		}
	}
	if nil != instance.ezmqPublisher {
		result := instance.ezmqPublisher.Stop()
		if result != EZMQX_OK {
			Logger.Error("Stop EZMQ publisher: failed")
			atomic.StoreUint32(&instance.status, INITIALIZED)
			return EZMQX_UNKNOWN_STATE
		}
		Logger.Debug("Stopped EZMQ publisher")
	}
	atomic.StoreUint32(&instance.status, CREATED)
	return EZMQX_OK
}

func (instance *EZMQXPublisher) isTerminated() bool {
	if atomic.LoadUint32(&instance.status) == CREATED {
		return true
	}
	return false
}

func (instance *EZMQXPublisher) getTopic() *EZMQXTopic {
	return instance.topic
}
