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

// Structure represents EZMQX topic.
type EZMQXTopic struct {
	name      string
	dataModel string
	endPoint  *EZMQXEndpoint
}

// Get EZMQX topic instance.
func GetEZMQXTopic(name string, dataModel string, endPoint *EZMQXEndpoint) *EZMQXTopic {
	var instance *EZMQXTopic
	instance = &EZMQXTopic{}
	instance.name = name
	instance.dataModel = dataModel
	instance.endPoint = endPoint
	return instance
}

// Get topic name.
func (topic *EZMQXTopic) GetName() string {
	return topic.name
}

// Get AML data model id.
func (topic *EZMQXTopic) GetDataModel() string {
	return topic.dataModel
}

// Get EZMQX end point.
func (topic *EZMQXTopic) GetEndPoint() *EZMQXEndpoint {
	return topic.endPoint
}
