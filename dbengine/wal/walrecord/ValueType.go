// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package walrecord

import "strconv"

type ValueType byte

const (
	ValueTypeFull    ValueType = 0
	ValueTypeChunked ValueType = 1
	ValueTypeColumn  ValueType = 2
)

var EnumNamesValueType = map[ValueType]string{
	ValueTypeFull:    "Full",
	ValueTypeChunked: "Chunked",
	ValueTypeColumn:  "Column",
}

var EnumValuesValueType = map[string]ValueType{
	"Full":    ValueTypeFull,
	"Chunked": ValueTypeChunked,
	"Column":  ValueTypeColumn,
}

func (v ValueType) String() string {
	if s, ok := EnumNamesValueType[v]; ok {
		return s
	}
	return "ValueType(" + strconv.FormatInt(int64(v), 10) + ")"
}
