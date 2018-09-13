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
	"go/ezmqx"
	"go/ezmqx_unittests/utils"
	"testing"

	"container/list"
)

func TestGetPublisherStandAlone(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestGetPublisherStandAlone1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestGetPublisherStandAlone2(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)

	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.PUB_TNS_URL, []byte(utils.VALID_PUB_TNS_RESPONSE))

	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestGetPublisherStandAlone3(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)

	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.PUB_TNS_URL, []byte(utils.VALID_PUB_TNS_RESPONSE))

	_, result := ezmqx.GetAMLPublisher("", ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Get publisher failed")
	}
	configInstance.Reset()
}

func TestGetPublisherStandAlone4(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)

	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.PUB_TNS_URL, []byte(utils.VALID_PUB_TNS_RESPONSE))

	_, result := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, "", utils.PORT)
	if result != ezmqx.EZMQX_INVALID_AML_MODEL {
		t.Errorf("Get publisher failed")
	}
	_, result = ezmqx.GetAMLPublisher(utils.TOPIC, 5, utils.AML_FILE_PATH, utils.PORT)
	if result != ezmqx.EZMQX_UNKNOWN_STATE {
		t.Errorf("Get publisher failed")
	}
	_, result = ezmqx.GetAMLPublisher(utils.TOPIC, 5, "", utils.PORT)
	if result != ezmqx.EZMQX_UNKNOWN_STATE {
		t.Errorf("Get publisher failed")
	}
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	_, _ = configInstance.AddAmlModel(*amlFilePath)
	_, result = ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_MODEL_ID, "", utils.PORT)
	if result != ezmqx.EZMQX_UNKNOWN_AML_MODEL {
		t.Errorf("publisher is nil")
	}
	configInstance.Reset()
}

func TestGetSecuredPublisher(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	publisher, errorCode := ezmqx.GetSecuredAMLPublisher(utils.TOPIC, utils.SERVER_SECRET_KEY, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if errorCode != ezmqx.EZMQX_OK {
		t.Errorf("GetSecuredAMLPublisher failed")
	}
	isSecured, _ := publisher.IsSecured()
	if !isSecured {
		t.Errorf("publisher is secured failed")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestGetSecuredPublisherNegative(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	//Invalid key
	_, errorCode := ezmqx.GetSecuredAMLPublisher(utils.TOPIC, " ", ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if errorCode != ezmqx.EZMQX_INVALID_PARAM {
		t.Errorf("Get publisher failed")
	}
	//Invalid topic
	_, errorCode = ezmqx.GetSecuredAMLPublisher("topic", utils.SERVER_SECRET_KEY, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if errorCode != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Get publisher failed")
	}
	configInstance.Reset()
	//Without config
	_, errorCode = ezmqx.GetSecuredAMLPublisher(utils.TOPIC, utils.SERVER_SECRET_KEY, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if errorCode != ezmqx.EZMQX_NOT_INITIALIZED {
		t.Errorf("Get publisher failed")
	}
	//With tns
	configInstance = ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, true, "")
	amlFilePath = list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	idList, _ = configInstance.AddAmlModel(*amlFilePath)
	_, errorCode = ezmqx.GetSecuredAMLPublisher(utils.TOPIC, utils.SERVER_SECRET_KEY, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), 5566)
	if ezmqx.EZMQX_OK == errorCode {
		t.Errorf("GetSecuredAMLPublisher wrong error code")
	}
	configInstance.Reset()
}

func TestIsSecuredPublisher(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	isSecured, _ := publisher.IsSecured()
	if isSecured {
		t.Errorf("publisher is secured failed")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestStandAlonePublish(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	result := publisher.Publish(utils.GetAMLObject())
	if result != ezmqx.EZMQX_OK {
		t.Errorf("publish failed")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestStandAlonePublishNegative(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	result := publisher.Publish(nil)
	if result != ezmqx.EZMQX_UNKNOWN_STATE {
		t.Errorf("publish failed")
	}
	publisher.Terminate()
	configInstance.Reset()
	result = publisher.Publish(nil)
	if result != ezmqx.EZMQX_TERMINATED {
		t.Errorf("publish failed")
	}
}

func TestDockerModePublish(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.CONFIG_URL, []byte(utils.VALID_CONFIG_RESPONSE))
	utils.SetRestResponse(utils.TNS_INFO_URL, []byte(utils.VALID_TNS_INFO_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APPS_URL, []byte(utils.VALID_RUNNING_APPS_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APP_INFO_URL, []byte(utils.RUNNING_APP_INFO_RESPONSE))
	utils.SetRestResponse(utils.PUB_TNS_URL, []byte(utils.VALID_PUB_TNS_RESPONSE))

	configInstance.StartDockerMode(utils.TNS_CONFIG_FILE_PATH)
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestTerminate(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	isTerminated, _ := publisher.IsTerminated()
	if true == isTerminated {
		t.Errorf("Publisher terminated")
	}
	result := publisher.Terminate()
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Terminate failed")
	}
	isTerminated, _ = publisher.IsTerminated()
	if false == isTerminated {
		t.Errorf("Terminate failed")
	}
	configInstance.Reset()
}

func TestGetTopic(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	publisher, _ := ezmqx.GetAMLPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.AML_FILE_PATH, utils.PORT)
	result := publisher.Publish(utils.GetAMLObject())
	if result != ezmqx.EZMQX_OK {
		t.Errorf("publish failed")
	}
	topic, _ := publisher.GetTopic()
	if topic.GetName() != utils.TOPIC {
		t.Errorf("Topic mismatch")
	}
	publisher.Terminate()
	configInstance.Reset()
}
