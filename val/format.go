package val

type Format int

// matches list in browse.h
const (
	FmtBinary Format = iota + 1
	FmtBits
	FmtBool
	FmtDecimal64
	FmtEmpty
	FmtEnum
	FmtIdentityRef
	FmtInstanceRef
	FmtInt8
	FmtInt16
	FmtInt32
	FmtInt64
	FmtLeafRef
	FmtString
	FmtUInt8
	FmtUInt16
	FmtUInt32
	FmtUInt64
	FmtUnion
	FmtAny
)

const (
	FmtBinaryList      Format = iota + 1025
	FmtBitsList               //1026
	FmtBoolList               //1027
	FmtDecimalList            //1028
	FmtEmptyList              //1029
	FmtEnumList               //1030
	FmtIdentityRefList        //1031
	FmtInstanceRefList        //1032
	FmtInt8List               //1033
	FmtInt16List              //1034
	FmtInt32List              //1035
	FmtInt64List              //1036
	FmtLeafRefList            //1037
	FmtStringList             //1038
	FmtUInt8List              //1039
	FmtUInt16List             //1040
	FmtUInt32List             //1041
	FmtUInt64List             //1042
	FmtUnionList              //1043
	FmtAnyList                //1044
)

var internalTypes = map[string]Format{
	"binary":              FmtBinary,
	"bits":                FmtBits,
	"boolean":             FmtBool,
	"decimal64":           FmtDecimal64,
	"enumeration":         FmtEnum,
	"identitydef":         FmtIdentityRef,
	"instance-identifier": FmtInstanceRef,
	"int8":                FmtInt8,
	"int16":               FmtInt16,
	"int32":               FmtInt32,
	"int64":               FmtInt64,
	"leafref":             FmtLeafRef,
	"string":              FmtString,
	"uint8":               FmtUInt8,
	"uint16":              FmtUInt16,
	"uint32":              FmtUInt32,
	"uint64":              FmtUInt64,
	"union":               FmtUnion,
	"any":                 FmtAny,
}

func (f Format) String() string {
	for name, candidate := range internalTypes {
		if f == candidate {
			return name
		}
		if f-1024 == candidate {
			return name + "-list"
		}
	}
	return "?unknown?"
}
