package async

import "fmt"

// RunVE runs FuncVE asynchronously and returns a promise for obtaining results.
func RunVE[T any](fn FuncVE[T]) *PromiseVE[T] {
	resultCh := make(chan T, 1)
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("recovered: %v", r)
			}
		}()
		v, err := fn()
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- v
	}()
	return &PromiseVE[T]{
		resultCh: resultCh,
		errCh:    errCh,
	}
}

// RunV runs FuncV asynchronously and returns a promise for obtaining results.
func RunV[T any](fn FuncV[T]) *PromiseV[T] {
	resultCh := make(chan T, 1)
	go func() {
		resultCh <- fn()
	}()
	return &PromiseV[T]{resultCh: resultCh}
}

// RunE runs FuncE asynchronously and returns a promise for obtaining results.
func RunE(fn FuncE) *PromiseE {
	errCh := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("recovered: %v", r)
			}
		}()
		errCh <- fn()
	}()
	return &PromiseE{resultCh: errCh}
}

// Run runs Func asynchronously and returns a promise for obtaining results.
func Run(fn Func) *Promise {
	waitCh := make(chan struct{}, 1)
	go func() {
		fn()
		waitCh <- struct{}{}
	}()
	return &Promise{waitCh: waitCh}
}
