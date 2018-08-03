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

import "time"

var restFactoryInstance *RestFactory

type RestFactory struct {
	restInterface RestClientFactoryInterface
	timeout       time.Duration
}

func GetRestFactory() *RestFactory {
	if nil == restFactoryInstance {
		restFactoryInstance = &RestFactory{}
		restFactoryInstance.restInterface = RestClientFactory{}
		restFactoryInstance.timeout = time.Duration(CONNECTION_TIMEOUT * time.Second)
	}
	return restFactoryInstance
}

func (instance *RestFactory) SetFactory(factory RestClientFactoryInterface) {
	instance.restInterface = factory
}

func (instance *RestFactory) Get(url string) (*RestResponse, EZMQXErrorCode) {
	restClient := instance.restInterface.GetRestClient(instance.timeout)
	return restClient.Get(url)
}

func (instance *RestFactory) Put(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	restClient := instance.restInterface.GetRestClient(instance.timeout)
	return restClient.Put(url, data)
}

func (instance *RestFactory) Post(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	restClient := instance.restInterface.GetRestClient(instance.timeout)
	return restClient.Post(url, data)
}

func (instance *RestFactory) Post1(url string, data []byte, timeout time.Duration) (*RestResponse, EZMQXErrorCode) {
	restClient := instance.restInterface.GetRestClient(timeout)
	return restClient.Post(url, data)
}

func (instance *RestFactory) Delete(url string, data []byte) (*RestResponse, EZMQXErrorCode) {
	restClient := instance.restInterface.GetRestClient(instance.timeout)
	return restClient.Delete(url, data)
}
