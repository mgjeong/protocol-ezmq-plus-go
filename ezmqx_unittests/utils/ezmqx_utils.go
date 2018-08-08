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

package utils

import (
	"fmt"
	"go/aml"
	"go/ezmqx"
	"io/ioutil"
	"os"
	"time"
)

const ADDRESS = "127.0.0.1"
const TEST_LOCAL_HOST = "localhost"
const TNS_ADDRESS = "http://192.168.0.1:80/tns-server"

const PORT = 5562
const IP_PORT = "127.0.0.1:5562"
const TOPIC = "/topic"
const DATA_MODEL = "Robot_1.1"
const AML_FILE_PATH = "sample_data_model.aml"
const TNS_CONFIG_FILE_PATH = "tnsConf.json"
const NUMBER_OF_EVENTS = 5

const CONFIG_URL = "http://pharos-node:48098/api/v1/management/device/configuration"
const VALID_CONFIG_RESPONSE = `{ "properties": [{ "pinginterval": "10", "readOnly": false }, { "anchoraddress": "10.113.66.234", "readOnly": true }, { "deviceid": "71e8707c-f93b-4b77-a606-2860868429b7", "readOnly": true }, { "devicename": "MgmtServer", "readOnly": false }, { "nodeaddress": "10.113.66.234", "readOnly": true }, { "readOnly": false, "reverseproxy": { "enabled": true } }, { "anchorendpoint": "http://10.113.66.234:80/pharos-anchor/api/v1", "readOnly": true }, { "os": "linux", "readOnly": true }, { "platform": "Ubuntu 16.04.4 LTS", "readOnly": true }, { "processor": [{ "cpu": "0", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "1", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "2", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "3", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }], "readOnly": true }] }`
const INVALID_CONFIG_RESPONSE = `{ "property": [{ "pinginterval": "10", "readOnly": false }, { "anchoraddress": "10.113.66.234", "readOnly": true }, { "deviceid": "71e8707c-f93b-4b77-a606-2860868429b7", "readOnly": true }, { "devicename": "MgmtServer", "readOnly": false }, { "nodeaddress": "10.113.66.234", "readOnly": true }, { "readOnly": false, "reverseproxy": { "enabled": true } }, { "anchorendpoint": "http://10.113.66.234:80/pharos-anchor/api/v1", "readOnly": true }, { "os": "linux", "readOnly": true }, { "platform": "Ubuntu 16.04.4 LTS", "readOnly": true }, { "processor": [{ "cpu": "0", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "1", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "2", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }, { "cpu": "3", "modelname": "Intel(R) Core(TM) i5 CPU 750  @ 2.67GHz" }], "readOnly": true }] }`

const TNS_INFO_URL = "http://10.113.66.234:80/pharos-anchor/api/v1/search/nodes?imageName=docker.sec.samsung.net:5000/edge/system-tns-server-go/ubuntu_x86_64"
const VALID_TNS_INFO_RESPONSE = `{ "nodes": [ { "id": "node_id_sample", "ip": "192.168.0.1", "status": "connected", "apps": [ "app_id_sample1", "app_id_sample2" ], "config": { "properties": [ { "deviceid": "00000000-0000-0000-0000-000000000000", "readOnly": true }, { "devicename": "EdgeDevice", "readOnly": false }, { "pinginterval": "10", "readOnly": false }, { "os": "linux", "readOnly": true }, { "processor": [ { "cpu": "0", "modelname": "Intel(R) Core(TM) i7-2600 CPU @ 3.40GHz" } ], "readOnly": true }, { "platform": "Ubuntu 16.04.3 LTS", "readOnly": true }, { "reverseproxy": { "enabled": true }, "readOnly": true } ] } } ] }`

const RUNNING_APPS_URL = "http://pharos-node:48098/api/v1/management/apps"
const VALID_RUNNING_APPS_RESPONSE = `{ "apps": [{ "id": "103dd8cca769ce1aee520511f7379fdfe2a909cc", "state": "running" }] }`

const RUNNING_APP_INFO_URL = "http://pharos-node:48098/api/v1/management/apps/103dd8cca769ce1aee520511f7379fdfe2a909cc"

