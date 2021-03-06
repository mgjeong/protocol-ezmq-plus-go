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

// ezmqx package which provides simplified APIs for publisher and subscriber.
package ezmqx

import (
	"container/list"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Structure represents EZMQX configuration.
type EZMQXConfig struct {
	context *EZMQXContext
	status  uint32
}

var configInstance *EZMQXConfig
var configMutex = &sync.Mutex{}

// Get EZMQX Config instance.
func GetConfigInstance() *EZMQXConfig {
	configMutex.Lock()
	defer configMutex.Unlock()
	if nil == configInstance {
		configInstance = &EZMQXConfig{}
		configInstance.context = getContextInstance()
		configInstance.status = CREATED
		rand.Seed(time.Now().UnixNano())
		InitLogger()
		factory := GetRestFactory()
		factory.SetFactory(RestClientFactory{})
	}
	return configInstance
}

// Start/Configure EZMQX in docker mode.
// It works with Pharos system. In DockerMode, stack automatically use Tns service.
func (configInstance *EZMQXConfig) StartDockerMode(tnsConfPath string) EZMQXErrorCode {
	if false == atomic.CompareAndSwapUint32(&configInstance.status, CREATED, INITIALIZING) {
		Logger.Error("Initialize docker mode failed: Invalid state")
		return EZMQX_UNKNOWN_STATE
	}
	result := configInstance.context.initializeDockerMode(tnsConfPath)
	if result != EZMQX_OK {
		Logger.Error("Initialize docker mode failed")
		atomic.StoreUint32(&configInstance.status, CREATED)
		return result
	}
	atomic.StoreUint32(&configInstance.status, INITIALIZED)
	Logger.Debug("Started docker mode")
	return EZMQX_OK
}

// Start/Configure EZMQX in stand-alone mode.
// It works without pharos system.
// Note: TNS address should be complete Rest address of TNS.
func (configInstance *EZMQXConfig) StartStandAloneMode(hostAddr string, useTns bool, tnsAddr string) EZMQXErrorCode {
	if false == atomic.CompareAndSwapUint32(&configInstance.status, CREATED, INITIALIZING) {
		Logger.Error("Initialize standalone mode failed: Invalid state")
		return EZMQX_UNKNOWN_STATE
	}
	result := configInstance.context.initializeStandAloneMode(hostAddr, useTns, tnsAddr)
	if result != EZMQX_OK {
		Logger.Error("Initialize standalone mode failed")
		atomic.StoreUint32(&configInstance.status, CREATED)
		return result
	}
	atomic.StoreUint32(&configInstance.status, INITIALIZED)
	Logger.Debug("Started standalone mode")
	return EZMQX_OK
}

// Add aml model file for publish or subscribe AML data.
func (configInstance *EZMQXConfig) AddAmlModel(amlFilePath list.List) (*list.List, EZMQXErrorCode) {
	if atomic.LoadUint32(&configInstance.status) != INITIALIZED {
		Logger.Error("Not initialized")
		return nil, EZMQX_NOT_INITIALIZED
	}
	return configInstance.context.addAmlRep(amlFilePath)
}

// Reset/Terminate EZMQX stack.
func (configInstance *EZMQXConfig) Reset() EZMQXErrorCode {
	if false == atomic.CompareAndSwapUint32(&configInstance.status, INITIALIZED, TERMINATING) {
		Logger.Error("Reset failed: invalid state")
		return EZMQX_UNKNOWN_STATE
	}
	result := configInstance.context.terminate()
	if result != EZMQX_OK {
		Logger.Error("context terminate failed")
		atomic.StoreUint32(&configInstance.status, INITIALIZED)
		return result
	}
	atomic.StoreUint32(&configInstance.status, CREATED)
	Logger.Debug("EZMQX reset done")
	return EZMQX_OK
}
