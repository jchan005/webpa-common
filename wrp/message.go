package wrp

import (
	"fmt"
	"github.com/ugorji/go/codec"
	"io"
	"strconv"
)

var (
	wrpHandle = codec.MsgpackHandle{
		BasicHandle: codec.BasicHandle{
			TypeInfos: codec.NewTypeInfos([]string{"wrp"}),
		},
		WriteExt:    true,
		RawToString: true,
	}
)

// MessageType indicates the kind of WRP message
type MessageType int64

const (
	AuthMessageType                  = MessageType(2)
	SimpleRequestResponseMessageType = MessageType(3)
	SimpleEventMessageType           = MessageType(4)

	InvalidMessageTypeString = "!!INVALID!!"
)

var (
	messageTypeStrings = []string{
		InvalidMessageTypeString,
		InvalidMessageTypeString,
		"Auth",
		"SimpleRequestResponse",
		"SimpleEvent",
	}
)

func (mt MessageType) String() string {
	if int(mt) < len(messageTypeStrings) {
		return messageTypeStrings[mt]
	}

	return InvalidMessageTypeString
}

// Message represents a single WRP message.  The Type field determines how the other fields
// are interpreted.  For example, if the Type is AuthMessageType, then only Status will be set.
type Message struct {
	Type            MessageType `wrp:"msg_type" json:"-"`
	Status          *int64      `wrp:"status" json:"status,omitempty"`
	TransactionUUID string      `wrp:"transaction_uuid" json:"transaction_uuid,omitempty"`
	Source          string      `wrp:"source" json:"source,omitempty"`
	Destination     string      `wrp:"dest" json:"dest,omitempty"`
	Payload         []byte      `wrp:"payload" json:"payload,omitempty"`
}

// String returns a useful string representation of this message
func (m *Message) String() string {
	if m == nil {
		return "nil"
	}

	status := "nil"
	if m.Status != nil {
		status = strconv.FormatInt(*m.Status, 10)
	}

	return fmt.Sprintf(
		`{Type: %s, Status: %s, Source: %s, Destination: %s, Payload: %v}`,
		m.Type,
		status,
		m.Source,
		m.Destination,
		m.Payload,
	)
}

// Valid performs a basic validation check on a given message
func (m *Message) Valid() error {
	switch m.Type {
	case AuthMessageType:
		// nothing to validate here

	case SimpleRequestResponseMessageType:
		fallthrough

	case SimpleEventMessageType:
		if len(m.Destination) == 0 {
			return fmt.Errorf("Missing destination for message type: %s", m.Type)
		}

	default:
		return fmt.Errorf("Invalid message type: %d", m.Type)
	}

	return nil
}

// NewEncoder returns a codec.Encoder configured for WRP msgpack output.
func NewEncoder(output io.Writer) *codec.Encoder {
	return codec.NewEncoder(output, &wrpHandle)
}

// NewDecoder returns a codec.Decoder configured for WRP msgpack input.
func NewDecoder(input io.Reader) *codec.Decoder {
	return codec.NewDecoder(input, &wrpHandle)
}
