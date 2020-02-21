package results

type FieldType uint8

const (
	TypeInvalid FieldType = iota
	TypeArrayACL
	TypeArrayBool
	TypeArrayBPChar
	TypeArrayByteA
	TypeArrayCIDR
	TypeArrayDate
	TypeArrayFloat32
	TypeArrayFloat64
	TypeArrayInet
	TypeArrayInt16
	TypeArrayInt32
	TypeArrayInt64
	TypeArrayJSON
	TypeArrayNumeric
	TypeArrayOID
	TypeArrayText
	TypeArrayTimestamp
	TypeArrayTimestampTZ
	TypeArrayUUID
	TypeArrayVarchar
	TypeArrayUnknown
	TypeACL
	TypeBit
	TypeBitVarying
	TypeBool
	TypeBox
	TypeBpchar
	TypeByteA
	TypeChar
	TypeCID
	TypeCIDR
	TypeCircle
	TypeDate
	TypeDateRange
	TypeFloat32
	TypeFloat64
	TypeHStore
	TypeInet
	TypeInt8
	TypeInt16
	TypeInt32
	TypeInt32Range
	TypeInt64
	TypeInt64Range
	TypeInterval
	TypeJSON
	TypeJSONB
	TypeLine
	TypeLineSegment
	TypeMacAddr
	TypeMoney
	TypeName
	TypeNumeric
	TypeNumRange
	TypeOID
	TypePath
	TypePoint
	TypePolygon
	TypeRecord
	TypeTID
	TypeText
	TypeTime
	TypeTimeTZ
	TypeTimestamp
	TypeTimestampTZ
	TypeTsRange
	TypeTsQuery
	TypeTsVector
	TypeTsTzRange
	TypeUnknown
	TypeUUID
	TypeVarchar
	TypeXID
	TypeXML
	endTypes
)

func (t FieldType) String() string {
	if t < endTypes {
		return typeNames[t].Key
	}
	return typeNames[TypeInvalid].Key
}

func (t FieldType) Desc() FieldTypeDesc {
	if t < endTypes {
		return typeNames[t]
	}
	return typeNames[TypeInvalid]
}

func (t FieldType) Valid() bool {
	return t > TypeInvalid && t < endTypes
}
