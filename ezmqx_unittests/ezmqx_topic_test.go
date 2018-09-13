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
)

func TestGetEZMQXTopic(t *testing.T) {
	var endPoint *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.IP_PORT)
	var instance *ezmqx.EZMQXTopic = ezmqx.GetEZMQXTopic(utils.TOPIC, utils.DATA_MODEL, false, endPoint)
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
}

func TestGetName(t *testing.T) {
	var endPoint *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.IP_PORT)
	var instance *ezmqx.EZMQXTopic = ezmqx.GetEZMQXTopic(utils.TOPIC, utils.DATA_MODEL, false, endPoint)
	if instance.GetName() != utils.TOPIC {
		t.Errorf("Error Topic name mismatch")
	}
}

func TestGetDataModel(t *testing.T) {
	var endPoint *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.IP_PORT)
	var instance *ezmqx.EZMQXTopic = ezmqx.GetEZMQXTopic(utils.TOPIC, utils.DATA_MODEL, false, endPoint)
	if instance.GetDataModel() != utils.DATA_MODEL {
		t.Errorf("Error data model mismatch")
	}
}

func TestGetEndPoint(t *testing.T) {
	var endPoint *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.IP_PORT)
	var instance *ezmqx.EZMQXTopic = ezmqx.GetEZMQXTopic(utils.TOPIC, utils.DATA_MODEL, false, endPoint)
	endPoint = instance.GetEndPoint()
	if nil == endPoint {
		t.Errorf("Error endpoint is NULL")
	}
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
	if endPoint.GetAddr() != utils.ADDRESS {
		t.Errorf("Error Address mismatch")
	}
	if endPoint.GetPort() != utils.PORT {
		t.Errorf("Error Address mismatch")
	}
	if endPoint.ToString() != utils.IP_PORT {
		t.Errorf("Error Address mismatch")
	}
}
