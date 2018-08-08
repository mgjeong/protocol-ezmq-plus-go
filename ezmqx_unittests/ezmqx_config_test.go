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
	"go/ezmqx"
	"go/ezmqx_unittests/utils"
	"testing"
)

func TestGetConfigInstance(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
}

func TestStartStandAloneMode(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	if ezmqx.EZMQX_OK != result {
		t.Errorf("StartStandAloneMode: Error")
	}
	instance.Reset()
}

func TestStartStandAloneModeWithTNS(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	if ezmqx.EZMQX_OK != result {
		t.Errorf("StartStandAloneMode: Error")
	}
	instance.Reset()
}

func TestMultipleStart(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(utils.TEST_LOCAL_HOST, true, "")
	if ezmqx.EZMQX_OK != result {
		t.Errorf("StartStandAloneMode: Error")
	}
	result = instance.StartStandAloneMode(utils.TEST_LOCAL_HOST, true, "")
	if ezmqx.EZMQX_OK == result {
		t.Errorf("StartStandAloneMode: Error")
	}
	instance.Reset()
}

func TestStartDockerMode(t *testing.T) {
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.CONFIG_URL, []byte(utils.VALID_CONFIG_RESPONSE))
	utils.SetRestResponse(utils.TNS_INFO_URL, []byte(utils.VALID_TNS_INFO_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APPS_URL, []byte(utils.VALID_RUNNING_APPS_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APP_INFO_URL, []byte(utils.RUNNING_APP_INFO_RESPONSE))
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartDockerMode(utils.TNS_CONFIG_FILE_PATH)
	if ezmqx.EZMQX_OK != result {
		t.Errorf("Start docker mode: Error")
	}
	instance.Reset()
}

func TestStartDockerModeNegative(t *testing.T) {
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.CONFIG_URL, []byte(utils.VALID_CONFIG_RESPONSE))
	utils.SetRestResponse(utils.TNS_INFO_URL, []byte(utils.VALID_TNS_INFO_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APPS_URL, []byte(utils.VALID_RUNNING_APPS_RESPONSE))
	utils.SetRestResponse(utils.RUNNING_APP_INFO_URL, []byte(utils.RUNNING_APP_INFO_RESPONSE))
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartDockerMode("")
	if result != ezmqx.EZMQX_UNKNOWN_STATE {
		t.Errorf("Start docker mode: Error")
	}
	instance.Reset()
}

func TestAddAmlModel(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	instance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	_, result := instance.AddAmlModel(*amlFilePath)
	if ezmqx.EZMQX_OK != result {
		t.Errorf("AddAmlModel: Error")
	}
	instance.Reset()
}

func TestAddAmlModelNegative(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	amlFilePath := list.New()
	amlFilePath.PushBack(utils.AML_FILE_PATH)
	_, result := instance.AddAmlModel(*amlFilePath)
	if ezmqx.EZMQX_OK == result {
		t.Errorf("AddAmlModel: Error")
	}
}

func TestReset(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	result = instance.Reset()
	if ezmqx.EZMQX_OK != result {
		t.Errorf("Reset [Standalone]: Error")
	}
	instance.StartDockerMode(utils.TNS_CONFIG_FILE_PATH)
	result = instance.Reset()
	if ezmqx.EZMQX_OK != result {
		t.Errorf("Reset [Docker]: Error")
	}
}

func TestResetNegative(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.Reset()
	if ezmqx.EZMQX_OK == result {
		t.Errorf("Reset: Error")
	}
}
