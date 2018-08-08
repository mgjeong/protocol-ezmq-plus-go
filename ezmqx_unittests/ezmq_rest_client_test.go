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

package ezmqx_unittests

import (
	"go/ezmqx"
	"testing"
)

const TEST_URL = "localhost:80/test"

func TestGetRestFactory(t *testing.T) {
	factory := ezmqx.GetRestFactory()
	if nil == factory {
		t.Errorf("Error get rest factory failed")
	}
}

func TestGetRequestN(t *testing.T) {
	factory := ezmqx.GetRestFactory()
	factory.SetFactory(ezmqx.RestClientFactory{})
	factory.Get(TEST_URL)
}

func TestPutRequestN(t *testing.T) {
	factory := ezmqx.GetRestFactory()
	factory.SetFactory(ezmqx.RestClientFactory{})
	factory.Put(TEST_URL, nil)
}

func TestPostRequestN(t *testing.T) {
	factory := ezmqx.GetRestFactory()
	factory.SetFactory(ezmqx.RestClientFactory{})
	factory.Post(TEST_URL, nil)
}

func TestDeleteRequestN(t *testing.T) {
	factory := ezmqx.GetRestFactory()
	factory.SetFactory(ezmqx.RestClientFactory{})
	factory.Delete(TEST_URL, nil)
}