var hostName = ReadHostName(ezmqx.HOST_NAME_FILE_PATH)
var RUNNING_APP_INFO_RESPONSE = `{ "description": "services:  system-provisioning-director:    container_name: system-provisioning-director   image: docker.sec.samsung.net:5000/edge/system-provisioning-director/ubuntu_x86_64    labels:   - traefik.frontend.rule=PathPrefixStrip:/provisioning-service    - traefik.port=48198    network_mode: proxy    ports:    - 4000:4000    volumes:    - /system-provisioning-director/data/db:/data/dbversion: \"2\"", "images": [{ "name": "docker.sec.samsung.net:5000/edge/system-provisioning-director/ubuntu_x86_64" }], "services": [{ "cid": "` + hostName + `", "name": "system-provisioning-director", "ports": [{ "IP": "0.0.0.0", "PrivatePort": 4000, "PublicPort": 4000, "Type": "tcp" }], "state": { "exitcode": "0", "status": "running" } }], "state": "running" }`

const TOPIC_DISCOVERY_H_URL = "http://192.168.0.1:80/tns-server/api/v1/tns/topic?name=/topic&hierarchical=yes"
const TOPIC_DISCOVERY_URL = "http://192.168.0.1:80/tns-server/api/v1/tns/topic?name=/topic&hierarchical=no"
const VALID_TOPIC_DISCOVERY_RESPONSE = `{ "topics": [  {"name":  "topicName", "datamodel": "GTC_Robot_0.0.1", "endpoint": "localhost:5562" } ] }`
const INVALID_TOPIC_DISCOVERY_RESPONSE = `{ "topic": [  {"name":  "topicName", "datamodel": "GTC_Robot_0.0.1", "endpoint": "localhost:5562" } ] }`

const PUB_TNS_URL = "http://192.168.0.1:80/tns-server/api/v1/tns/topic"
const VALID_PUB_TNS_RESPONSE = `{ "ka_interval": 200 }`

const SUB_TOPIC_H_URL = "http://192.168.0.1:80/tns-server/api/v1/tns/topic?name=/topic&hierarchical=yes"
const SUB_TOPIC_RESPONSE = `{ "topics": [  {"name":  "/topic", "datamodel": "GTC_Robot_0.0.1", "endpoint": "localhost:5562" } ] }`
const SUB_TOPIC_URL = "http://192.168.0.1:80/tns-server/api/v1/tns/topic?name=/topic&hierarchical=no"

var Factory = ezmqx.GetRestFactory()

var Exit_Chan = make(chan bool, 1)

func Publish() {
	publisher, errorCode := ezmqx.GetAMLPublisher(TOPIC, ezmqx.AML_FILE_PATH, AML_FILE_PATH, PORT)
	if errorCode != ezmqx.EZMQX_OK {
		fmt.Println("Get publiser failed")
		os.Exit(-1)
	}
	amlObject := GetAMLObject()
	for i := 0; i < NUMBER_OF_EVENTS; i++ {
		time.Sleep(1500 * time.Millisecond)
		result := publisher.Publish(amlObject)
		if result != ezmqx.EZMQX_OK {
			fmt.Println("Error while publishing")
		}
		fmt.Println("Published event result:", result)
	}
	result := publisher.Terminate()
	if result != ezmqx.EZMQX_OK {
		fmt.Printf("Error while terminating publisher")
		os.Exit(-1)
	}
	Exit_Chan <- true
}

func GetAMLObject() *aml.AMLObject {
	// create "Model" data
	model, _ := aml.CreateAMLData()
	model.SetValueStr("ctname", "Model_107.113.97.248")
	model.SetValueStr("con", "SR-P7-970")

	// create "Sample" data
	axis, _ := aml.CreateAMLData()
	axis.SetValueStr("x", "20")
	axis.SetValueStr("y", "110")
	axis.SetValueStr("z", "80")

	info, _ := aml.CreateAMLData()
	info.SetValueStr("id", "f437da3b")
	info.SetValueAMLData("axis", axis)

	sample, _ := aml.CreateAMLData()
	sample.SetValueAMLData("info", info)
	appendix := [3]string{"935", "52303", "1442"}
	sample.SetValueStrArr("appendix", appendix[:])

	// set data to object
	amlObj, _ := aml.CreateAMLObject("Robot0001", time.Now().Format("20060102150405"))
	amlObj.AddData("Model", model)
	amlObj.AddData("Sample", sample)
	return amlObj
}

func ReadHostName(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	hostName := string(data)
	//remove trailing /n
	hostName = hostName[0 : len(hostName)-1]
	return hostName
}
