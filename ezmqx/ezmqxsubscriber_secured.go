// +build !unsecure

package ezmqx

import (
	"go.uber.org/zap"
	"go/ezmq"
	"sync/atomic"
)

func (instance *EZMQXSubscriber) storeSecuredTopics(ezmqxTopic EZMQXTopic, serverPublicKey string, clientPublicKey string, clientSecretKey string) EZMQXErrorCode {
	context := instance.context
	if false == context.isCtxInitialized() {
		return EZMQX_NOT_INITIALIZED
	}
	var result EZMQXErrorCode
	//validate topic
	isValid := validateTopic(ezmqxTopic.GetName())
	if !isValid {
		Logger.Error("Invalid topic")
		return EZMQX_INVALID_TOPIC
	}
	instance.amlRepDic[ezmqxTopic.GetName()], result = context.getAmlRep(ezmqxTopic.GetDataModel())
	if result != EZMQX_OK {
		Logger.Error("getAmlRep failed", zap.Int("Error code:", int(result)))
		return result
	}
	result = instance.subscribeSecured(ezmqxTopic, serverPublicKey, clientPublicKey, clientSecretKey)
	if result != EZMQX_OK {
		Logger.Error("subscribe failed", zap.Int("Error code:", int(result)))
		return result
	}
	instance.storedTopics.PushBack(ezmqxTopic)
	atomic.StoreUint32(&instance.status, INITIALIZED)
	return EZMQX_OK
}

func (instance *EZMQXSubscriber) subscribeSecured(topic EZMQXTopic, serverPublicKey string, clientPublicKey string, clientSecretKey string) EZMQXErrorCode {
	if len(serverPublicKey) != KEY_LENGTH || len(clientPublicKey) != KEY_LENGTH || len(clientSecretKey) != KEY_LENGTH {
		return EZMQX_INVALID_PARAM
	}
	endPoint := topic.GetEndPoint()
	if nil == instance.ezmqSubscriber {
		result := instance.createSubscriber(endPoint)
		if result != EZMQX_OK {
			Logger.Error("Create subscriber failed", zap.Int("Error code:", int(result)))
			return result
		}
		//set server key
		ezmqResult := instance.ezmqSubscriber.SetServerPublicKey([]byte(serverPublicKey))
		if ezmqResult != ezmq.EZMQ_OK {
			Logger.Error("SetServerPublicKey failed", zap.Int("Error code:", int(result)))
			return EZMQX_UNKNOWN_STATE
		}
		//set client keys
		ezmqResult = instance.ezmqSubscriber.SetClientKeys([]byte(clientSecretKey), []byte(clientPublicKey))
		if ezmqResult != ezmq.EZMQ_OK {
			Logger.Error("SetClientKeys failed", zap.Int("Error code:", int(result)))
			return EZMQX_UNKNOWN_STATE
		}
		//start subscriber
		ezmqResult = instance.ezmqSubscriber.Start()
		if ezmqResult != ezmq.EZMQ_OK {
			Logger.Error("Start ezmq subscriber failed", zap.Int("Error code:", int(result)))
			return EZMQX_UNKNOWN_STATE
		}
		Logger.Debug("Started ezmq subscriber", zap.Int("Error code:", int(result)))
		//Subscribe
		errorCode := instance.ezmqSubscriber.SubscribeForTopic(topic.GetName())
		if errorCode != ezmq.EZMQ_OK {
			Logger.Error("Subscribe failed")
			return EZMQX_SESSION_UNAVAILABLE
		}
		Logger.Debug("Subscribed for topic", zap.String("Topic: ", topic.GetName()))
	} else {
		//set server key
		ezmqResult := instance.ezmqSubscriber.SetServerPublicKey([]byte(serverPublicKey))
		if ezmqResult != ezmq.EZMQ_OK {
			Logger.Error("SetServerPublicKey failed", zap.Int("Error code:", int(ezmqResult)))
			return EZMQX_UNKNOWN_STATE
		}
		errorCode := instance.ezmqSubscriber.SubscribeWithIPPort(endPoint.GetAddr(), endPoint.GetPort(), topic.GetName())
		if errorCode != ezmq.EZMQ_OK {
			Logger.Error("Subscribe with IP port failed")
			return EZMQX_SESSION_UNAVAILABLE
		}
		Logger.Debug("Subscribed for topic [With Ip and port]", zap.String("Topic: ", topic.GetName()))
	}
	return EZMQX_OK
}
