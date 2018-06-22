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

func TestGetPublisherTest(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	idList, _ := configInstance.AddAmlModel(*amlFilePath)
	publisher, _ := ezmqx.GetPublisher(utils.TOPIC, ezmqx.AML_MODEL_ID, idList.Front().Value.(string), utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestGetPublisherTest1(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	publisher, _ := ezmqx.GetPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.FILE_PATH, utils.PORT)
	if nil == publisher {
		t.Errorf("publisher is nil")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestPublish(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	publisher, _ := ezmqx.GetPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.FILE_PATH, utils.PORT)
	result := publisher.Publish(utils.GetAMLObject())
	if result != ezmqx.EZMQX_OK {
		t.Errorf("publish failed")
	}
	publisher.Terminate()
	configInstance.Reset()
}

func TestTerminate(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(false, "")
	publisher, _ := ezmqx.GetPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.FILE_PATH, utils.PORT)
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
	configInstance.StartStandAloneMode(false, "")
	publisher, _ := ezmqx.GetPublisher(utils.TOPIC, ezmqx.AML_FILE_PATH, utils.FILE_PATH, utils.PORT)
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