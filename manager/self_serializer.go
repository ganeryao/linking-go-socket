package manager

import (
	lkCommon "github.com/ganeryao/linking-go-agile/common"
	"github.com/ganeryao/linking-go-agile/serialize/json"
	"github.com/ganeryao/linking-go-agile/serialize/protobuf"
)

// Serializer implements the serialize.Serializer interface
type SelfSerializer struct {
	protobuf *protobuf.Serializer
	json     *json.Serializer
}

// NewSerializer returns a new Serializer.
func NewSerializer() *SelfSerializer {
	return &SelfSerializer{protobuf: protobuf.NewSerializer(), json: json.NewSerializer()}
}

// Marshal returns the JSON encoding of v.
func (s *SelfSerializer) Marshal(v interface{}) ([]byte, error) {
	var protocolType = lkCommon.SelfRuntime.GetProtocolType()
	switch protocolType {
	case lkCommon.ProtocolProtobuf.String():
		return s.protobuf.Marshal(v)
	case lkCommon.ProtocolJson.String():
		return s.json.Marshal(v)
	default:
		return s.json.Marshal(v)
	}
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func (s *SelfSerializer) Unmarshal(data []byte, v interface{}) error {
	var protocolType = lkCommon.SelfRuntime.GetProtocolType()
	switch protocolType {
	case lkCommon.ProtocolProtobuf.String():
		return s.protobuf.Unmarshal(data, v)
	case lkCommon.ProtocolJson.String():
		return s.json.Unmarshal(data, v)
	default:
		return s.json.Unmarshal(data, v)
	}
}

// GetName returns the name of the serializer.
func (s *SelfSerializer) GetName() string {
	var protocolType = lkCommon.SelfRuntime.GetProtocolType()
	switch protocolType {
	case lkCommon.ProtocolProtobuf.String():
		return "protos"
	case lkCommon.ProtocolJson.String():
		return "json"
	default:
		return "json"
	}
}
