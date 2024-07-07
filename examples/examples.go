package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kogutich/go-async"
)

func main() {
	f1 := func() (string, error) {
		time.Sleep(time.Second)
		return "Hello", nil
	}
	f2 := func() (string, error) {
		time.Sleep(time.Second)
		return "world!", nil
	}
	// async run
	promise1 := async.RunVE(f1)
	promise2 := async.RunVE(f2)
	// wait
	value1, err := promise1.Wait()
	if err != nil {
		log.Fatal(err)
	}
	value2, err := promise2.Wait()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", value1, value2)
}
