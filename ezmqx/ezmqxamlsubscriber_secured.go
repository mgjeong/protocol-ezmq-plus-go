// +build !unsecure

package ezmqx

import (
	"go.uber.org/zap"
)

// Get secured AML subscriber instance for given topic.
//
// Note:
// (1) Key should be 40-character string encoded in the Z85 encoding format
func GetSecuredAMLSubscriber(topic EZMQXTopic, serverPublicKey string, clientPublicKey string, clientSecretKey string, subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) (*EZMQXAMLSubscriber, EZMQXErrorCode) {
	if !topic.IsSecured() {
		return nil, EZMQX_INVALID_PARAM
	}
	instance := createAmlSubscriber(subCallback, errorCallback)
	result := instance.subscriber.storeSecuredTopics(topic, serverPublicKey, clientPublicKey, clientSecretKey)
	if result != EZMQX_OK {
		Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
		return nil, result
	}
	instance.isSecured = true
	return instance, result
}

// Get secured AML subscriber instance for given topic.
//
// Note:
// (1) Key should be 40-character string encoded in the Z85 encoding format
func GetSecuredAMLSubscriber1(topicKeyMap map[EZMQXTopic]string, clientPublicKey string, clientSecretKey string, subCallback EZMQXAmlSubCB, errorCallback EZMQXAmlErrorCB) (*EZMQXAMLSubscriber, EZMQXErrorCode) {
	for topic, _ := range topicKeyMap {
		if !topic.IsSecured() {
			return nil, EZMQX_INVALID_PARAM
		}
	}
	instance := createAmlSubscriber(subCallback, errorCallback)
	var result EZMQXErrorCode = EZMQX_INVALID_PARAM
	for topic, serverKey := range topicKeyMap {
		result = instance.subscriber.storeSecuredTopics(topic, serverKey, clientPublicKey, clientSecretKey)
		if result != EZMQX_OK {
			Logger.Error("Store topic failed", zap.Int("Error code:", int(result)))
			return nil, result
		}
	}
	instance.isSecured = true
	return instance, result
}
