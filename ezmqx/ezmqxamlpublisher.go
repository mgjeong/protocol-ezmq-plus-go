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
	"go/aml"
	"go/ezmq"
)

// Structure represents EZMQX publisher.
type EZMQXAMLPublisher struct {
	publisher      *EZMQXPublisher
	representation *aml.Representation
	isSecured      bool
}

// Get EZMQX publisher instance.
func GetAMLPublisher(topic string, modelInfo EZMQXAmlModelInfo, modelId string, optionalPort int) (*EZMQXAMLPublisher, EZMQXErrorCode) {
	var instance *EZMQXAMLPublisher
	instance = &EZMQXAMLPublisher{}
	instance.publisher = getPublisher()
	result := instance.publisher.initialize(optionalPort)
	if result != EZMQX_OK {
		return nil, result
	}
	result = instance.registerTopic(topic, modelInfo, modelId, false)
	if result != EZMQX_OK {
		Logger.Error("Register topic failed, stopping ezmq publisher")
		instance.publisher.ezmqPublisher.Stop()
		return nil, result
	}
	instance.isSecured = false
	return instance, EZMQX_OK
}

// Publish AMLObject on the socket for subscribers.
func (instance *EZMQXAMLPublisher) Publish(object *aml.AMLObject) EZMQXErrorCode {
	publisher := instance.publisher
	if nil == publisher {
		Logger.Error("Publisher is null")
		return EZMQX_UNKNOWN_STATE
	}
	if publisher.context.isCtxTerminated() {
		Logger.Error("Context terminated")
		instance.Terminate()
		return EZMQX_TERMINATED
	}
	byteData, errorCode := instance.representation.DataToByte(object)
	if errorCode != aml.AML_OK {
		Logger.Error("AML DataToByte failed")
		return EZMQX_UNKNOWN_STATE
	}
	ezmqByteData := ezmq.EZMQByteData{byteData}
	ezmqPublisher := publisher.ezmqPublisher
	if nil == ezmqPublisher {
		Logger.Error("Ezmq Publisher failed")
		return EZMQX_UNKNOWN_STATE
	}
	result := ezmqPublisher.PublishOnTopic(publisher.topic.GetName(), ezmqByteData)
	if result != ezmq.EZMQ_OK {
		Logger.Error("Publish failed")
		return EZMQX_UNKNOWN_STATE
	}
	return EZMQX_OK
}

// Terminate EZMQX publisher.
func (instance *EZMQXAMLPublisher) Terminate() EZMQXErrorCode {
	publisher := instance.publisher
	if nil == publisher {
		return EZMQX_UNKNOWN_STATE
	}
	return publisher.terminate()
}

// Check whether publisher is terminated or not.
func (instance *EZMQXAMLPublisher) IsTerminated() (bool, EZMQXErrorCode) {
	publisher := instance.publisher
	if nil == publisher {
		return false, EZMQX_UNKNOWN_STATE
	}
	return publisher.isTerminated(), EZMQX_OK
}

// Get instance of Topic that used on this publisher.
func (instance *EZMQXAMLPublisher) GetTopic() (*EZMQXTopic, EZMQXErrorCode) {
	publisher := instance.publisher
	if nil == publisher {
		return nil, EZMQX_UNKNOWN_STATE
	}
	return publisher.getTopic(), EZMQX_OK
}

// Check whether publisher is secured or not.
func (instance *EZMQXAMLPublisher) IsSecured() (bool, EZMQXErrorCode) {
	return instance.isSecured, EZMQX_OK
}

func (instance *EZMQXAMLPublisher) registerTopic(topic string, modelInfo EZMQXAmlModelInfo, modelId string, isSecured bool) EZMQXErrorCode {
	var errorCode EZMQXErrorCode
	publisher := instance.publisher
	context := publisher.context
	if AML_MODEL_ID == modelInfo {
		instance.representation, errorCode = context.getAmlRep(modelId)
		if errorCode != EZMQX_OK {
			Logger.Error("Get aml representation failed [AML_MODEL_ID]")
			return errorCode
		}
	} else if AML_FILE_PATH == modelInfo {
		amlFilePath := list.New()
		amlFilePath.PushBack(modelId)
		idList, error := context.addAmlRep(*amlFilePath)
		if error != EZMQX_OK {
			Logger.Error("Add aml representation failed")
			return error
		}
		id := idList.Front().Value.(string)
		instance.representation, error = context.getAmlRep(id)
		if error != EZMQX_OK {
			Logger.Error("Get aml representation failed [AML_FILE_PATH]")
			return error
		}
	} else {
		Logger.Error("Unknown aml model info")
		return EZMQX_UNKNOWN_STATE
	}

	repId, amlCode := instance.representation.GetRepresentationId()
	if amlCode != aml.AML_OK {
		Logger.Error("Get representation ID failed")
		return EZMQX_UNKNOWN_STATE
	}
	hostEP, errorCode := context.getHostEp(publisher.localPort)
	if errorCode != EZMQX_OK {
		Logger.Error("Get hostEP failed")
		return EZMQX_UNKNOWN_STATE
	}
	ezmqxTopic := GetEZMQXTopic(topic, repId, isSecured, hostEP)
	return publisher.registerTopic(ezmqxTopic)
}
