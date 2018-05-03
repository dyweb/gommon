package log

import (
	"fmt"
)

type FieldType uint8

const (
	UnknownType FieldType = iota
	IntType
	StringType
)

type Fields []Field

// TODO: we can specify the type in field ... how zap do it, using pointer?
type Field struct {
	Key       string
	Type      FieldType
	Int       int64
	Str       string
	Interface interface{}
}

func Int(k string, v int) Field {
	return Field{
		Key:  k,
		Type: IntType,
		Int:  int64(v),
	}
}

func Str(k string, v string) Field {
	return Field{
		Key:  k,
		Type: StringType,
		Str:  v,
	}
}

func Stringer(k string, v fmt.Stringer) Field {
	return Field{
		Key:  k,
		Type: StringType,
		Str:  v.String(),
	}
}
