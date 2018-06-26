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
	"container/list"
	"encoding/json"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"go.uber.org/zap"
	"go/ezmq"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type EZMQXTopicHandler struct {
	context            *zmq.Context
	topicServer        *zmq.Socket
	topicClient        *zmq.Socket
	poller             *zmq.Poller
	ezmqxContext       *EZMQXContext
	handlerAddress     string
	initialized        atomic.Value
	keepAliveInterval  atomic.Value
	isKeepAliveStarted atomic.Value
	isRoutineStarted   atomic.Value
	tnsAddress         string
	topicList          *list.List
	shutdownChan       chan string
	mutex              *sync.Mutex
}

var topicHandler *EZMQXTopicHandler
var handlerMutex = &sync.Mutex{}

func getTopicHandler() *EZMQXTopicHandler {
	handlerMutex.Lock()
	defer handlerMutex.Unlock()
	if nil == topicHandler {
		topicHandler = &EZMQXTopicHandler{}
		topicHandler.context = ezmq.GetInstance().GetContext()
		topicHandler.tnsAddress = getContextInstance().ctxGetTnsAddr()
		topicHandler.initialized.Store(false)
		var interval int64 = -1
		topicHandler.keepAliveInterval.Store(interval)
		topicHandler.isKeepAliveStarted.Store(false)
		topicHandler.isRoutineStarted.Store(false)
		topicHandler.topicList = list.New()
		topicHandler.shutdownChan = nil
		topicHandler.mutex = &sync.Mutex{}
	}
	return topicHandler
}

func (instance *EZMQXTopicHandler) initHandler() {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if true == instance.initialized.Load() {
		return
	}
	// Topic server
	if nil == instance.topicServer {
		address := getInProcUniqueAddress()
		Logger.Debug("initHandler", zap.String("Topic server address: ", address))
		instance.topicServer, _ = zmq.NewSocket(zmq.PAIR)
		if nil != instance.topicServer {
			instance.topicServer.Bind(address)
		}
		instance.handlerAddress = address
	}
	// Topic client
	if nil == instance.topicClient {
		address := instance.handlerAddress
		Logger.Debug("initHandler", zap.String("Topic client address: ", address))
		instance.topicClient, _ = zmq.NewSocket(zmq.PAIR)
		if nil != instance.topicClient {
			instance.topicClient.Connect(address)
		}
	}
	// Poller
	if nil == instance.poller {
		instance.poller = zmq.NewPoller()
		instance.poller.Add(instance.topicServer, zmq.POLLIN)
	}
	//call a go routine [new thread] for handler
	if false == instance.isRoutineStarted.Load() {
		topicHandler.isRoutineStarted.Store(true)
		go handleEvents(instance)
		Logger.Debug("Topic Handler thread started")
	}
	instance.initialized.Store(true)
}

func (instance *EZMQXTopicHandler) parseSocketData(topicServer *zmq.Socket) bool {
	var data string
	requestType, err := topicServer.Recv(0)
	if err != nil {
		return false
	}
	more, err := topicServer.GetRcvmore()
	if err != nil {
		return false
	}
	if true == more {
		data, err = topicServer.Recv(0)
		if err != nil {
			return false
		}
	}
	Logger.Debug("parseSocketData", zap.String("requestType: ", requestType))
	Logger.Debug("parseSocketData", zap.String("Data: ", data))
	if 0 == strings.Compare(requestType, SHUTDOWN) {
		return true
	} else if 0 == strings.Compare(requestType, REGISTER) {
		instance.addTopic(data)
	} else if 0 == strings.Compare(requestType, UNREGISTER) {
		instance.removeTopic(data)
	} else if 0 == strings.Compare(requestType, KEEPALIVE) {
		instance.isKeepAliveStarted.Store(true)
	} else {
		Logger.Error("Unknown request type")
	}
	return false
}

func handleEvents(instance *EZMQXTopicHandler) {
	var sockets []zmq.Polled
	var socket zmq.Polled
	var soc *zmq.Socket
	var err error
	var lastKeepAlive int64 = 0
	duration := time.Duration(1 * time.Second)
	for instance.poller != nil {
		duration = time.Duration(instance.keepAliveInterval.Load().(int64)) * time.Second
		fmt.Println("parseSocketData Duration is : ", duration)
		sockets, err = instance.poller.Poll(duration)
		Logger.Debug("Received register/unregister/keepalive/shutdown request")
		if err == nil {
			for _, socket = range sockets {
				switch soc = socket.Socket; soc {
				case instance.topicServer:
					isShutDown := instance.parseSocketData(instance.topicServer)
					if true == isShutDown {
						Logger.Debug("Received shut down request")
						goto End
					}
				default:
					Logger.Debug("Poller timeout occured")
				}
			}
		}
		if true == instance.isKeepAliveStarted.Load() {
			currentTime := time.Now().UnixNano() / int64(time.Millisecond)
			fmt.Println("parseSocketData current time: ", currentTime)
			difference := currentTime - lastKeepAlive
			fmt.Println("parseSocketData difference: ", difference)
			fmt.Println("parseSocketData instance.getKeepAliveInterval(): ", instance.getKeepAliveInterval())
			if difference >= instance.getKeepAliveInterval() {
				Logger.Debug("Sending keep alive request [timer expired]")
				instance.sendKeepAlive()
				lastKeepAlive = time.Now().UnixNano() / int64(time.Millisecond)
				fmt.Println("parseSocketData New Keep alive time: ", lastKeepAlive)
			}
		}
	}
End:
	if nil != instance.shutdownChan {
		Logger.Debug("[handleEvents] Go routine stopped: signaling channel")
		instance.shutdownChan <- "shutdown"
	}
}

