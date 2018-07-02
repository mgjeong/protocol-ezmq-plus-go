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
	"regexp"
	"strings"
)

const INPROC_PREFIX = "inproc://topicHandler"
const LOCAL_HOST = "localhost"
const LOCAL_PORT_START = 4000
const LOCAL_PORT_MAX = 100
const F_SLASH = "/"
const F_DOUBLE_SLASH = "//"
const TOPIC_PATTERN = "^(/)[a-zA-Z0-9-_./]+$"

//const TOPIC_WILD_CARD = "*"
//const TOPIC_WILD_PATTERN = "/*/";
const CREATED = 0
const INITIALIZING = 1
const INITIALIZED = 2
const TERMINATING = 3

func validateTopic(topic string) bool {
	if 0 == len(topic) {
		return false
	}
	if strings.Contains(topic, F_DOUBLE_SLASH) || strings.HasSuffix(topic, F_SLASH) {
		return false
	}
	result, _ := regexp.MatchString(TOPIC_PATTERN, topic)
	if false == result {
		return result
	}
	return true
}
