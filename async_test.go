package async

import (
	"errors"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunVE(t *testing.T) {
	fnGenerate := func(i int, e error) FuncVE[int] {
		t.Helper()
		return func() (int, error) {
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Nanosecond)
			return i, e
		}
	}

	t.Run("async success", func(t *testing.T) {
		promise := RunVE(fnGenerate(10, nil))
		result, err := promise.Wait()
		require.NoError(t, err)
		assert.Equal(t, 10, result)
	})

	t.Run("async success with multiple promises", func(t *testing.T) {
		promise1 := RunVE(fnGenerate(1, nil))
		promise2 := RunVE(fnGenerate(2, nil))
		promise3 := RunVE(fnGenerate(3, nil))

		result1, err1 := promise1.Wait()
		require.NoError(t, err1)
		assert.Equal(t, 1, result1)
		result2, err2 := promise2.Wait()
		require.NoError(t, err2)
		assert.Equal(t, 2, result2)
		result3, err3 := promise3.Wait()
		require.NoError(t, err3)
		assert.Equal(t, 3, result3)
	})

	t.Run("async with error", func(t *testing.T) {
		promise := RunVE(fnGenerate(0, errors.New("some error")))
		_, err := promise.Wait()
		assert.ErrorContains(t, err, "some error")
	})

	t.Run("async with panic", func(t *testing.T) {
		promise := RunVE(func() (int, error) {
			panic("runtime error")
		})
		_, err := promise.Wait()
		assert.ErrorContains(t, err, "recovered: runtime error")
	})

	t.Run("multiple Wait() calls", func(t *testing.T) {
		promise := RunVE(fnGenerate(123, nil))
		result, err := promise.Wait()
		require.NoError(t, err)
		assert.Equal(t, 123, result)
		assert.Panics(t, func() { _, _ = promise.Wait() })
	})
}

func TestRunV(t *testing.T) {
	fnGenerate := func(i int) FuncV[int] {
		t.Helper()
		return func() int {
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Nanosecond)
			return i
		}
	}

	t.Run("async with multiple promises", func(t *testing.T) {
		promise1 := RunV(fnGenerate(1))
		promise2 := RunV(fnGenerate(2))
		promise3 := RunV(fnGenerate(3))

		assert.Equal(t, 1, promise1.Wait())
		assert.Equal(t, 2, promise2.Wait())
		assert.Equal(t, 3, promise3.Wait())
	})

	t.Run("multiple Wait() calls", func(t *testing.T) {
		promise := RunV(fnGenerate(123))
		assert.Equal(t, 123, promise.Wait())
		assert.Panics(t, func() { _ = promise.Wait() })
	})
}

func TestRunE(t *testing.T) {
	fnGenerate := func(e error) FuncE {
		t.Helper()
		return func() error {
			time.Sleep(time.Duration(rand.Int63n(100)) * time.Nanosecond)
			return e
		}
	}

	t.Run("async success", func(t *testing.T) {
		promise := RunE(fnGenerate(nil))
		assert.NoError(t, promise.Wait())
	})

	t.Run("async error with multiple promises", func(t *testing.T) {
		promise1 := RunE(fnGenerate(errors.New("1")))
		promise2 := RunE(fnGenerate(errors.New("2")))
		promise3 := RunE(fnGenerate(errors.New("3")))

		require.ErrorContains(t, promise1.Wait(), "1")
		require.ErrorContains(t, promise2.Wait(), "2")
		require.ErrorContains(t, promise3.Wait(), "3")
	})

	t.Run("async with panic", func(t *testing.T) {
		promise := RunE(func() error {
			panic("runtime error")
		})
		assert.ErrorContains(t, promise.Wait(), "recovered: runtime error")
	})

	t.Run("multiple Wait() calls", func(t *testing.T) {
		promise := RunE(fnGenerate(nil))
		require.NoError(t, promise.Wait())
		assert.Panics(t, func() { _ = promise.Wait() })
	})
}

func TestRun(t *testing.T) {
	t.Run("async with multiple promises", func(t *testing.T) {
		var i int64
		testFunc := func() {
			atomic.AddInt64(&i, 1)
		}
		promise1 := Run(testFunc)
		promise2 := Run(testFunc)
		promise3 := Run(testFunc)
		promise1.Wait()
		promise2.Wait()
		promise3.Wait()
		assert.Equal(t, int64(3), i)
	})

	t.Run("multiple Wait() calls", func(t *testing.T) {
		testFunc := func() {}
		promise := Run(testFunc)
		assert.NotPanics(t, promise.Wait)
		assert.Panics(t, promise.Wait)
	})
}
