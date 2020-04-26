package results

import (
	"strings"

	"logur.dev/logur"
)

func FieldTypeForName(logger logur.LoggerFacade, name string, t string) FieldType {
	t = strings.ToLower(t)
	if strings.HasPrefix(t, "_") {
		return arrayTypeFor(logger, name, strings.TrimPrefix(t, "_"))
	}
	r := TypeUnknown
	switch t {
	case "aclitem":
		r = TypeACL
	case "bit":
		r = TypeBit
	case "varbit", "bit varying":
		r = TypeBitVarying
	case "bool", "boolean":
		r = TypeBool
	case "box":
		r = TypeBox
	case "bpchar":
		r = TypeBpchar
	case "bytea":
		r = TypeByteA
	case "char":
		r = TypeChar
	case "cid":
		r = TypeCID
	case "cidr":
		r = TypeCIDR
	case "circle":
		r = TypeCircle
	case "date":
		r = TypeDate
	case "daterange":
		r = TypeDateRange
	case "float4", "real", "float":
		r = TypeFloat32
	case "float8", "double precision", "double":
		r = TypeFloat64
	case "hstore":
		r = TypeHStore
	case "inet":
		r = TypeInet
	case "int1", "tinyint":
		r = TypeInt8
	case "int2", "smallint":
		r = TypeInt16
	case "int4", "integer", "int", "mediumint":
		r = TypeInt32
	case "int4range":
		r = TypeInt32Range
	case "int8", "bigint":
		r = TypeInt64
	case "int8range":
		r = TypeInt64Range
	case "interval":
		r = TypeInterval
	case "json":
		r = TypeJSON
	case "jsonb":
		r = TypeJSONB
	case "line":
		r = TypeLine
	case "lseg":
		r = TypeLineSegment
	case "macaddr":
		r = TypeMacAddr
	case "money":
		r = TypeMoney
	case "name":
		r = TypeName
	case "numeric", "decimal":
		r = TypeNumeric
	case "numrange":
		r = TypeNumRange
	case "oid":
		r = TypeOID
	case "path":
		r = TypePath
	case "point":
		r = TypePoint
	case "polygon":
		r = TypePolygon
	case "record":
		r = TypeRecord
	case "text", "character varying", "character":
		r = TypeText
	case "tid":
		r = TypeTID
	case "time", "time without time zone":
		r = TypeTime
	case "timetz", "time with time zone":
		r = TypeTimeTZ
	case "timestamp", "timestamp without time zone", "datetime":
		r = TypeTimestamp
	case "timestamptz", "timestamp with time zone":
		r = TypeTimestampTZ
	case "tsrange":
		r = TypeTsRange
	case "tsquery":
		r = TypeTsQuery
	case "tsvector":
		r = TypeTsVector
	case "tstzrange":
		r = TypeTsTzRange
	case "unknown":
		r = TypeUnknown
	case "uuid":
		r = TypeUUID
	case "varchar":
		r = TypeVarchar
	case "xid":
		r = TypeXID
	case "xml":
		r = TypeXML
	case "year":
		r = TypeYear
	default:
		logger.Debug("Unhandled data type [" + t + "] for column [" + name + "]")
	}
	return r
}

func arrayTypeFor(logger logur.LoggerFacade, name string, t string) FieldType {
	r := TypeArrayUnknown
	switch t {
	case "aclitem":
		r = TypeArrayACL
	case "bool", "boolean":
		r = TypeArrayBool
	case "bpchar":
		r = TypeArrayBPChar
	case "bytea", "blob", "binary":
		r = TypeArrayByteA
	case "cidr":
		r = TypeArrayCIDR
	case "float4":
		r = TypeArrayFloat32
	case "float8", "double precision":
		r = TypeArrayFloat64
	case "inet":
		r = TypeArrayInet
	case "int2":
		r = TypeArrayInt16
	case "int4":
		r = TypeArrayInt32
	case "int8", "integer":
		r = TypeArrayInt64
	case "json", "jsonb":
		r = TypeArrayJSON
	case "numeric":
		r = TypeArrayNumeric
	case "oid":
		r = TypeArrayOID
	case "text":
		r = TypeArrayText
	case "varchar", "character varying":
		r = TypeArrayVarchar
	case "timestamp", "timestamp without time zone":
		r = TypeArrayTimestamp
	case "timestamptz", "timestamp with time zone":
		r = TypeArrayTimestampTZ
	case "uuid":
		r = TypeArrayUUID
	default:
		logger.Debug("Unhandled array data type [" + t + "] for column [" + name + "]")
	}
	return r
}
