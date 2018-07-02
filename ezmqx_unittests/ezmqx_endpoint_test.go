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

func TestGetEZMQXEndPoint(t *testing.T) {
	var instance *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.IP_PORT)
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
	if instance.GetAddr() != utils.ADDRESS {
		t.Errorf("Error Address mismatch")
	}
	if instance.GetPort() != utils.PORT {
		t.Errorf("Error Address mismatch")
	}
	if instance.ToString() != utils.IP_PORT {
		t.Errorf("Error Address mismatch")
	}
}

func TestGetEZMQXEndPoint1(t *testing.T) {
	var instance *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint(utils.ADDRESS)
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
	if instance.GetAddr() != utils.ADDRESS {
		t.Errorf("Error Address mismatch")
	}
	if instance.GetPort() != -1 {
		t.Errorf("Error Address mismatch")
	}
}

func TestGetEZMQXEndPoint2(t *testing.T) {
	var instance *ezmqx.EZMQXEndpoint = ezmqx.GetEZMQXEndPoint1(utils.ADDRESS, utils.PORT)
	if nil == instance {
		t.Errorf("Error config instance is NULL")
	}
	if instance.GetAddr() != utils.ADDRESS {
		t.Errorf("Error Address mismatch")
	}
	if instance.GetPort() != utils.PORT {
		t.Errorf("Error Address mismatch")
	}
	if instance.ToString() != utils.IP_PORT {
		t.Errorf("Error Address mismatch")
	}
}