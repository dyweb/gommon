package errors

import (
	"os"
	"sync"
	"testing"

	asst "github.com/stretchr/testify/assert"
)

func TestMultiErrSafe_Append(t *testing.T) {
	assert := asst.New(t)
	merr := NewMultiErrSafe()
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
	assert.Equal(nRoutine*errPerRoutine, len(merr.Errors()))
}

func TestMultiErr_Flatten(t *testing.T) {
	errs := map[string]MultiErr{
		"unsafe": NewMultiErr(),
		"safe":   NewMultiErrSafe(),
	}
	for name := range errs {
		t.Run(name, func(t *testing.T) {
			assert := asst.New(t)
			merr := NewMultiErr()
			merr.Append(os.ErrPermission)
			merr.Append(os.ErrClosed)
			assert.Equal(2, len(merr.Errors()))

			merr2 := NewMultiErr()
			merr2.Append(merr)
			merr2.Append(os.ErrNotExist)
			assert.Equal(3, len(merr2.Errors()))
			t.Log(merr2.Error())
		})
	}
}
