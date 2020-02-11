package results

type FieldType uint8

const (
	TypeInvalid FieldType = iota
	TypeBool
	TypeTime
	TypeJSON
	TypeUUID
	TypeBytes
	TypeEnum
	TypeString
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt
	TypeInt64
	TypeUint8
	TypeUint16
	TypeUint32
	TypeUint
	TypeUint64
	TypeFloat32
	TypeFloat64
	endTypes
)

func (t FieldType) String() string {
	if t < endTypes {
		return typeNames[t]
	}
	return typeNames[TypeInvalid]
}

func (t FieldType) Numeric() bool {
	return t >= TypeInt8 && t < endTypes
}

func (t FieldType) Valid() bool {
	return t > TypeInvalid && t < endTypes
}

var (
	typeNames = [...]string{
		TypeInvalid: "invalid",
		TypeBool:    "bool",
		TypeTime:    "time.Time",
		TypeJSON:    "json.RawMessage",
		TypeUUID:    "[16]byte",
		TypeBytes:   "[]byte",
		TypeEnum:    "string",
		TypeString:  "string",
		TypeInt:     "int",
		TypeInt8:    "int8",
		TypeInt16:   "int16",
		TypeInt32:   "int32",
		TypeInt64:   "int64",
		TypeUint:    "uint",
		TypeUint8:   "uint8",
		TypeUint16:  "uint16",
		TypeUint32:  "uint32",
		TypeUint64:  "uint64",
		TypeFloat32: "float32",
		TypeFloat64: "float64",
	}
)
