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

func TestGetEZMQXTopicDiscovery(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.TEST_LOCAL_HOST, false, "")
	_, result := ezmqx.GetEZMQXTopicDiscovery()
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Error get EZMQX topic discovery failed")
	}
	configInstance.Reset()
}

func TestGetTopicDiscoveryNegative(t *testing.T) {
	_, result := ezmqx.GetEZMQXTopicDiscovery()
	if result == ezmqx.EZMQX_OK {
		t.Errorf("Error EZMQX topic query failed")
	}
}

func TestQuery(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()

	//Set fake rest client
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.TOPIC_DISCOVERY_URL, []byte(utils.VALID_TOPIC_DISCOVERY_RESPONSE))

	_, result := topicDiscovery.Query(utils.TOPIC)
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Error EZMQX topic query failed")
	}
	configInstance.Reset()
}

func TestQueryNegative(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()
	configInstance.Reset()
	_, result := topicDiscovery.Query(utils.TOPIC)
	if result != ezmqx.EZMQX_TERMINATED {
		t.Errorf("Error EZMQX topic query failed")
	}

	configInstance = ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, false, "")
	topicDiscovery, _ = ezmqx.GetEZMQXTopicDiscovery()
	_, result = topicDiscovery.Query(utils.TOPIC)
	if result != ezmqx.EZMQX_TNS_NOT_AVAILABLE {
		t.Errorf("Error EZMQX topic query failed")
	}
	configInstance.Reset()
}

func TestQueryNegative2(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()

	//Set fake rest client
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.TOPIC_DISCOVERY_URL, []byte(utils.INVALID_TOPIC_DISCOVERY_RESPONSE))

	_, result := topicDiscovery.Query(utils.TOPIC)
	if result != ezmqx.EZMQX_REST_ERROR {
		t.Errorf("Error EZMQX topic query failed")
	}
	configInstance.Reset()
}

func TestHierarchicalQuery(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()

	//Set fake rest client
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.TOPIC_DISCOVERY_H_URL, []byte(utils.VALID_TOPIC_DISCOVERY_RESPONSE))

	_, result := topicDiscovery.HierarchicalQuery(utils.TOPIC)
	if result != ezmqx.EZMQX_OK {
		t.Errorf("Error EZMQX topic query failed")
	}
	configInstance.Reset()
}

func TestQueryTopicValidation(t *testing.T) {
	configInstance := ezmqx.GetConfigInstance()
	configInstance.StartStandAloneMode(utils.ADDRESS, true, utils.TNS_ADDRESS)
	topicDiscovery, _ := ezmqx.GetEZMQXTopicDiscovery()

	//Set fake rest client
	utils.Factory.SetFactory(utils.FakeRestClientFactory{})
	utils.SetRestResponse(utils.TOPIC_DISCOVERY_URL, []byte(utils.VALID_TOPIC_DISCOVERY_RESPONSE))

	// Valid topics
	topic := "/sdf/sdf/sdfgfg/12"
	_, result := topicDiscovery.Query(topic)
	if result == ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdfsdf/123sdfs/t312xsdf/213lkj_"
	_, result = topicDiscovery.Query(topic)
	if result == ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/-0/ssd............fdfadsf/fdsg-/-0-"
	_, result = topicDiscovery.Query(topic)
	if result == ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdf/sdf/sdfgfg"
	_, result = topicDiscovery.Query(topic)
	if result == ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/123/sdafdsaf/44___kk/2232/abicls"
	_, result = topicDiscovery.Query(topic)
	if result == ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	// Invalid topics
	topic = "//////)///(///////////"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "//3/1/2/1/2/3"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "312123/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/fds+dsfg-23-/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/312?_!_12--3/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/312?_!_12--3/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/312_!_123/sda"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/        /312123"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "-=0/ssdfdfadsf/fdsg-/-0-"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdfsdf*/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdf/sdf/sdfgfg/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "sdf/sdf/sdf/sdfgfg/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdf/sdf/sdf gfg/12"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdfsdf/123sdfs/t312xsdf*/213lkj_+/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "sdfsdf/123sdfs/t312xsdf*/213lkj_+/"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/sdfsdf/123sdfs//t312xsdf*/213lkj_+"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = ""
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	topic = "/#topic"
	_, result = topicDiscovery.Query(topic)
	if result != ezmqx.EZMQX_INVALID_TOPIC {
		t.Errorf("Error EZMQX topic validation failed")
	}

	configInstance.Reset()
}
