package csproto

import (
	"fmt"

	gogo "github.com/gogo/protobuf/proto"
	google "github.com/golang/protobuf/proto" //nolint: staticcheck // we're using this deprecated package intentionally
	googlev2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// HasExtension returns true if msg contains the specified proto2 extension field, delegating to the
// appropriate underlying Protobuf API based on the concrete type of msg.
func HasExtension(msg interface{}, ext interface{}) bool {
	switch MsgType(msg) {
	case MessageTypeGoogleV1:
		ed, ok := ext.(*google.ExtensionDesc)
		if !ok {
			return false
		}
		return google.HasExtension(msg.(google.Message), ed)
	case MessageTypeGoogle:
		et, ok := ext.(protoreflect.ExtensionType)
		if !ok {
			return false
		}
		return googlev2.HasExtension(msg.(googlev2.Message), et)
	case MessageTypeGogo:
		ed, ok := ext.(*gogo.ExtensionDesc)
		if !ok {
			return false
		}
		return gogo.HasExtension(msg.(gogo.Message), ed)
	default:
		return false
	}
}

// ClearExtension removes a proto2 extension field from msg, if it exists, delegating to the appropriate
// underlying Protobuf API based on the concrete type of msg.
//
// This function panics if the provded parameters are invalid, rather than returning an error, to be
// consistent with the signature of the ClearExtension() functions in the underlying Protobuf runtimes.
func ClearExtension(msg interface{}, ext interface{}) {
	switch MsgType(msg) {
	case MessageTypeGoogleV1:
		if ed, ok := ext.(*google.ExtensionDesc); ok {
			google.ClearExtension(msg.(google.Message), ed)
			return
		}
	case MessageTypeGoogle:
		if et, ok := ext.(protoreflect.ExtensionType); ok {
			googlev2.ClearExtension(msg.(googlev2.Message), et)
			return
		}
	case MessageTypeGogo:
		if ed, ok := ext.(*gogo.ExtensionDesc); ok {
			gogo.ClearExtension(msg.(gogo.Message), ed)
			return
		}
	default:
		panic(fmt.Sprintf("unsupported message type %T", msg))
	}
	// mismatched message and extension defintion types
	// - ex: a Google V2 message and a Gogo extension definition
	panic(fmt.Sprintf("invalid proto2 extension definition type %T for a message of type %T", ext, msg))
}

// GetExtension returns a proto2 extension field from msg, delegating to the appropriate underlying
// Protobuf API based on the concrete type of msg.
func GetExtension(msg interface{}, ext interface{}) (interface{}, error) {
	switch MsgType(msg) {
	case MessageTypeGoogleV1:
		ed, ok := ext.(*google.ExtensionDesc)
		if !ok {
			return nil, fmt.Errorf("invalid extension description type %T", ext)
		}
		return google.GetExtension(msg.(google.Message), ed)
	case MessageTypeGoogle:
		et, ok := ext.(protoreflect.ExtensionType)
		if !ok {
			return nil, fmt.Errorf("invalid extension type %T", ext)
		}
		return googlev2.GetExtension(msg.(googlev2.Message), et), nil
	case MessageTypeGogo:
		ed, ok := ext.(*gogo.ExtensionDesc)
		if !ok {
			return nil, fmt.Errorf("invalid extension description type %T", ext)
		}
		return gogo.GetExtension(msg.(gogo.Message), ed)
	default:
		return nil, fmt.Errorf("unsupported message type %T", ext)
	}
}

// SetExtension sets a proto2 extension field in msg to the provided value, delegating to the
// appropriate underlying Protobuf API based on the concrete type of msg.
func SetExtension(msg interface{}, ext interface{}, val interface{}) error {
	switch MsgType(msg) {
	case MessageTypeGoogleV1:
		ed, ok := ext.(*google.ExtensionDesc)
		if !ok {
			return fmt.Errorf("invalid extension description type %T", ext)
		}
		return google.SetExtension(msg.(google.Message), ed, val)
	case MessageTypeGoogle:
		et, ok := ext.(protoreflect.ExtensionType)
		if !ok {
			return fmt.Errorf("invalid extension type %T", ext)
		}
		googlev2.SetExtension(msg.(googlev2.Message), et, val)
		return nil
	case MessageTypeGogo:
		ed, ok := ext.(*gogo.ExtensionDesc)
		if !ok {
			return fmt.Errorf("invalid extension description type %T", ext)
		}
		return gogo.SetExtension(msg.(gogo.Message), ed, val)
	default:
		return fmt.Errorf("unsupported message type %T", ext)
	}
}

// ClearAllExtensions removes all proto2 extension fields from msg.
//
// This operation can be very inefficient, especially for Google V2 messages, so it is generally better
// to explicitly clear individual extension fields.
func ClearAllExtensions(msg interface{}) {
	switch MsgType(msg) {
	case MessageTypeGoogleV1:
		google.ClearAllExtensions(msg.(google.Message))
	case MessageTypeGoogle:
		// no ClearAllExtensions() API in Google V2 so we have to brute force it :(
		m := msg.(googlev2.Message)
		googlev2.RangeExtensions(m, func(xt protoreflect.ExtensionType, _ interface{}) bool {
			googlev2.ClearExtension(m, xt)
			return true
		})
	case MessageTypeGogo:
		gogo.ClearAllExtensions(msg.(gogo.Message))
	default:
		// no-op
	}
}