func (instance *EZMQXTopicHandler) send(requestType string, payload string) EZMQXErrorCode {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	result, err := instance.topicClient.Send(requestType, zmq.SNDMORE)
	if err != nil {
		Logger.Error("Error while sending requestType", zap.Int("Result: ", result))
		return EZMQX_UNKNOWN_STATE
	}
	result, err = instance.topicClient.Send(payload, 0)
	if err != nil {
		Logger.Error("Error while sending payload", zap.Int("Result: ", result))
		return EZMQX_UNKNOWN_STATE
	}
	return EZMQX_OK
}

func (instance *EZMQXTopicHandler) addTopic(topic string) {
	if nil == instance.topicList {
		Logger.Error("Topic list is nil")
	}
	instance.topicList.PushBack(topic)
	Logger.Debug("Added topic to list", zap.String("Topic: ", topic))
}

func (instance *EZMQXTopicHandler) removeTopic(topic string) {
	topicList := instance.topicList
	if nil == topicList {
		Logger.Error("Topic list is nil")
	}
	for element := topicList.Front(); element != nil; element = element.Next() {
		if 0 == strings.Compare(element.Value.(string), topic) {
			topicList.Remove(element)
			Logger.Debug("Removed topic from list", zap.String("Topic: ", topic))
		}
	}
}

func (instance *EZMQXTopicHandler) sendKeepAlive() {
	instance.mutex.Lock()
	topicArray := make([]string, instance.topicList.Len())
	i := 0
	for topic := instance.topicList.Front(); topic != nil; topic = topic.Next() {
		topicArray[i] = topic.Value.(string)
		i++
	}
	instance.mutex.Unlock()
	payload := make(map[string]interface{})
	payload[PAYLOAD_TOPIC_KA] = topicArray
	fmt.Println("Payload to send: \n\n", payload)
	jsonPayload, error := json.Marshal(payload)
	if error != nil {
		Logger.Error("send Keep alive: json marshal failed")
		return
	}
	keepAliveURL := HTTP_PREFIX + instance.tnsAddress + COLON + TNS_KNOWN_PORT + PREFIX + TNS_KEEP_ALIVE
	Logger.Debug("[Send Keep Alive]", zap.String("Rest URL:", keepAliveURL))
	response, _ := http.Post(keepAliveURL, APPLICATION_JSON, bytes.NewBuffer(jsonPayload))
	Logger.Debug("[Send Keep Alive] ", zap.Int("Response Status code: ", response.StatusCode))
}

func (instance *EZMQXTopicHandler) terminateHandler() {
	instance.mutex.Lock()
	defer instance.mutex.Unlock()
	if false == instance.initialized.Load() {
		return
	}
	//shut down channel to stop go routine
	instance.shutdownChan = make(chan string)
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	instance.send(SHUTDOWN, "")
	select {
	case <-instance.shutdownChan:
		Logger.Debug("Received success shutdown signal")
	case <-timeout:
		Logger.Debug("Timeout occured for shutdown socket")
	}
	// remove socket from poller
	if nil != instance.poller {
		instance.poller.RemoveBySocket(instance.topicServer)
	}
	// close topicServer
	if nil != instance.topicServer {
		err := instance.topicServer.Close()
		if nil != err {
			Logger.Error("Error while closing topicServer socket")
		}
	}
	// close topicClient
	if nil != instance.topicClient {
		err := instance.topicClient.Close()
		if nil != err {
			Logger.Error("Error while closing topicClient socket")
		}
	}
	instance.poller = nil
	instance.topicServer = nil
	instance.topicClient = nil
	instance.topicList.Init()
	var interval int64 = -1
	instance.keepAliveInterval.Store(interval)
	instance.isKeepAliveStarted.Store(false)
	instance.isRoutineStarted.Store(false)
	instance.initialized.Store(false)
}

func (instance *EZMQXTopicHandler) updateKeepAliveInterval(interval int64) {
	instance.keepAliveInterval.Store(interval)
}

func (instance *EZMQXTopicHandler) getKeepAliveInterval() int64 {
	return instance.keepAliveInterval.Load().(int64)
}

func (instance *EZMQXTopicHandler) isKeepAliveServiceStarted() bool {
	return instance.isKeepAliveStarted.Load().(bool)
}

func getInProcUniqueAddress() string {
	return string(INPROC_PREFIX) + strconv.Itoa(rand.Intn(10000000))
}
