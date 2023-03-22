package snippet

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

type TestStruct struct {
	Name string      `bson:"name"`
	Age  int         `bson:"age"`
	Key  UniformType `bson:"key"`
}

type UniformType struct {
	UT interface{}
}

func (ut *UniformType) UnmarshalBSONValue(t bsontype.Type, value []byte) error {
	if t == bsontype.Int32 {
		deci, _, ok := bsoncore.ReadInt32(value)
		if !ok {
			return fmt.Errorf("invalid Decimal128")
		}
		ut.UT = deci
	} else if t == bsontype.EmbeddedDocument {
		doc, _, ok := bsoncore.ReadDocument(value)
		if !ok {
			return fmt.Errorf("invalid embeded doc")
		}
		is := InnerStruct{}
		bson.Unmarshal(doc, &is)
		ut.UT = is
	}
	return nil
}

type InnerStruct struct {
	Sk  string `json:"sk"`
	Sk1 int32  `json:"sk1"`
}

type MStruct struct {
	Name string      `bson:"name"`
	Age  int32       `bson:"age"`
	Key  interface{} `bson:"key"`
}

func (ms *MStruct) UnmarshalBSON(b []byte) error {
	type MStructAlias MStruct
	err := bson.Unmarshal(b, (*MStructAlias)(ms))
	if err != nil {
		return err
	}

	switch t := ms.Key.(type) {
	case int32:
		ms.Key = t
	case primitive.D:
		d, err := bson.Marshal(t)
		if err != nil {
			return err
		}
		is := InnerStruct{}
		err = bson.Unmarshal(d, &is)
		if err != nil {
			return err
		}
		ms.Key = is
	}
	return nil
}
