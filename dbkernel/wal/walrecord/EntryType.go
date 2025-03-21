// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package walrecord

import "strconv"

type EntryType byte

const (
	EntryTypeKV      EntryType = 0
	EntryTypeChunked EntryType = 1
	EntryTypeRow     EntryType = 2
)

var EnumNamesEntryType = map[EntryType]string{
	EntryTypeKV:      "KV",
	EntryTypeChunked: "Chunked",
	EntryTypeRow:     "Row",
}

var EnumValuesEntryType = map[string]EntryType{
	"KV":      EntryTypeKV,
	"Chunked": EntryTypeChunked,
	"Row":     EntryTypeRow,
}

func (v EntryType) String() string {
	if s, ok := EnumNamesEntryType[v]; ok {
		return s
	}
	return "EntryType(" + strconv.FormatInt(int64(v), 10) + ")"
}
