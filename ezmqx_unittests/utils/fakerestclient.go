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

package utils

import (
	"go/ezmqx"
	"time"
)

type FakeRestClient struct {
}

var restResponse = make(map[string][]byte)

func GetRestResponse(url string) []byte {
	return restResponse[url]
}

func SetRestResponse(url string, payload []byte) {
	restResponse[url] = payload
}

func GetFakeClient(timeout time.Duration) *FakeRestClient {
	var instance *FakeRestClient
	instance = &FakeRestClient{}
	return instance
}

func (instance *FakeRestClient) Get(url string) (*ezmqx.RestResponse, ezmqx.EZMQXErrorCode) {
	return ezmqx.GetRestResponse(200, restResponse[url]), ezmqx.EZMQX_OK
}

func (instance *FakeRestClient) Put(url string, data []byte) (*ezmqx.RestResponse, ezmqx.EZMQXErrorCode) {
	return ezmqx.GetRestResponse(200, restResponse[url]), ezmqx.EZMQX_OK
}

func (instance *FakeRestClient) Post(url string, data []byte) (*ezmqx.RestResponse, ezmqx.EZMQXErrorCode) {
	return ezmqx.GetRestResponse(201, restResponse[url]), ezmqx.EZMQX_OK
}

func (instance *FakeRestClient) Delete(url string, data []byte) (*ezmqx.RestResponse, ezmqx.EZMQXErrorCode) {
	return ezmqx.GetRestResponse(200, restResponse[url]), ezmqx.EZMQX_OK
}
