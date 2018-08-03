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
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type RestClient struct {
	client http.Client
}

func GetRestClient(timeout time.Duration) *RestClient {
	var instance *RestClient
	instance = &RestClient{}
	instance.client = http.Client{
		Timeout: timeout,
	}
	InitLogger()
	return instance
}

func (instance *RestClient) Get(url string) (*RestResponse, EZMQXErrorCode) {
	response, err := instance.client.Get(url)
	if err != nil {
		Logger.Error("HTTP request failed")
		return nil, EZMQX_REST_ERROR
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		Logger.Error("Failed to read response body")
		return nil, EZMQX_REST_ERROR
	}
	res := GetRestResponse(response.StatusCode, data)
	return res, EZMQX_OK
}

func (instance *RestClient) Put(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		Logger.Error("Form delete request failed")
		return nil, EZMQX_REST_ERROR
	}
	response, err := instance.client.Do(req)
	if err != nil {
		Logger.Error("Delete request failed")
		return nil, EZMQX_REST_ERROR
	}
	resData, error := ioutil.ReadAll(response.Body)
	if error != nil {
		Logger.Error("Failed to read response body")
		return nil, EZMQX_REST_ERROR
	}
	res := GetRestResponse(response.StatusCode, resData)
	return res, EZMQX_OK
}

func (instance *RestClient) Post(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	response, err := instance.client.Post(url, APPLICATION_JSON, bytes.NewBuffer(data))
	if err != nil {
		Logger.Error("Post request failed")
		return nil, EZMQX_REST_ERROR
	}
	resData, error := ioutil.ReadAll(response.Body)
	if error != nil {
		Logger.Error("Failed to read response body")
		return nil, EZMQX_REST_ERROR
	}
	res := GetRestResponse(response.StatusCode, resData)
	return res, EZMQX_OK
}

func (instance *RestClient) Delete(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		Logger.Error("Form delete request failed")
		return nil, EZMQX_REST_ERROR
	}
	response, err := instance.client.Do(req)
	if err != nil {
		Logger.Error("Delete request failed")
		return nil, EZMQX_REST_ERROR
	}
	res := GetRestResponse(response.StatusCode, nil)
	return res, EZMQX_OK
}
