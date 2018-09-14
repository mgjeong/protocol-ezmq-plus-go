// +build !unsecure

package ezmqx

import (
	"go/ezmq"
	"sync/atomic"
)

func (instance *EZMQXPublisher) initializeSecured(optionalPort int, serverPrivateKey string) EZMQXErrorCode {
	if !instance.context.isCtxInitialized() {
		return EZMQX_NOT_INITIALIZED
	}
	if instance.context.isCtxStandAlone() {
		instance.localPort = optionalPort
	} else {
		var error EZMQXErrorCode
		instance.localPort, error = instance.context.assignDynamicPort()
		if error != EZMQX_OK {
			return error
		}
	}
	// create ezmq publisher
	instance.ezmqPublisher = ezmq.GetEZMQPublisher(instance.localPort, nil, nil, nil)
	if nil == instance.ezmqPublisher {
		Logger.Error("Could not create ezmq publisher")
		return EZMQX_UNKNOWN_STATE
	}
	//Set server key
	result := instance.ezmqPublisher.SetServerPrivateKey([]byte(serverPrivateKey))
	if result != ezmq.EZMQ_OK {
		return EZMQX_INVALID_PARAM
	}
	// Start ezmq publisher
	if ezmq.EZMQ_OK != instance.ezmqPublisher.Start() {
		Logger.Error("Could not start ezmq publisher")
		return EZMQX_UNKNOWN_STATE
	}
	// Init topic handler
	if instance.context.isCtxTnsEnabled() {
		instance.topicHandler = getTopicHandler()
		instance.topicHandler.initHandler()
		Logger.Debug("Initialized topic handler")
	}
	atomic.StoreUint32(&instance.status, INITIALIZED)
	return EZMQX_OK
}
