package json

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/dyweb/gommon/log"

	"github.com/stretchr/testify/assert"
)

var tm = time.Unix(1517861935, 0)

type entry struct {
	Level string `json:"l"`
	Time  int64  `json:"t"` // TODO: if we use time.Time, go is expecting to use string time format...
	Msg   string `json:"m"`
	Src   string `json:"s"`
}

func TestHandler_HandleLog(t *testing.T) {
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLog(log.DebugLevel, tm, "hi", log.EmptyCaller(), nil, nil)
	h.HandleLog(log.InfoLevel, tm, "hello", log.EmptyCaller(), nil, nil)
	//t.Log(b.String())
	expected := `{"l":"debug","t":1517861935,"m":"hi"}
{"l":"info","t":1517861935,"m":"hello"}
`
	assert.Equal(t, expected, b.String())

	dec := json.NewDecoder(b)
	e1 := &entry{}
	if err := dec.Decode(e1); err != nil {
		t.Fatalf("can't decode using encoding/json %s", err)
	}
	assert.Equal(t, "debug", e1.Level)
	assert.Equal(t, tm.Unix(), e1.Time)
	e2 := &entry{}
	if err := dec.Decode(e2); err != nil {
		t.Fatalf("can't decode using encoding/json %s", err)
	}
	assert.Equal(t, "info", e2.Level)
	assert.Equal(t, tm.Unix(), e2.Time)
}

func TestHandler_HandleLogWithSource(t *testing.T) {
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLog(log.DebugLevel, tm, "hi", log.Caller{File: "abc.go", Line: 12}, nil, nil)
	//t.Log(b.String())
	assert.Equal(t, `{"l":"debug","t":1517861935,"m":"hi","s":"abc.go:12"}
`, b.String())
	validJSON(t, b.Bytes())
}

func TestHandler_HandleLogWithFields(t *testing.T) {
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLog(log.DebugLevel, tm, "hi", log.EmptyCaller(), log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	}, nil)
	assert.Equal(t, `{"l":"debug","t":1517861935,"m":"hi","num":1,"str":"rts"}
`, b.String())
	validJSON(t, b.Bytes())
}

func TestHandler_HandleLogWithSourceFields(t *testing.T) {
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLog(log.DebugLevel, tm, "hi", log.Caller{File: "abc.go", Line: 12}, log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	}, nil)
	assert.Equal(t, `{"l":"debug","t":1517861935,"m":"hi","s":"abc.go:12","num":1,"str":"rts"}
`, b.String())
	validJSON(t, b.Bytes())
}

func TestJsonEscape(t *testing.T) {
	var buf bytes.Buffer
	h := New(&buf)
	h.HandleLog(log.DebugLevel, tm, `I have "quote"`, log.EmptyCaller(), nil, nil)
	validJSON(t, buf.Bytes())
}

func Test_FormatFields(t *testing.T) {
	var b []byte
	assert.Equal(t, `"num":1,"str":"rts"`, string(formatFields(b, log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	})))
}

func TestEncodeString(t *testing.T) {
	strs := []string{
		"normal string",
		`has "quote"`,
		"support 中文 me?",
	}
	// TODO: add assert here
	for _, s := range strs {
		var b []byte
		b = EncodeString(b, s)
		t.Log(string(b))
	}

	t.Run("no escape", func(t *testing.T) {
		var b []byte
		b = EncodeString(b, "nothing")
		t.Log(string(b))
	})
}

func validJSON(t *testing.T, b []byte) {
	e := &entry{}
	if err := json.Unmarshal(b, e); err != nil {
		t.Fatalf("invalid json %v", err)
	}
}
