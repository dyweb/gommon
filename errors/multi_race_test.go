// +build !race

package errors

import (
	"os"
	"sync"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

// NOTE: when adding build tag, remember to have a blank line between it and package name, otherwise it is treated as comment
func TestMultiErr_Append(t *testing.T) {
	assert := asst.New(t)
	merr := NewMultiErr()
	nRoutine := 10
	errPerRoutine := 20
	var wg sync.WaitGroup
	wg.Add(nRoutine)
	for i := 0; i < nRoutine; i++ {
		go func() {
			for j := 0; j < errPerRoutine; j++ {
				merr.Append(os.ErrClosed)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	// due to data race, the actual number of errors should be smaller than what we expected
	// NOTE: we add = because sometimes there is no race ....
	assert.True(nRoutine*errPerRoutine >= len(merr.Errors()))
	t.Logf("expect %d actual %d", nRoutine*errPerRoutine, len(merr.Errors()))
}
