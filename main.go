package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type asyncError struct {
	RunnerName string
	Err        error
}

var wg sync.WaitGroup

func runner(name string, err chan<- asyncError, f interface{}, params ...interface{}) {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()

	if fType.Kind() != reflect.Func {
		err <- asyncError{name, errors.New("We need a valid function to work")}
		return
	}

	r := make(chan error)
	rValue := reflect.ValueOf(r)

	in := []reflect.Value{rValue}

	for _, param := range params {
		pValue := reflect.ValueOf(param)
		in = append(in, pValue)
	}

	wg.Add(2)
	go fValue.Call(in)
	go func() {
		select {
		case e := <-r:
			defer wg.Done()
			if e != nil {
				supervisor(asyncError{name, e})
			}
		}
	}()
	wg.Wait()
}

func supervisor(err asyncError) {
	fmt.Println(err.RunnerName)
	fmt.Println(err.Err)
}

func main() {
	err := make(chan asyncError)
	runner("HARD TASK", err, hardTask, []int{1, 4, 6, 7})
}

func hardTask(err chan<- error, numbers []int) {
	defer wg.Done()
	defer close(err)
	fmt.Println(numbers)
	// err <- errors.New("EROOOOOR!!!")
}
