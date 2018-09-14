// +build !unsecure

package ezmqx

// Get Secured EZMQX publisher instance.
//
// Note:
// (1) Key should be 40-character string encoded in the Z85 encoding format
func GetSecuredAMLPublisher(topic string, serverPrivateKey string, modelInfo EZMQXAmlModelInfo, modelId string, optionalPort int) (*EZMQXAMLPublisher, EZMQXErrorCode) {
	var instance *EZMQXAMLPublisher
	instance = &EZMQXAMLPublisher{}
	instance.publisher = getPublisher()
	result := instance.publisher.initializeSecured(optionalPort, serverPrivateKey)
	if result != EZMQX_OK {
		return nil, result
	}
	result = instance.registerTopic(topic, modelInfo, modelId, true)
	if result != EZMQX_OK {
		Logger.Error("Register topic failed, stopping ezmq publisher")
		instance.publisher.ezmqPublisher.Stop()
		return nil, result
	}
	instance.isSecured = true
	return instance, EZMQX_OK
}
