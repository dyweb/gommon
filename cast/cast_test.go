package cast

import (
    "testing"
    asst "github.com/stretchr/testify/assert"
)

func TestToStringMap(t *testing.T) {
    assert := asst.New(t)
    m := make(map[interface{}]interface{})
    m[1] = "123"
    m["milk"]= "cow"
    converted := ToStringMap(m)
    assert.Equal("cow", converted["milk"])
    assert.Equal(1, len(converted))
}
