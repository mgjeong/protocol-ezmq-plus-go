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
	"testing"

	"container/list"
	"go/ezmqx"
)

const AML_FILE_PATH = "sample_data_model.aml"

func TestGetConfigInstance(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
}

func TestStartStandAloneMode(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(false, "")
	if ezmqx.EZMQX_OK != result {
		t.Errorf("StartStandAloneMode: Error")
	}
	instance.Reset()
}

/*func TestStartDockerMode(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartDockerMode()
	if ezmqx.EZMQX_OK != result {
		t.Errorf("StartDockerMode: Error")
	}
	instance.Reset()
}*/

func TestAddAmlModel(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	instance.StartStandAloneMode(false, "")
	amlFilePath := list.New()
	amlFilePath.PushBack(AML_FILE_PATH)
	_, result := instance.AddAmlModel(*amlFilePath)
	if ezmqx.EZMQX_OK != result {
		t.Errorf("AddAmlModel: Error")
	}
	instance.Reset()
}

func TestReset(t *testing.T) {
	var instance *ezmqx.EZMQXConfig = ezmqx.GetConfigInstance()
	result := instance.StartStandAloneMode(false, "")
	result = instance.Reset()
	if ezmqx.EZMQX_OK != result {
		t.Errorf("Reset [Standalone]: Error")
	}
	/*instance.StartDockerMode()
	result = instance.Reset()
	if ezmqx.EZMQX_OK != result {
		t.Errorf("Reset [Docker]: Error")
	}*/
}
