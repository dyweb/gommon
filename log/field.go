package log

import (
	"fmt"
)

// FieldType avoids doing type assertion or calling reflection
// TODO: difference between type assertion and reflection
type FieldType uint8

const (
	UnknownType FieldType = iota
	IntType
	StringType
)

// Fields is a slice of Field
type Fields []Field

// Field is based on uber-go/zap https://github.com/uber-go/zap/blob/master/zapcore/field.go
// It can be treated as a Union, the value is stored in either Int, Str or Interface
// TODO: interface{} is actually not used ...
type Field struct {
	Key  string
	Type FieldType

	// values
	Int       int64
	Str       string
	Interface interface{}
}

// Int creates a field with int value, it uses int64 internally
func Int(k string, v int) Field {
	return Field{
		Key:  k,
		Type: IntType,
		Int:  int64(v),
	}
}

// TODO: reuse same int64 and rely on type? ...
//func Float(k string, v float64) Field {
//
//}

// Str creates a field with string value
func Str(k string, v string) Field {
	return Field{
		Key:  k,
		Type: StringType,
		Str:  v,
	}
}

// Stringer calls the String() method and stores return value
func Stringer(k string, v fmt.Stringer) Field {
	return Field{
		Key:  k,
		Type: StringType,
		Str:  v.String(),
	}
}

// copyFields make a copy of the slice so modifying one won't have effect on another,
func copyFields(fields Fields) Fields {
	copied := make([]Field, len(fields))
	copy(copied, fields)
	return copied
}
