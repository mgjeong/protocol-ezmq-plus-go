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
	context     *EZMQXContext
	initialized atomic.Value
	mutex       *sync.Mutex
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
		configInstance.initialized.Store(false)
		rand.Seed(time.Now().UnixNano())
		configInstance.mutex = &sync.Mutex{}
		InitLogger()
	}
	return configInstance
}

// Start/Configure EZMQX in docker mode.
// It works with Pharos system. In DockerMode, stack automatically use Tns service.
func (configInstance *EZMQXConfig) StartDockerMode() EZMQXErrorCode {
	configInstance.mutex.Lock()
	defer configInstance.mutex.Unlock()
	if true == configInstance.initialized.Load() {
		Logger.Debug("Already started")
		return EZMQX_INITIALIZED
	}
	result := configInstance.context.initializeDockerMode()
	if result != EZMQX_OK {
		Logger.Error("Initialize docker mode failed")
		return result
	}
	configInstance.initialized.Store(true)
	Logger.Debug("Started docker mode")
	return EZMQX_OK
}

// Start/Configure EZMQX in stand-alone mode.
// It works without pharos system.
func (configInstance *EZMQXConfig) StartStandAloneMode(useTns bool, tnsAddr string) EZMQXErrorCode {
	configInstance.mutex.Lock()
	defer configInstance.mutex.Unlock()
	if true == configInstance.initialized.Load() {
		Logger.Debug("Already started")
		return EZMQX_INITIALIZED
	}
	result := configInstance.context.initializeStandAloneMode(useTns, tnsAddr)
	if result != EZMQX_OK {
		Logger.Error("Initialize standalone mode failed")
		return result
	}
	configInstance.initialized.Store(true)
	Logger.Debug("Started standalone mode")
	return EZMQX_OK
}

// Add aml model file for publish or subscribe AML data.
func (configInstance *EZMQXConfig) AddAmlModel(amlFilePath list.List) (*list.List, EZMQXErrorCode) {
	configInstance.mutex.Lock()
	defer configInstance.mutex.Unlock()
	if false == configInstance.initialized.Load() {
		Logger.Error("Not initialized")
		return nil, EZMQX_NOT_INITIALIZED
	}
	return configInstance.context.addAmlRep(amlFilePath)
}

// Reset/Terminate EZMQX stack.
func (configInstance *EZMQXConfig) Reset() EZMQXErrorCode {
	configInstance.mutex.Lock()
	defer configInstance.mutex.Unlock()
	if false == configInstance.initialized.Load() {
		Logger.Error("Not initialized")
		return EZMQX_NOT_INITIALIZED
	}
	result := configInstance.context.terminate()
	if result != EZMQX_OK {
		Logger.Error("context terminate failed")
		return result
	}
	configInstance.initialized.Store(false)
	Logger.Debug("EZMQX reset done")
	return EZMQX_OK
}
