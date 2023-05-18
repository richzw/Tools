package snippet

import (
	"fmt"
	"reflect"
)

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

// New Instances of Generic Types
// https://blog.openziti.io/golang-aha-moments-generics
type Example[T any] interface {
	*T
	Init(config map[string]interface{})
}

type ExampleFactory[T any, P Example[T]] struct {
	config map[string]interface{}
}

func (e *ExampleFactory[T, P]) Get() P {
	var result P = new(T)
	result.Init(e.config)
	return result
}

type MapInt map[string]*int

func (m MapInt) Init(config map[string]interface{}) {
	m = make(map[string]*int)
	for k, v := range config {
		(m)[k] = v.(*int)
	}
}

// another example
type GenericA[T any] interface {
	*T
}

type GenericB[U any, V GenericA[U]] struct{}

func (*GenericB[U, V]) Hello(v V) {
	fmt.Println("Hello B")
}

type GenericC[U any, V GenericA[U]] struct {
	b *GenericB[U, V]
}

func (*GenericC[U, V]) Hello(v V) {
	b := &GenericB[U, V]{}
	b.Hello(v)
	fmt.Println("Hello C")
}

func TestGenericV1() {
	a := 1
	c := &GenericC[int, *int]{}
	c.Hello(&a)

	// construct one instance through ExampleFactory
	factory := ExampleFactory[MapInt, *MapInt]{}

	factory.Get()
}
