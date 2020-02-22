package results

import (
	"fmt"
	"strings"
)

type FieldTypeDesc struct {
	T        FieldType
	Key      string
	Title    string
	Examples string
	ToString func(interface{}) string
}

var (
	typeNames = [...]FieldTypeDesc{
		TypeArrayACL:         FieldTypeDesc{T: TypeArrayACL, Key: "[]acl", Title: "ACL Array", Examples: "", ToString: fmtString},
		TypeArrayBool:        FieldTypeDesc{T: TypeArrayBool, Key: "[]bool", Title: "Boolean Array", Examples: "", ToString: fmtString},
		TypeArrayBPChar:      FieldTypeDesc{T: TypeArrayBPChar, Key: "[]bpchar", Title: "Character String Array", Examples: "", ToString: fmtString},
		TypeArrayByteA:       FieldTypeDesc{T: TypeArrayByteA, Key: "[]bytea", Title: "Bytes Array", Examples: "", ToString: fmtString},
		TypeArrayCIDR:        FieldTypeDesc{T: TypeArrayCIDR, Key: "[]cidr", Title: "CIDR Array", Examples: "", ToString: fmtString},
		TypeArrayDate:        FieldTypeDesc{T: TypeArrayDate, Key: "[]date", Title: "Date Array", Examples: "", ToString: fmtString},
		TypeArrayFloat32:     FieldTypeDesc{T: TypeArrayFloat32, Key: "[]float32", Title: "32-bit Float Array", Examples: "", ToString: fmtString},
		TypeArrayFloat64:     FieldTypeDesc{T: TypeArrayFloat64, Key: "[]float64", Title: "64-bit Float Array", Examples: "", ToString: fmtString},
		TypeArrayInet:        FieldTypeDesc{T: TypeArrayInet, Key: "[]inet", Title: "Inet Address Array", Examples: "", ToString: fmtString},
		TypeArrayInt16:       FieldTypeDesc{T: TypeArrayInt16, Key: "[]int16", Title: "16-bit Integer Array", Examples: "", ToString: fmtString},
		TypeArrayInt32:       FieldTypeDesc{T: TypeArrayInt32, Key: "[]int32", Title: "32-bit Integer Array", Examples: "", ToString: fmtString},
		TypeArrayInt64:       FieldTypeDesc{T: TypeArrayInt64, Key: "[]int64", Title: "64-bit Integer Array", Examples: "", ToString: fmtString},
		TypeArrayJSON:        FieldTypeDesc{T: TypeArrayJSON, Key: "[]json", Title: "JSON Array", Examples: "", ToString: fmtString},
		TypeArrayNumeric:     FieldTypeDesc{T: TypeArrayNumeric, Key: "[]numeric", Title: "Numeric Array", Examples: "", ToString: fmtString},
		TypeArrayOID:         FieldTypeDesc{T: TypeArrayOID, Key: "[]oid", Title: "OID Array", Examples: "", ToString: fmtString},
		TypeArrayText:        FieldTypeDesc{T: TypeArrayText, Key: "[]text", Title: "Text Array", Examples: "", ToString: fmtString},
		TypeArrayTimestamp:   FieldTypeDesc{T: TypeArrayTimestamp, Key: "[]timestamp", Title: "Timestamp Array", Examples: "", ToString: fmtString},
		TypeArrayTimestampTZ: FieldTypeDesc{T: TypeArrayTimestampTZ, Key: "[]timestamptz", Title: "Zoned Timestamp Array", Examples: "", ToString: fmtString},
		TypeArrayUUID:        FieldTypeDesc{T: TypeArrayUUID, Key: "[]uuid", Title: "UUID Array", Examples: "", ToString: fmtString},
		TypeArrayVarchar:     FieldTypeDesc{T: TypeArrayVarchar, Key: "[]varchar", Title: "Varchar Array", Examples: "", ToString: fmtString},
		TypeArrayUnknown:     FieldTypeDesc{T: TypeArrayUnknown, Key: "[]unknown", Title: "Unknown Array", Examples: "", ToString: fmtString},
		TypeACL:              FieldTypeDesc{T: TypeACL, Key: "acl", Title: "ACL", Examples: "", ToString: fmtString},
		TypeBit:              FieldTypeDesc{T: TypeBit, Key: "bit", Title: "Bit", Examples: "", ToString: fmtString},
		TypeBitVarying:       FieldTypeDesc{T: TypeBitVarying, Key: "bit varying", Title: "Bit Varying", Examples: "", ToString: fmtString},
		TypeBool:             FieldTypeDesc{T: TypeBool, Key: "bool", Title: "Boolean", Examples: "", ToString: fmtString},
		TypeBox:              FieldTypeDesc{T: TypeBox, Key: "box", Title: "Box", Examples: "", ToString: fmtString},
		TypeBpchar:           FieldTypeDesc{T: TypeBpchar, Key: "bpchar", Title: "Character String", Examples: "", ToString: fmtString},
		TypeByteA:            FieldTypeDesc{T: TypeByteA, Key: "bytea", Title: "Bytes", Examples: "", ToString: fmtString},
		TypeChar:             FieldTypeDesc{T: TypeChar, Key: "char", Title: "Character", Examples: "", ToString: fmtString},
		TypeCID:              FieldTypeDesc{T: TypeCID, Key: "cid", Title: "CID", Examples: "", ToString: fmtString},
		TypeCIDR:             FieldTypeDesc{T: TypeCIDR, Key: "cidr", Title: "CIDR", Examples: "", ToString: fmtString},
		TypeCircle:           FieldTypeDesc{T: TypeCircle, Key: "circle", Title: "Circle", Examples: "", ToString: fmtString},
		TypeDate:             FieldTypeDesc{T: TypeDate, Key: "date", Title: "Date", Examples: "", ToString: fmtString},
		TypeDateRange:        FieldTypeDesc{T: TypeDateRange, Key: "daterange", Title: "Date Range", Examples: "", ToString: fmtString},
		TypeFloat32:          FieldTypeDesc{T: TypeFloat32, Key: "float32", Title: "32-bit Float", Examples: "", ToString: fmtString},
		TypeFloat64:          FieldTypeDesc{T: TypeFloat64, Key: "float64", Title: "64-bit Float", Examples: "", ToString: fmtString},
		TypeHStore:           FieldTypeDesc{T: TypeHStore, Key: "hstore", Title: "HStore", Examples: "", ToString: fmtString},
		TypeInet:             FieldTypeDesc{T: TypeInet, Key: "inet", Title: "Inet Address", Examples: "", ToString: fmtString},
		TypeInt8:             FieldTypeDesc{T: TypeInt8, Key: "int8", Title: "8-bit Integer", Examples: "", ToString: fmtString},
		TypeInt16:            FieldTypeDesc{T: TypeInt16, Key: "int16", Title: "16-bit Integer", Examples: "", ToString: fmtString},
		TypeInt32:            FieldTypeDesc{T: TypeInt32, Key: "int32", Title: "32-bit Integer", Examples: "", ToString: fmtString},
		TypeInt32Range:       FieldTypeDesc{T: TypeInt32Range, Key: "int32range", Title: "32-bit Integer Range", Examples: "", ToString: fmtString},
		TypeInt64:            FieldTypeDesc{T: TypeInt64, Key: "int64", Title: "64-bit Integer", Examples: "", ToString: fmtString},
		TypeInt64Range:       FieldTypeDesc{T: TypeInt64Range, Key: "int64range", Title: "64-bit Integer Range", Examples: "", ToString: fmtString},
		TypeInterval:         FieldTypeDesc{T: TypeInterval, Key: "interval", Title: "Interval", Examples: "", ToString: fmtString},
		TypeJSON:             FieldTypeDesc{T: TypeJSON, Key: "json", Title: "JSON", Examples: "", ToString: fmtString},
		TypeJSONB:            FieldTypeDesc{T: TypeJSONB, Key: "jsonb", Title: "JSONB", Examples: "", ToString: fmtString},
		TypeLine:             FieldTypeDesc{T: TypeLine, Key: "line", Title: "Line", Examples: "", ToString: fmtString},
		TypeLineSegment:      FieldTypeDesc{T: TypeLineSegment, Key: "linesegment", Title: "Line Segment", Examples: "", ToString: fmtString},
		TypeMacAddr:          FieldTypeDesc{T: TypeMacAddr, Key: "macaddr", Title: "MAC Address", Examples: "", ToString: fmtString},
		TypeMoney:            FieldTypeDesc{T: TypeMoney, Key: "money", Title: "Money", Examples: "", ToString: fmtString},
		TypeName:             FieldTypeDesc{T: TypeName, Key: "name", Title: "Name", Examples: "", ToString: fmtString},
		TypeNumeric:          FieldTypeDesc{T: TypeNumeric, Key: "numeric", Title: "Numeric", Examples: "", ToString: fmtTruncZeros},
		TypeNumRange:         FieldTypeDesc{T: TypeNumRange, Key: "numrange", Title: "Numeric Range", Examples: "", ToString: fmtString},
		TypeOID:              FieldTypeDesc{T: TypeOID, Key: "oid", Title: "OID", Examples: "", ToString: fmtString},
		TypePath:             FieldTypeDesc{T: TypePath, Key: "path", Title: "Path", Examples: "", ToString: fmtString},
		TypePoint:            FieldTypeDesc{T: TypePoint, Key: "point", Title: "Point", Examples: "", ToString: fmtString},
		TypePolygon:          FieldTypeDesc{T: TypePolygon, Key: "polygon", Title: "Polygon", Examples: "", ToString: fmtString},
		TypeRecord:           FieldTypeDesc{T: TypeRecord, Key: "record", Title: "Record", Examples: "", ToString: fmtString},
		TypeTID:              FieldTypeDesc{T: TypeTID, Key: "tid", Title: "TID", Examples: "", ToString: fmtString},
		TypeText:             FieldTypeDesc{T: TypeText, Key: "text", Title: "Text", Examples: "", ToString: fmtString},
		TypeTime:             FieldTypeDesc{T: TypeTime, Key: "time", Title: "Time", Examples: "", ToString: fmtString},
		TypeTimeTZ:           FieldTypeDesc{T: TypeTimeTZ, Key: "timetz", Title: "Zoned Time", Examples: "", ToString: fmtString},
		TypeTimestamp:        FieldTypeDesc{T: TypeTimestamp, Key: "timestamp", Title: "Timestamp", Examples: "", ToString: fmtString},
		TypeTsRange:          FieldTypeDesc{T: TypeTsRange, Key: "tsrange", Title: "Timestamp Range", Examples: "", ToString: fmtString},
		TypeTimestampTZ:      FieldTypeDesc{T: TypeTimestampTZ, Key: "timestamptz", Title: "Zoned Timestamp", Examples: "", ToString: fmtString},
		TypeTsTzRange:        FieldTypeDesc{T: TypeTsTzRange, Key: "tstzrange", Title: "Zoned Timestamp Range", Examples: "", ToString: fmtString},
		TypeTsQuery:          FieldTypeDesc{T: TypeTsQuery, Key: "tsquery", Title: "TS Query", Examples: "", ToString: fmtString},
		TypeTsVector:         FieldTypeDesc{T: TypeTsVector, Key: "tsvector", Title: "TS Vector", Examples: "", ToString: fmtString},
		TypeUUID:             FieldTypeDesc{T: TypeUUID, Key: "uuid", Title: "UUID", Examples: "", ToString: fmtString},
		TypeVarchar:          FieldTypeDesc{T: TypeVarchar, Key: "varchar", Title: "Varchar", Examples: "", ToString: fmtString},
		TypeXID:              FieldTypeDesc{T: TypeXID, Key: "xid", Title: "XID", Examples: "", ToString: fmtString},
		TypeXML:              FieldTypeDesc{T: TypeXML, Key: "xml", Title: "XML", Examples: "", ToString: fmtString},
		TypeUnknown:          FieldTypeDesc{T: TypeUnknown, Key: "unknown", Title: "Unknown", Examples: "", ToString: fmtString},
	}
)

func fmtString(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	if s == "<nil>" {
		s = "∅"
	}
	return s
}

func fmtTruncZeros(v interface{}) string {
	s := fmtString(v)
	for strings.HasSuffix(s, "0") {
		s = strings.TrimSuffix(s, "0")
	}
	for strings.HasSuffix(s, ".") {
		s = strings.TrimSuffix(s, ".")
	}
	return s
}
