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
	"strconv"
	"strings"
)

//Structure represents EZMQX end point.
type EZMQXEndpoint struct {
	address string
	port    int
}

// Get EZMQX end point instance for the given address.
//
// Example: 127.0.0.0:4545
// If address contains only IP address then port will be initialized to -1.
func GetEZMQXEndPoint(address string) *EZMQXEndpoint {
	var instance *EZMQXEndpoint
	instance = &EZMQXEndpoint{}
	if false == strings.Contains(address, COLON) {
		instance.address = address
		instance.port = -1
		return instance
	}
	addr := strings.Split(address, COLON)
	instance.address = addr[0]
	instance.port, _ = strconv.Atoi(addr[1])
	return instance
}

// Get EZMQX end point instance for the given Ip and port.
func GetEZMQXEndPoint1(address string, port int) *EZMQXEndpoint {
	var instance *EZMQXEndpoint
	instance = &EZMQXEndpoint{}
	instance.address = address
	instance.port = port
	return instance
}

// Get address of end point.
func (endpoint *EZMQXEndpoint) GetAddr() string {
	return endpoint.address
}

// Get port of end point.
func (endpoint *EZMQXEndpoint) GetPort() int {
	return endpoint.port
}

// Get endpoint as string.
func (endpoint *EZMQXEndpoint) ToString() string {
	return endpoint.address + COLON + strconv.Itoa(endpoint.port)
}
