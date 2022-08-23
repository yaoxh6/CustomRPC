package rpc

import (
	"encoding/json"
	log "github.com/hyahm/golog"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
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

func DecodeArchiverWithTrace(rpcName string, d Codec, pak *transport.Package, args ...interface{}) error {
	tempArgs := []interface{}{}
	err := d.Decode(pak.Data, tempArgs)
	if err != nil {
		log.Fatal("Decode Failed")
		return err
	}
	args = tempArgs[1:]
	return nil
}