package async

// Function types supported by this library.
// V - means function returns value, E - error.
type (
	FuncVE[T any] func() (T, error)
	FuncV[T any]  func() T
	FuncE         FuncV[error]
	Func          func()
)
