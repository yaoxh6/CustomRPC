package rpc

import (
	"encoding/json"
)

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte, interface{}) error
}

type JsonCodec struct {
}
func (c *JsonCodec)Encode(data interface{}) ([]byte, error){
	 return json.Marshal(data)
}

func (c *JsonCodec)Decode(data []byte, v interface{}) error{
	return json.Unmarshal(data, v)
}