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

type EZMQXErrorCode int

// Constants represents EZMQX error codes.
const (
	EZMQX_OK                  = 0
	EZMQX_INVALID_PARAM       = 1
	EZMQX_INITIALIZED         = 2
	EZMQX_NOT_INITIALIZED     = 3
	EZMQX_TERMINATED          = 4
	EZMQX_UNKNOWN_STATE       = 5
	EZMQX_SERVICE_UNAVAILABLE = 6
	EZMQX_INVALID_TOPIC       = 7
	EZMQX_DUPLICATED_TOPIC    = 8
	EZMQX_UNKNOWN_TOPIC       = 9
	EZMQX_INVALID_ENDPOINT    = 10
	EZMQX_BROKEN_PAYLOAD      = 11
	EZMQX_REST_ERROR          = 12
	EZMQX_MAXIMUM_PORT_EXCEED = 13
	EZMQX_RELEASE_WRONG_PORT  = 14
	EZMQX_NO_TOPIC_MATCHED    = 15
	EZMQX_TNS_NOT_AVAILABLE   = 16
	EZMQX_UNKNOWN_AML_MODEL   = 17
	EZMQX_INVALID_AML_MODEL   = 18
	EZMQX_SESSION_UNAVAILABLE = 19
)
