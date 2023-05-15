package snippet

func UnmarshalProtoMsgInGenericWay(body []byte, msg proto.Message) error {
	msgType := reflect.TypeOf(msg).Elem()
	msg = reflect.New(msgType).Interface().(proto.Message)
	return proto.Unmarshal(body, msg)
}

func Sample() {
	var msg T // Constrained to proto.Message

	// Peek the type inside T (as T= *SomeProtoMsgType)
	msgType := reflect.TypeOf(msg).Elem()

	// Make a new one, and throw it back into T
	msg = reflect.New(msgType).Interface().(T)

	errUnmarshal := proto.Unmarshal(body, msg)
}
