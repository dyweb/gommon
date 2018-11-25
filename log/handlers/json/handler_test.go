// +build ignore

// TODO: recover the test after refactor
package json

import (
	"bytes"
	"testing"
	"time"

	"github.com/dyweb/gommon/log"

	"encoding/json"

	asst "github.com/stretchr/testify/assert"
)

var tm = time.Unix(1517861935, 0)

type entry struct {
	Level string `json:"l"`
	Time  int64  `json:"t"` // TODO: if we use time.Time, go is expecting to use string time format...
	Msg   string `json:"m"`
	Src   string `json:"s"`
}

func TestHandler_HandleLog(t *testing.T) {
	assert := asst.New(t)
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLog(log.DebugLevel, tm, "hi")
	h.HandleLog(log.InfoLevel, tm, "hello")
	//t.Log(b.String())
	expected := `{"l":"debug","t":1517861935,"m":"hi"}
{"l":"info","t":1517861935,"m":"hello"}
`
	assert.Equal(expected, b.String())

	dec := json.NewDecoder(b)
	e1 := &entry{}
	if err := dec.Decode(e1); err != nil {
		t.Fatalf("can't decode using encoding/json %s", err)
	}
	assert.Equal("debug", e1.Level)
	assert.Equal(tm.Unix(), e1.Time)
	e2 := &entry{}
	if err := dec.Decode(e2); err != nil {
		t.Fatalf("can't decode using encoding/json %s", err)
	}
	assert.Equal("info", e2.Level)
	assert.Equal(tm.Unix(), e2.Time)
}

func TestHandler_HandleLogWithSource(t *testing.T) {
	assert := asst.New(t)
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLogWithSource(log.DebugLevel, tm, "hi", "abc.go:12")
	//t.Log(b.String())
	assert.Equal(`{"l":"debug","t":1517861935,"m":"hi","s":"abc.go:12"}
`, b.String())
	validJSON(t, b.Bytes())
}

func TestHandler_HandleLogWithFields(t *testing.T) {
	assert := asst.New(t)
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLogWithFields(log.DebugLevel, tm, "hi", log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	})
	assert.Equal(`{"l":"debug","t":1517861935,"m":"hi","num":1,"str":"rts"}
`, b.String())
	validJSON(t, b.Bytes())
}

func TestHandler_HandleLogWithSourceFields(t *testing.T) {
	assert := asst.New(t)
	b := &bytes.Buffer{}
	h := New(b)
	h.HandleLogWithSourceFields(log.DebugLevel, tm, "hi", "abc.go:12", log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	})
	assert.Equal(`{"l":"debug","t":1517861935,"m":"hi","s":"abc.go:12","num":1,"str":"rts"}
`, b.String())
	validJSON(t, b.Bytes())
}

func Test_FormatHead(t *testing.T) {
	assert := asst.New(t)
	var b []byte
	assert.Equal(`{"l":"debug","t":1517861935,"m":"hi"`, string(formatHead(b, log.DebugLevel, tm.Unix(), "hi")))
}

func Test_FormatFields(t *testing.T) {
	assert := asst.New(t)
	var b []byte
	assert.Equal(`"num":1,"str":"rts"`, string(formatFields(b, log.Fields{
		log.Int("num", 1),
		log.Str("str", "rts"),
	})))
}

func validJSON(t *testing.T, b []byte) {
	e := &entry{}
	if err := json.Unmarshal(b, e); err != nil {
		t.Fatalf("invalid json %v", err)
	}
}
