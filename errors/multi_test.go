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

func TestMultiErr_AppendReturn(t *testing.T) {
	// NOTE: inspired by https://github.com/uber-go/multierr/issues/21
	errs := map[string]func() MultiErr{
		"unsafe": NewMultiErr,
		"safe":   NewMultiErrSafe,
	}
	for name := range errs {
		t.Run(name, func(t *testing.T) {
			assert := asst.New(t)
			merr := errs[name]()
			assert.False(merr.Append(nil))
			assert.True(merr.Append(os.ErrPermission))
			assert.False(merr.Append(nil))
		})
	}
}

func TestMultiErr_Flatten(t *testing.T) {
	errs := map[string]func() MultiErr{
		"unsafe": NewMultiErr,
		"safe":   NewMultiErrSafe,
	}
	for name := range errs {
		t.Run(name, func(t *testing.T) {
			assert := asst.New(t)
			merr := errs[name]()
			merr.Append(os.ErrPermission)
			merr.Append(os.ErrClosed)
			assert.Equal(2, len(merr.Errors()))

			merr2 := errs[name]()
			merr2.Append(merr)
			merr2.Append(os.ErrNotExist)
			merr2.Append(nil) // nil is not appended
			assert.Equal(3, len(merr2.Errors()))
			t.Log(merr2.Error())
		})
	}
}

func TestMultiErr_ErrorOrNil(t *testing.T) {
	assert := asst.New(t)

	merr := NewMultiErr()
	assert.Nil(merr.ErrorOrNil())

	merr.Append(os.ErrPermission)
	assert.NotNil(merr.ErrorOrNil())
}
