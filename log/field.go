package log

type Fields []Field

// TODO: we can specify the type in field ... how zap do it, using pointer?
type Field struct {
	Key   string
	Value interface{}
}
