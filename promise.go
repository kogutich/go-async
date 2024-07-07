package async

import "sync/atomic"

// PromiseVE represents a promise for function returning (value, error).
type PromiseVE[T any] struct {
	resultCh   <-chan T
	errCh      <-chan error
	calledOnce atomic.Bool
}

// PromiseV represents a promise for function returning value.
type PromiseV[T any] struct {
	resultCh   <-chan T
	calledOnce atomic.Bool
}

// PromiseE represents a promise for function returning error.
type PromiseE = PromiseV[error]

// Promise represents a promise for a void function.
type Promise struct {
	waitCh     <-chan struct{}
	calledOnce atomic.Bool
}

// Wait blocks until function completes, then return results.
// Note: only 1 call to this method allowed.
func (p *PromiseVE[T]) Wait() (T, error) {
	if !p.calledOnce.CompareAndSwap(false, true) {
		panic("only 1 call to Wait allowed")
	}
	select {
	case v := <-p.resultCh:
		return v, nil
	case err := <-p.errCh:
		return *new(T), err
	}
}

// Wait blocks until function completes, then return result.
// Note: only 1 call to this method allowed.
func (p *PromiseV[T]) Wait() T {
	if !p.calledOnce.CompareAndSwap(false, true) {
		panic("only 1 call to Wait allowed")
	}
	return <-p.resultCh
}

// Wait blocks until function completes.
// Note: only 1 call to this method allowed.
func (p *Promise) Wait() {
	if !p.calledOnce.CompareAndSwap(false, true) {
		panic("only 1 call to Wait allowed")
	}
	<-p.waitCh
}
