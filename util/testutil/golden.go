package testutil

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	enabledGoldenMu sync.RWMutex
	enabledGolden   map[*testing.T]bool
)

// GenGolden check if env GOLDEN or GEN_GOLDEN is set, sometimes you need to generate test fixture in test
func GenGolden() Condition {
	return Or(EnvTrue("GOLDEN"), EnvTrue("GEN_GOLDEN"))
}

// GenGoldenT check if current test is manually set to generate golden file
func GenGoldenT(t *testing.T) Condition {
	return Or(GenGolden(), &con{
		stmt: func() (res bool, msg string, err error) {
			enabledGoldenMu.RLock()
			enabled := enabledGolden[t]
			enabledGoldenMu.RUnlock()
			if enabled {
				return enabled, "golden is enabled manually for this test", nil
			}
			return enabled, "golden is not set for this test", nil
		},
	})
}

func EnableGenGolden(t *testing.T) {
	enabledGoldenMu.Lock()
	enabledGolden[t] = true
	enabledGoldenMu.Unlock()
}

func DisableGenGolden(t *testing.T) {
	enabledGoldenMu.Lock()
	enabledGolden[t] = false
	enabledGoldenMu.Unlock()
}

func WriteOrCompare(t *testing.T, file string, data []byte) {
	if GenGoldenT(t).B() {
		WriteFixture(t, file, data)
	} else {
		b := ReadFixture(t, file)
		assert.Equal(t, b, data, file)
	}
}

func WriteOrCompareAsString(t *testing.T, file string, data string) {
	if GenGoldenT(t).B() {
		WriteFixture(t, file, []byte(data))
	} else {
		b := ReadFixture(t, file)
		assert.Equal(t, string(b), data)
	}
}

func init() {
	enabledGolden = make(map[*testing.T]bool)
}
