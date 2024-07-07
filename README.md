# go-async

Library that provides promise-like functionality for running functions asynchronously.

## Usage
Imagine a function that takes some time to complete.  
And we want to run it asynchronously and later take the results.
```go
f := func() (string, error) {
    time.Sleep(time.Second)
    return "Hello world", nil
}
``` 

1) Start a function and save the promise (under the hood it runs in a separate goroutine):
```go
promise := async.RunVE(f)
```
2) Wait for results:
```go
value, err := promise.Wait()
```

## Naming
`RunVE` means we run function that returns **V**alue and **E**rror.  
Supported functions described [here](https://github.com/kogutich/go-async/blob/master/func_types.go). Each function type has its own `Run*` method.
